package repository

import (
	"bankingsystem/models"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
)

type Accounts interface {
	CreateAccount(input models.CreateAccount) (int, error)
	DeleteAccount(id int) error
	GetAccountById(id int) (models.GetAccount, error)
	GetAccounts() ([]models.GetAccount, error)
}

type Transactions interface{}

type Repository struct {
	Accounts
	Transactions
}

func NewRepository(db *pgxpool.Pool, redis *redis.Client) *Repository {
	return &Repository{
		Accounts: NewAccountRepository(db, redis),
	}
}
