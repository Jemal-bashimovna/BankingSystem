package main

import (
	"bankingsystem/constants"
	"bankingsystem/deps"
	"bankingsystem/pkg/repository"
	"bankingsystem/pkg/repository/listeners"
	"os"

	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

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

	transactionRepo := repository.NewTransactionRepository(db, rdb)

	brokers := viper.GetString("kafka.brokers")
	groupId := "banking-system-consumer"

	depCons := deps.NewConsumer(brokers, groupId, []string{constants.Deposit})
	transCons := deps.NewConsumer(brokers, groupId, []string{constants.Transfer})
	withdrawCons := deps.NewConsumer(brokers, groupId, []string{constants.Withdraw})

	depositConsumer := listeners.NewDepositConsumer(depCons, groupId, transactionRepo)
	withdrawConsumer := listeners.NewWithdrawConsumer(withdrawCons, groupId, transactionRepo)
	transferConsumer := listeners.NewTransferConsumer(transCons, groupId, transactionRepo)

	go depositConsumer.StartListening()
	go withdrawConsumer.StartListening()
	go transferConsumer.StartListening()

	select {}
}

func initConfig() error {
	viper.AddConfigPath("configs")
	viper.SetConfigName("config")
	return viper.ReadInConfig()
}
