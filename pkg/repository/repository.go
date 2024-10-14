package repository

import (
	"bankingsystem/models"
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
)

type Accounts interface {
	CreateAccount(input models.CreateAccount) (int, error)
	DeleteAccount(id int) error
	GetAccountById(id int) (models.GetAccount, error)
	GetAccounts() ([]models.GetAccount, error)
}

type Transactions interface {
	AddDeposit(deposit models.InputDeposit) (int, error)
	Withdraw(deposit models.InputWithdraw) (int, error)
	Transfer(deposit models.InputTransfer) (int, error)

	IsExistAccount(id int) error
	IsLockedAccount(id int) error
	CheckBalance(id int) (float64, error)
}

type Repository struct {
	Accounts
	Transactions
}

func NewRepository(db *pgxpool.Pool, redis *redis.Client, ctx context.Context) *Repository {
	return &Repository{
		Accounts:     NewAccountRepository(db, redis),
		Transactions: NewTransactionRepository(db, redis),
	}
}
