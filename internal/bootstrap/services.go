package bootstrap

import (
	"github.com/ntm/internal/domain/news"
	"github.com/ntm/internal/domain/tag"
	"github.com/ntm/internal/domain/topic"
)

func registerTagAPIService() {
	// Initialize Tag Service
	tagRepo := tag.NewRepository(db)
	tagUseCase := tag.NewUseCase(tagRepo)
	tagController := tag.NewHTTPController(tagUseCase, cacher)
	// Build API
	registerTagRoute(router, tagController)
}

func registerNewsAPIService() {
	// Initialize Tag Service
	newsRepo := news.NewRepository(db)
	newsUseCase := news.NewUseCase(newsRepo)
	newsController := news.NewHTTPController(newsUseCase, cacher)
	// Build API
	registerNewsRoute(router, newsController)
}

func registerTopicAPIService() {
	// Initialize Tag Service
	topicRepo := topic.NewRepository(db)
	topicUseCase := topic.NewUseCase(topicRepo)
	topicController := topic.NewHTTPController(topicUseCase, cacher)
	// Build API
	registerTopicRoute(router, topicController)
}
