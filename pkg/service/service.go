package service

import (
	"bankingsystem/deps"
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
	DepositProducer(id int, sum models.InputDeposit) error
	WithdrawProducer(id int, sum models.InputWithdraw) error
	TransferProducer(id int, sum models.InputTransfer) error
	GetAll(id int) ([]models.GetTransactions, error)
}

type Service struct {
	Accounts
	Transactions
}

func NewService(repo *repository.Repository, producer *deps.Producer) *Service {
	return &Service{
		Accounts:     NewAccountService(repo.Accounts),
		Transactions: NewTransactionService(repo, producer),
	}
}
