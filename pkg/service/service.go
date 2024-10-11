package service

import (
	"bankingsystem/models"
	"bankingsystem/pkg/repository"
	kafkaproducer "bankingsystem/pkg/repository/kafka-producer"
)

type Accounts interface {
	CreateAccount(inputAccount models.CreateAccount) (int, error)
	DeleteAccount(id int) error
	GetAccountById(id int) (models.GetAccount, error)
	GetAccounts() ([]models.GetAccount, error)
}

type Transactions interface {
	Deposit(id int) error
}

type Service struct {
	Accounts
	Transactions
}

func NewService(repo *repository.Repository, producer *kafkaproducer.Producer) *Service {
	return &Service{
		Accounts:     NewAccountService(repo.Accounts),
		Transactions: NewTransactionService(repo, producer),
	}
}
