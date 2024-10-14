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
	AddDeposit(deposit models.InputDeposit) error
	Withdraw(deposit models.InputDeposit) error
	Transfer(deposit models.InputDeposit) error
}

type Repository struct {
	Accounts
	Transactions
}

func NewRepository(db *pgxpool.Pool, redis *redis.Client, ctx context.Context) *Repository {
	return &Repository{
		Accounts: NewAccountRepository(db, redis),
	}
}
