package tag

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
	tagUseCase UseCase
	cacher     cache.Cacher
}

func NewHTTPController(tagUseCase UseCase, cacher cache.Cacher) *HTTPController {
	return &HTTPController{
		tagUseCase: tagUseCase,
		cacher:     cacher,
	}
}

func (controller *HTTPController) FindAll(c *gin.Context) {
	var filter Filter
	if err := c.ShouldBindQuery(&filter); err != nil {
		response.ErrorWithMessage(c, http.StatusBadRequest, "invalid parameters", err)
		return
	}

	// get from cache
	cache_key := tools.MD5([]byte(fmt.Sprintf("tags:%s", c.Request.URL.Query().Encode())))
	if controller.cacher.IsExist(cache_key) {
		log.Println("tag | findAll | serve by redis")
		payload := []Tag{}
		if err := json.Unmarshal([]byte(controller.cacher.Get(cache_key).(string)), &payload); err != nil {
			response.Error(c, http.StatusInternalServerError, err)
			return
		}
		response.Success(c, http.StatusOK, payload)
		return
	}

	// create context
	ctx := context.WithValue(context.Background(), ContextKey("tags_filter"), filter)

	// get from db
	tags, err := controller.tagUseCase.FindAll(ctx)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err)
		return
	}

	// save in cache
	cache_val, _ := json.Marshal(tags)
	if err := controller.cacher.Put(cache_key, cache_val, 600); err != nil {
		log.Println(err.Error())
	}

	response.Success(c, http.StatusOK, tags)
}

func (controller *HTTPController) FindByID(c *gin.Context) {
	id := c.Param("id")

	// get from cache
	cache_key := tools.MD5([]byte(fmt.Sprintf("tag_id:" + id)))
	if controller.cacher.IsExist(cache_key) {
		log.Println("tag | findByID | serve by redis")
		payload := Tag{}
		if err := json.Unmarshal([]byte(controller.cacher.Get(cache_key).(string)), &payload); err != nil {
			response.Error(c, http.StatusInternalServerError, err)
			return
		}
		response.Success(c, http.StatusOK, payload)
		return
	}

	// get from db
	tag, err := controller.tagUseCase.FindByID(c.Request.Context(), tools.StringsToInt(id))
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err)
		return
	}

	// save in cache
	cache_val, _ := json.Marshal(tag)
	if err := controller.cacher.Put(cache_key, cache_val, 600); err != nil {
		log.Println(err.Error())
	}

	response.Success(c, http.StatusOK, tag)
}

func (controller *HTTPController) Add(c *gin.Context) {
	var err error
	var tag Tag

	err = c.Bind(&tag)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err)
		return
	}

	_, err = controller.tagUseCase.Add(c.Request.Context(), tag)
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
	var tag Tag

	err = c.Bind(&tag)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err)
		return
	}

	_, err = controller.tagUseCase.Update(c.Request.Context(), tag, tools.StringsToInt(c.Param("id")))
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

	id := c.Param("id")
	err = controller.tagUseCase.Delete(c.Request.Context(), tools.StringsToInt(id))
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
