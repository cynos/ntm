package topic

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
	topicUseCase UseCase
	cacher       cache.Cacher
}

func NewHTTPController(topicUseCase UseCase, cacher cache.Cacher) *HTTPController {
	return &HTTPController{
		topicUseCase: topicUseCase,
		cacher:       cacher,
	}
}

func (controller *HTTPController) FindAll(c *gin.Context) {
	var filter Filter
	if err := c.ShouldBindQuery(&filter); err != nil {
		response.ErrorWithMessage(c, http.StatusBadRequest, "invalid parameters", err)
		return
	}

	// get from cache
	cache_key := tools.MD5([]byte(fmt.Sprintf("topics:%s", c.Request.URL.Query().Encode())))
	if controller.cacher.IsExist(cache_key) {
		log.Println("topic | findAll | serve by redis")
		payload := []Topic{}
		if err := json.Unmarshal([]byte(controller.cacher.Get(cache_key).(string)), &payload); err != nil {
			response.Error(c, http.StatusInternalServerError, err)
			return
		}
		response.Success(c, http.StatusOK, payload)
		return
	}

	// create context
	ctx := context.WithValue(context.Background(), ContextKey("topics_filter"), filter)

	// get from db
	topics, err := controller.topicUseCase.FindAll(ctx)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err)
		return
	}

	// save in cache
	cache_val, _ := json.Marshal(topics)
	if err := controller.cacher.Put(cache_key, cache_val, 600); err != nil {
		log.Println(err.Error())
	}

	response.Success(c, http.StatusOK, topics)
}

func (controller *HTTPController) FindByID(c *gin.Context) {
	id := c.Param("id")

	// get from cache
	cache_key := tools.MD5([]byte(fmt.Sprintf("topic_id:" + id)))
	if controller.cacher.IsExist(cache_key) {
		log.Println("topic | findByID | serve by redis")
		payload := Topic{}
		if err := json.Unmarshal([]byte(controller.cacher.Get(cache_key).(string)), &payload); err != nil {
			response.Error(c, http.StatusInternalServerError, err)
			return
		}
		response.Success(c, http.StatusOK, payload)
		return
	}

	// get from db
	topic, err := controller.topicUseCase.FindByID(c.Request.Context(), tools.StringsToInt(id))
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err)
		return
	}

	// save in cache
	cache_val, _ := json.Marshal(topic)
	if err := controller.cacher.Put(cache_key, cache_val, 600); err != nil {
		log.Println(err.Error())
	}

	response.Success(c, http.StatusOK, topic)
}

func (controller *HTTPController) Add(c *gin.Context) {
	var err error
	var topic Topic

	err = c.Bind(&topic)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err)
		return
	}

	_, err = controller.topicUseCase.Add(c.Request.Context(), topic)
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
	var topic Topic

	err = c.Bind(&topic)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err)
		return
	}

	_, err = controller.topicUseCase.Update(c.Request.Context(), topic, tools.StringsToInt(c.Param("id")))
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
	err = controller.topicUseCase.Delete(c.Request.Context(), tools.StringsToInt(id))
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
