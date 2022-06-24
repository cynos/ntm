package tools

import (
	"context"
	"fmt"
	"log"

	"github.com/go-redis/redis/v8"
)

type RedisConfig struct {
	Host     string
	Port     string
	DB       int
	Username string
	Password string
}

func RedisClient(config RedisConfig) *redis.Client {
	rdbClient := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", config.Host, config.Port),
		Username: config.Username,
		Password: config.Password, // no password set
		DB:       config.DB,       // use default DB
	})
	_, err := rdbClient.Ping(context.TODO()).Result()
	if err != nil {
		log.Fatal(err)
	}
	return rdbClient
}
