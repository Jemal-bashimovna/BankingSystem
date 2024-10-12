package service

import (
	"bankingsystem/models"
	"bankingsystem/pkg/repository"
)

type Accounts interface {
	CreateAccount(inputAccount models.CreateAccount) (int, error)
	DeleteAccount(id int) error
	GetAccountById(id int) (models.GetAccount, error)
	GetAccounts() ([]models.GetAccount, error)
}

type Transactions interface {
	Deposit(id int, sum models.InputDeposit) error
	Withdraw(id int, sum models.InputWithdraw) error
	Transfer(id int, sum models.InputTransfer) error
}

type Service struct {
	Accounts
	Transactions
}

func NewService(repo *repository.Repository) *Service {
	return &Service{
		Accounts:     NewAccountService(repo.Accounts),
		Transactions: NewTransactionService(repo),
	}
}
