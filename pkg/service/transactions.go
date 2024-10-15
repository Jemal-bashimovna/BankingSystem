package service

import (
	"bankingsystem/constants"
	"bankingsystem/deps"
	"bankingsystem/models"
	"bankingsystem/pkg/repository"
	"encoding/json"
	"fmt"
)

type TransactionService struct {
	repo     repository.Transactions
	producer *deps.Producer
}

func NewTransactionService(repo repository.Transactions, producer *deps.Producer) *TransactionService {
	return &TransactionService{
		repo:     repo,
		producer: producer,
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

	message, err := json.Marshal(sum)
	if sum.DepositSum <= 0.00 {
		return fmt.Errorf("deposit sum must be a positive valuse")
	}
	if err != nil {
		return fmt.Errorf("failed sending message to producer: %s", err)
	}
	err = s.producer.SendMessage([]byte(message), constants.Deposit)

	if err != nil {
		return fmt.Errorf("failed sending message to producer: %s", err)
	}
	return nil
}

func (s *TransactionService) WithdrawProducer(id int, sum models.InputWithdraw) error {

	// check account
	err := s.repo.IsExistAccount(sum.Id)
	if err != nil {
		return err
	}

	err = s.repo.IsLockedAccount(sum.Id)
	if err != nil {
		return err
	}

	balance, err := s.repo.CheckBalance(sum.Id)
	if err != nil {
		return err
	}

	if sum.WithDrawSum > balance {
		return fmt.Errorf("there are not enough funds in the account")
	}
	// check account balance

	message, err := json.Marshal(sum)
	if err != nil {
		return fmt.Errorf("failed sending message to producer: %s", err)
	}
	err = s.producer.SendMessage([]byte(message), constants.Withdraw)
	if err != nil {
		return fmt.Errorf("failed sending message to producer: %s", err)
	}
	return nil
}

func (s *TransactionService) TransferProducer(id int, sum models.InputTransfer) error {

	// check account
	err := s.repo.IsExistAccount(sum.Id)
	if err != nil {
		return err
	}

	err = s.repo.IsLockedAccount(sum.Id)
	if err != nil {
		return err
	}

	balance, err := s.repo.CheckBalance(sum.Id)
	if err != nil {
		return err
	}

	if sum.TransferSum > balance {
		return fmt.Errorf("there are not enough funds in the account")
	}

	// check target account
	err = s.repo.IsExistAccount(sum.TargetId)
	if err != nil {
		return err
	}

	err = s.repo.IsLockedAccount(sum.TargetId)
	if err != nil {
		return err
	}

	message, err := json.Marshal(sum)
	if err != nil {
		return fmt.Errorf("failed sending message to producer: %s", err)
	}
	err = s.producer.SendMessage([]byte(message), constants.Transfer)
	if err != nil {
		return fmt.Errorf("failed sending message to producer: %s", err)
	}
	return nil
}

func (s *TransactionService) GetAll(id int) ([]models.GetTransactions, error) {
	return s.repo.GetAll(id)
}
