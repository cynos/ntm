package news

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ntm/internal/infrastructure/cache"
	"github.com/ntm/internal/pkg/common/http/response"
	"github.com/ntm/internal/tools"
)

type HTTPController struct {
	newsUseCase UseCase
	cacher      cache.Cacher
}

func NewHTTPController(newsUseCase UseCase, cacher cache.Cacher) *HTTPController {
	return &HTTPController{
		newsUseCase: newsUseCase,
		cacher:      cacher,
	}
}

func (controller *HTTPController) FindAll(c *gin.Context) {
	var filter Filter
	if err := c.ShouldBindQuery(&filter); err != nil {
		response.ErrorWithMessage(c, http.StatusBadRequest, "invalid parameters", err)
		return
	}

	// get from cache
	cache_key := tools.MD5([]byte(fmt.Sprintf("news:%s", c.Request.URL.Query().Encode())))
	if controller.cacher.IsExist(cache_key) {
		log.Println("news | findAll | serve by redis")
		payload := []News{}
		if err := json.Unmarshal([]byte(controller.cacher.Get(cache_key).(string)), &payload); err != nil {
			response.Error(c, http.StatusInternalServerError, err)
			return
		}
		response.Success(c, http.StatusOK, payload)
		return
	}

	// create context
	ctx := context.WithValue(context.Background(), ContextKey("news_filter"), filter)

	// get from db
	news, err := controller.newsUseCase.FindAll(ctx)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err)
		return
	}

	// save in cache
	cache_val, _ := json.Marshal(news)
	if err := controller.cacher.Put(cache_key, cache_val, 60); err != nil {
		log.Println(err.Error())
	}

	response.Success(c, http.StatusOK, news)
}

func (controller *HTTPController) FindByID(c *gin.Context) {
	id := c.Param("id")

	// get from cache
	cache_key := tools.MD5([]byte(fmt.Sprintf("news_id:" + id)))
	if controller.cacher.IsExist(cache_key) {
		log.Println("news | findByID | serve by redis")
		payload := News{}
		if err := json.Unmarshal([]byte(controller.cacher.Get(cache_key).(string)), &payload); err != nil {
			response.Error(c, http.StatusInternalServerError, err)
			return
		}
		response.Success(c, http.StatusOK, payload)
		return
	}

	// get from db
	news, err := controller.newsUseCase.FindByID(c.Request.Context(), tools.StringsToInt(id))
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err)
		return
	}

	// save in cache
	cache_val, _ := json.Marshal(news)
	if err := controller.cacher.Put(cache_key, cache_val, 60); err != nil {
		log.Println(err.Error())
	}

	response.Success(c, http.StatusOK, news)
}

func (controller *HTTPController) Add(c *gin.Context) {
	var err error
	var dto NewsDTO

	err = c.Bind(&dto)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err)
		return
	}

	_, err = controller.newsUseCase.Save(c.Request.Context(), dto)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err)
		return
	}

	// flush cache
	err = controller.cacher.Flush()
	if err != nil {
		log.Println(err.Error())
	}

	response.Success(c, http.StatusOK, nil)
}

func (controller *HTTPController) Update(c *gin.Context) {
	var err error
	var dto NewsDTO

	err = c.Bind(&dto)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err)
		return
	}

	// get news by id
	news, err := controller.newsUseCase.FindByID(c.Request.Context(), tools.StringsToInt(c.Param("id")))
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err)
		return
	}

	// update tag
	dto.ID = news.ID
	_, err = controller.newsUseCase.Save(c.Request.Context(), dto)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err)
		return
	}

	// flush cache
	err = controller.cacher.Flush()
	if err != nil {
		log.Println(err.Error())
	}

	response.Success(c, http.StatusOK, nil)
}

func (controller *HTTPController) Delete(c *gin.Context) {
	var err error

	err = controller.newsUseCase.Delete(c.Request.Context(), tools.StringsToInt(c.Param("id")))
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err)
		return
	}

	// flush cache
	err = controller.cacher.Flush()
	if err != nil {
		log.Println(err.Error())
	}

	response.Success(c, http.StatusOK, nil)
}
