package service

import (
	"bankingsystem/constants"
	"bankingsystem/deps"
	"bankingsystem/models"
	"bankingsystem/pkg/repository"
	"encoding/json"
	"fmt"

	"github.com/spf13/viper"
)

type TransactionService struct {
	repo repository.Transactions
}

func NewTransactionService(repo repository.Transactions) *TransactionService {
	return &TransactionService{
		repo: repo,
	}
}

func (s *TransactionService) Deposit(id int, sum models.InputDeposit) error {
	p := deps.NewProducer(viper.GetString("kafka.brokers"))
	message, err := json.Marshal(sum)
	if err != nil {
		return fmt.Errorf("failed sending message to producer: %s", err)
	}
	err = p.SendMessage([]byte(message), constants.Deposit)

	if err != nil {
		return fmt.Errorf("failed sending message to producer: %s", err)
	}
	defer p.Close()
	return nil
}

func (s *TransactionService) Withdraw(id int, sum models.InputWithdraw) error {
	p := deps.NewProducer(viper.GetString("kafka.brokers"))
	message, err := json.Marshal(sum)
	if err != nil {
		return fmt.Errorf("failed sending message to producer: %s", err)
	}
	err = p.SendMessage([]byte(message), constants.Withdraw)
	if err != nil {
		return fmt.Errorf("failed sending message to producer: %s", err)
	}
	defer p.Close()
	return nil
}

func (s *TransactionService) Transfer(id int, sum models.InputTransfer) error {
	p := deps.NewProducer(viper.GetString("kafka.brokers"))
	message, err := json.Marshal(sum)
	if err != nil {
		return fmt.Errorf("failed sending message to producer: %s", err)
	}
	err = p.SendMessage([]byte(message), constants.Transfer)
	if err != nil {
		return fmt.Errorf("failed sending message to producer: %s", err)
	}
	defer p.Close()
	return nil
}
