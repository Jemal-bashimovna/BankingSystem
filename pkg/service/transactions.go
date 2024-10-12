package service

import (
	"bankingsystem/models"
	"bankingsystem/pkg/repository"
	"encoding/json"
	"fmt"

	"github.com/spf13/viper"
)

const (
	deposit  = "account-deposit"
	withdraw = "account-withdraw"
	transfer = "account-transfer"
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
	p := repository.NewProducer(viper.GetString("kafka.brokers"))
	message, err := json.Marshal(sum)
	if err != nil {
		return fmt.Errorf("failed sending message to producer: %s", err)
	}
	err = p.SendMessage([]byte(message), deposit)

	if err != nil {
		return fmt.Errorf("failed sending message to producer: %s", err)
	}
	defer p.Close()
	return nil
}

func (s *TransactionService) Withdraw(id int, sum models.InputWithdraw) error {
	p := repository.NewProducer(viper.GetString("kafka.brokers"))
	message, err := json.Marshal(sum)
	if err != nil {
		return fmt.Errorf("failed sending message to producer: %s", err)
	}
	err = p.SendMessage([]byte(message), withdraw)
	if err != nil {
		return fmt.Errorf("failed sending message to producer: %s", err)
	}
	defer p.Close()
	return nil
}

func (s *TransactionService) Transfer(id int, sum models.InputTransfer) error {
	p := repository.NewProducer(viper.GetString("kafka.brokers"))
	message, err := json.Marshal(sum)
	if err != nil {
		return fmt.Errorf("failed sending message to producer: %s", err)
	}
	err = p.SendMessage([]byte(message), transfer)
	if err != nil {
		return fmt.Errorf("failed sending message to producer: %s", err)
	}
	defer p.Close()
	return nil
}
