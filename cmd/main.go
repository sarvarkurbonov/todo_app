package main

import (
	todo "Todo_rest_api"
	"Todo_rest_api/pkg/cache"
	"Todo_rest_api/pkg/handler"
	"Todo_rest_api/pkg/repository"
	"Todo_rest_api/pkg/service"
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func main() {

	logrus.SetFormatter(new(logrus.JSONFormatter))

	if err := initConfig(); err != nil {
		logrus.Fatal("Error reading config:", err.Error())
	}

	if err := godotenv.Load(); err != nil {
		logrus.Fatal("Error loading .env file:", err.Error())
	}

	db, err := repository.NewPostgresDB(repository.Config{
		Host:     viper.GetString("db.host"),
		Port:     viper.GetString("db.port"),
		Username: viper.GetString("db.username"),
		Password: os.Getenv("DB_PASSWORD"),
		DbName:   viper.GetString("db.dbname"),
		SslMode:  viper.GetString("db.sslMode"),
	})
	rdb := cache.New()
	defer func(rdb *redis.Client) {
		err := rdb.Close()
		if err != nil {

		}
	}(rdb)
	if err != nil {
		logrus.Fatal("Error connecting to database:", err.Error())
	}
	repos := repository.NewRepository(db, rdb)
	services := service.NewService(repos)
	handlers := handler.NewHandler(services)

	srv := new(todo.Server)
	go func() {
		if err := srv.Run(viper.GetString("port"), handlers.InitRoutes()); err != nil {
			logrus.Fatal("Error starting server:", err.Error())
		}
	}()
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logrus.Println("Shutting down server...")

	if err := srv.Shutdown(context.Background()); err != nil {
		logrus.WithError(err).Fatal("server forced to shutdown")
	}
	logrus.Info("Server shut down")
	if err := db.Close(); err != nil {
		logrus.WithError(err).Fatal("Error closing database connection")
	}

}

func initConfig() error {
	viper.AddConfigPath("configs")
	viper.SetConfigName("config")
	return viper.ReadInConfig()
}
