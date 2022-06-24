package bootstrap

import (
	"github.com/gin-gonic/gin"
	"github.com/ntm/internal/domain/news"
	"github.com/ntm/internal/domain/tag"
	"github.com/ntm/internal/domain/topic"
)

func registerTagRoute(r *gin.Engine, tagController *tag.HTTPController) {
	tagRouter := r.Group("/v1/tag")
	tagRouter.GET("/", tagController.FindAll)
	tagRouter.GET("/:id", tagController.FindByID)
	tagRouter.POST("/", tagController.Add)
	tagRouter.PUT("/:id", tagController.Update)
	tagRouter.DELETE("/:id", tagController.Delete)
}

func registerNewsRoute(r *gin.Engine, newsController *news.HTTPController) {
	newsRouter := r.Group("/v1/news")
	newsRouter.GET("/", newsController.FindAll)
	newsRouter.GET("/:id", newsController.FindByID)
	newsRouter.POST("/", newsController.Add)
	newsRouter.PUT("/:id", newsController.Update)
	newsRouter.DELETE("/:id", newsController.Delete)
}

func registerTopicRoute(r *gin.Engine, topicController *topic.HTTPController) {
	topicRouter := r.Group("/v1/topic")
	topicRouter.GET("/", topicController.FindAll)
	topicRouter.GET("/:id", topicController.FindByID)
	topicRouter.POST("/", topicController.Add)
	topicRouter.PUT("/:id", topicController.Update)
	topicRouter.DELETE("/:id", topicController.Delete)
}
