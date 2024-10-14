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

	err := t.IsExistAccount(deposit.Id)
	if err != nil {
		return 0, err
	}

	err = t.IsLockedAccount(deposit.Id)
	if err != nil {
		return 0, err
	}

	var id int
	query := fmt.Sprintf("INSERT INTO %s (account_id, amount, transaction_type) VALUES ($1, $2, $3) RETURNING id", constants.TransactionsTable)
	row := t.db.QueryRow(ctx, query, deposit.Id, deposit.DepositSum, constants.Deposit)

	if err := row.Scan(&id); err != nil {
		return 0, err
	}

	return id, nil
}

func (t *TransactionRepository) Withdraw(withdraw models.InputWithdraw) (int, error) {

	err := t.IsExistAccount(withdraw.Id)
	if err != nil {
		return 0, err
	}

	err = t.IsLockedAccount(withdraw.Id)
	if err != nil {
		return 0, err
	}

	balance, err := t.CheckBalance(withdraw.Id)
	if err != nil {
		return 0, nil
	}

	if balance < withdraw.WithDrawSum {
		return 0, fmt.Errorf("there are not enough funds in the account")
	}
	newBalance := balance - withdraw.WithDrawSum

	tx, err := t.db.Begin(ctx)
	if err != nil {
		return 0, err
	}
	query := fmt.Sprintf("update %s set balance = $1 where id=$2", constants.AccountsTable)
	_, err = tx.Exec(ctx, query, newBalance, withdraw.Id)
	if err != nil {
		tx.Rollback(ctx)
		return 0, err
	}

	queryTransaction := fmt.Sprintf("insert into %s (account_id, amount, transaction_type) values ($1, $2, $3) returning id", constants.TransactionsTable)
	var transactionID int
	row := tx.QueryRow(ctx, queryTransaction, withdraw.Id, withdraw.WithDrawSum, constants.Withdraw)
	if err := row.Scan(&transactionID); err != nil {
		tx.Rollback(ctx)
		return 0, fmt.Errorf("failed to create transaction record: %v", err)
	}

	return transactionID, nil
}

func (t *TransactionRepository) Transfer(transfer models.InputTransfer) (int, error) {
	// ctx := context.Background()
	// _, err := d.db.Exec(ctx, "INSERT INTO transactions (account_id, amount, transaction_type, created_at) VALUES ($1, $2, 'deposit', NOW())",
	// 	deposit.Id, deposit.DepositSum)
	// if err != nil {
	return 0, nil
}

func (t *TransactionRepository) IsExistAccount(id int) error {

	var exists bool
	queryCheck := fmt.Sprintf("SELECT EXISTS(SELECT 1 FROM %s WHERE id=$1 AND deleted_at IS NULL)", constants.AccountsTable)
	row := t.db.QueryRow(ctx, queryCheck, id)

	if err := row.Scan(&exists); err != nil {
		return fmt.Errorf("failed to check account existence: %v", err)
	}

	if !exists {
		return fmt.Errorf("account with id %d does not exist", id)
	}
	return nil
}

func (t *TransactionRepository) IsLockedAccount(id int) error {

	var isLocked bool
	queryCheck := fmt.Sprintf("SELECT EXISTS(SELECT 1 FROM %s WHERE id=$1 AND is_locked=false)", constants.AccountsTable)
	row := t.db.QueryRow(ctx, queryCheck, id)

	if err := row.Scan(&isLocked); err != nil {
		return fmt.Errorf("failed to check account existence: %v", err)
	}

	if !isLocked {
		return fmt.Errorf("account with id %d is locked", id)
	}
	return nil
}

func (t *TransactionRepository) CheckBalance(id int) (float64, error) {
	var balance float64
	query := fmt.Sprintf("SELECT balance from %s WHERE id=$1 AND deleted_at IS NULL AND is_locked=false", constants.AccountsTable)
	row := t.db.QueryRow(ctx, query, id)
	if err := row.Scan(&balance); err != nil {
		return 0, fmt.Errorf("failed to get balance")
	}

	return balance, nil
}
