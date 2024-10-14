package repository

import (
	"bankingsystem/constants"
	"bankingsystem/models"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
)

type TransactionRepository struct {
	db    *pgxpool.Pool
	redis *redis.Client
}

func NewTransactionRepository(db *pgxpool.Pool, redis *redis.Client) *TransactionRepository {
	return &TransactionRepository{
		db:    db,
		redis: redis,
	}
}

func (t *TransactionRepository) AddDeposit(deposit models.InputDeposit) (int, error) {

	if deposit.DepositSum <= 0 {
		return 0, fmt.Errorf("deposit sum must be positive")
	}

	var exists bool
	queryCheck := fmt.Sprintf("SELECT EXISTS(SELECT 1 FROM %s WHERE id=$1)", constants.AccountsTable)
	row := t.db.QueryRow(ctx, queryCheck, deposit.Id)

	if err := row.Scan(&exists); err != nil {
		return 0, fmt.Errorf("failed to check account existence: %v", err)
	}

	if !exists {
		return 0, fmt.Errorf("account with id %d does not exist", deposit.Id)
	}

	var id int
	query := fmt.Sprintf("INSERT INTO %s (account_id, amount, transaction_type) VALUES ($1, $2) RETURNING id", constants.TransactionsTable)
	row = t.db.QueryRow(ctx, query, deposit.Id, deposit.DepositSum)

	if err := row.Scan(&id); err != nil {
		return 0, err
	}

	return id, nil
}

func (t *TransactionRepository) Withdraw(withdraw models.InputWithdraw) error {
	// ctx := context.Background()
	// _, err := t.db.Exec(ctx, "INSERT INTO transactions (account_id, amount, transaction_type, created_at) VALUES ($1, $2, 'deposit', NOW())",
	// 	deposit.Id, deposit.DepositSum)
	// if err != nil {
	// 	return err
	// }

	return nil
}

func (t *TransactionRepository) Transfer(transfer models.InputTransfer) error {
	// ctx := context.Background()
	// _, err := d.db.Exec(ctx, "INSERT INTO transactions (account_id, amount, transaction_type, created_at) VALUES ($1, $2, 'deposit', NOW())",
	// 	deposit.Id, deposit.DepositSum)
	// if err != nil {
	// 	return err
	// }

	return nil
}

// func (d *DepositConsumer) createDepositTransaction(deposit models.InputDeposit) error {

// }
