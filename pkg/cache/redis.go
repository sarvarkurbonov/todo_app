package cache

import (
	"context"
	"os"

	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func New() *redis.Client {
	rdb := redis.NewClient(&redis.Options{
		Addr:     viper.GetString("REDIS_HOST"),
		Password: os.Getenv("REDIS_PASSWORD"), // your password
		DB:       0,
	})
	if err := rdb.Ping(context.Background()).Err(); err != nil {
		logrus.WithError(err).Fatal("redis ping failed")
	}
	return rdb
}
