package main

import (
	"bankingsystem/cmd/server"
	"bankingsystem/deps"
	"bankingsystem/pkg/handler"
	"bankingsystem/pkg/repository"
	"bankingsystem/pkg/service"
	"context"
	"os"

	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

var ctx = context.Background()

func main() {

	if err := initConfig(); err != nil {
		logrus.Fatalf("Error initializing configs: %s", err)
	}

	if err := godotenv.Load(); err != nil {
		logrus.Fatalf("Error loading env values: %s", err)
	}

	db, err := deps.NewPostgres(deps.Config{
		Host:     viper.GetString("db.host"),
		Port:     viper.GetString("db.port"),
		User:     viper.GetString("db.user"),
		DBName:   viper.GetString("db.dbname"),
		Password: os.Getenv("DB_PASSWORD"),
		SSLMode:  viper.GetString("db.sslmode"),
	})
	if err != nil {
		logrus.Fatalf("Error connecting to postgres db: %s", err)
	}

	rdb := deps.NewRedis(deps.RedisConfig{
		Host:     viper.GetString("redis.host"),
		Port:     viper.GetString("redis.port"),
		Password: os.Getenv("REDIS_PASSWORD"),
		DB:       viper.GetInt("redis.db"),
	})

	repo := repository.NewRepository(db, rdb, ctx)
	services := service.NewService(repo)
	handlers := handler.NewHandler(services)

	srv := new(server.Server)

	if err := srv.Run(viper.GetString("port"), handlers.InitRoutes()); err != nil {
		logrus.Fatalf("Error running server: %s", err)
	}
}

func initConfig() error {
	viper.AddConfigPath("configs")
	viper.SetConfigName("config")
	return viper.ReadInConfig()
}
