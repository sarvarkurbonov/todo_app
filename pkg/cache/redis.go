package cache

import (
	"context"

	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
)

func New() *redis.Client {
	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "CHANGE_ME", // your password
		DB:       0,
	})
	if err := rdb.Ping(context.Background()).Err(); err != nil {
		logrus.WithError(err).Fatal("redis ping failed")
	}
	return rdb
}
