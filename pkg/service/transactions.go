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

func (s *TransactionService) DepositProducer(id int, sum models.InputDeposit) error {

	err := s.repo.IsExistAccount(sum.Id)
	if err != nil {
		return err
	}

	err = s.repo.IsLockedAccount(sum.Id)
	if err != nil {
		return err
	}

	p := deps.NewProducer(viper.GetString("kafka.brokers"))
	message, err := json.Marshal(sum)
	if sum.DepositSum <= 0.00 {
		return fmt.Errorf("deposit sum must be a positive valuse")
	}
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

func (s *TransactionService) WithdrawProducer(id int, sum models.InputWithdraw) error {

	// check account
	err := s.repo.IsExistAccount(sum.Id)
	if err != nil {
		return err
	}

	// check account balance

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

func (s *TransactionService) TransferProducer(id int, sum models.InputTransfer) error {
	err := s.repo.IsExistAccount(sum.Id)
	if err != nil {
		return err
	}

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
