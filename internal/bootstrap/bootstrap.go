package bootstrap

import (
	"os"

	"github.com/gin-gonic/gin"
	"github.com/ntm/internal/domain/news"
	"github.com/ntm/internal/domain/tag"
	"github.com/ntm/internal/domain/topic"
	"github.com/ntm/internal/infrastructure/cache"
	"github.com/ntm/internal/tools"
	"gorm.io/gorm"
)

var router *gin.Engine
var db *gorm.DB
var cacher *cache.CacheRedis

func init() {
	// init router
	router = gin.Default()

	// init database
	db, _ = tools.DBClient(tools.DBConfig{
		Host:            os.Getenv("DB_HOST"),
		Port:            os.Getenv("DB_PORT"),
		Name:            os.Getenv("DB_NAME"),
		User:            os.Getenv("DB_USER"),
		Password:        os.Getenv("DB_PASSWORD"),
		ApplicationName: os.Getenv("news_topic_management"),
		ConnectTimeout:  tools.StringsToInt(os.Getenv("DB_CONN_TIMEOUT")),
		MaxIdleConn:     tools.StringsToInt(os.Getenv("DB_MAX_IDLE_CONN")),
		MaxOpenConn:     tools.StringsToInt(os.Getenv("DB_MAX_OPEN_CONN")),
	})

	// init redis
	redisClient := tools.RedisClient(tools.RedisConfig{
		Host:     os.Getenv("REDIS_HOST"),
		Port:     os.Getenv("REDIS_PORT"),
		DB:       tools.StringsToInt(os.Getenv("REDIS_DB")),
		Username: os.Getenv("REDIS_USERNAME"),
		Password: os.Getenv("REDIS_PASSWORD"),
	})
	cacher = cache.NewCacheRedis(redisClient)

	// migrate tables
	db.AutoMigrate(
		tag.Tag{},
		news.News{},
		topic.Topic{},
	)
}

func LaunchApp() {
	// register services
	registerTagAPIService()
	registerNewsAPIService()
	registerTopicAPIService()

	// start server
	router.Run(os.Getenv("APP_PORT"))
}
