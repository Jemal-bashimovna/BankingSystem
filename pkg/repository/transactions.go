package repository

import (
	"bankingsystem/constants"
	"bankingsystem/models"
	"fmt"
	"strconv"
	"time"

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

	balance, err := t.CheckBalance(deposit.Id)
	if err != nil {
		return 0, err
	}

	tx, err := t.db.Begin(ctx)
	if err != nil {
		return 0, fmt.Errorf("failed to begin transaction: %v", err)
	}
	newBalance := balance + deposit.DepositSum

	query := fmt.Sprintf("update %s set balance=$1 where id=$2", constants.AccountsTable)
	_, err = tx.Exec(ctx, query, newBalance, deposit.Id)
	if err != nil {
		tx.Rollback(ctx)
		return 0, err
	}

	var transactionId int
	transactionQuery := fmt.Sprintf("INSERT INTO %s (account_id, amount, transaction_type) VALUES ($1, $2, $3) RETURNING id", constants.TransactionsTable)
	row := tx.QueryRow(ctx, transactionQuery, deposit.Id, deposit.DepositSum, constants.Deposit)

	if err := row.Scan(&transactionId); err != nil {
		return 0, err
	}

	if err := tx.Commit(ctx); err != nil {
		return 0, fmt.Errorf("failed to commit transaction: %v", err)
	}
	// add cache to redis

	cacheData := map[string]interface{}{
		"id":     transactionId,
		"amount": deposit.DepositSum,
		"type":   constants.Deposit,
		"date":   time.Now(),
	}

	cacheKey := fmt.Sprintf("transaction:%d:%d", deposit.Id, transactionId)

	err = t.redis.HSet(ctx, cacheKey, cacheData).Err()
	if err != nil {
		return 0, fmt.Errorf("failed to cache transfer data: %v", err)
	}

	err = t.redis.Expire(ctx, cacheKey, 24*time.Hour).Err()
	if err != nil {
		return 0, fmt.Errorf("failed to set cache expiration: %v", err)
	}

	return transactionId, nil
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
		return 0, err
	}

	if withdraw.WithDrawSum > balance {
		return 0, fmt.Errorf("there are not enough funds in the account")
	}

	tx, err := t.db.Begin(ctx)
	if err != nil {
		return 0, fmt.Errorf("failed to begin transaction: %v", err)
	}
	newBalance := balance - withdraw.WithDrawSum

	query := fmt.Sprintf("update %s set balance = $1 where id=$2", constants.AccountsTable)
	_, err = tx.Exec(ctx, query, newBalance, withdraw.Id)
	if err != nil {
		tx.Rollback(ctx)
		return 0, err
	}

	queryTransaction := fmt.Sprintf("insert into %s (account_id, amount, transaction_type) values ($1, $2, $3) returning id", constants.TransactionsTable)
	var transactionId int
	row := tx.QueryRow(ctx, queryTransaction, withdraw.Id, withdraw.WithDrawSum, constants.Withdraw)
	if err := row.Scan(&transactionId); err != nil {
		tx.Rollback(ctx)
		return 0, fmt.Errorf("failed to create transaction record: %v", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return 0, fmt.Errorf("failed to commit transaction: %v", err)
	}

	cacheData := map[string]interface{}{
		"id":     transactionId,
		"amount": withdraw.WithDrawSum,
		"type":   constants.Withdraw,
		"date":   time.Now(),
	}
	cacheKey := fmt.Sprintf("transaction:%d:%d", withdraw.Id, transactionId)
	err = t.redis.HSet(ctx, cacheKey, cacheData).Err()
	if err != nil {
		return 0, fmt.Errorf("failed to cache transfer data: %v", err)
	}

	err = t.redis.Expire(ctx, cacheKey, 24*time.Hour).Err()
	if err != nil {
		return 0, fmt.Errorf("failed to set cache expiration: %v", err)
	}

	return transactionId, nil
}

func (t *TransactionRepository) Transfer(transfer models.InputTransfer) (int, error) {
	// chek account
	err := t.IsExistAccount(transfer.Id)
	if err != nil {
		return 0, err
	}

	err = t.IsLockedAccount(transfer.Id)
	if err != nil {
		return 0, err
	}

	balance, err := t.CheckBalance(transfer.Id)
	if err != nil {
		return 0, err
	}

	if transfer.TransferSum > balance {
		return 0, fmt.Errorf("there are not enough funds in the account")
	}

	// check target account
	err = t.IsExistAccount(transfer.TargetId)
	if err != nil {
		return 0, err
	}

	err = t.IsLockedAccount(transfer.TargetId)
	if err != nil {
		return 0, err
	}

	targetBalance, err := t.CheckBalance(transfer.Id)
	if err != nil {
		return 0, err
	}

	tx, err := t.db.Begin(ctx)
	if err != nil {
		return 0, err
	}
	newBalance := balance - transfer.TransferSum
	newTargetBalance := targetBalance + transfer.TransferSum

	query := fmt.Sprintf("update %s set balance=$1 where id=$2", constants.AccountsTable)
	_, err = tx.Exec(ctx, query, newBalance, transfer.Id)
	if err != nil {
		tx.Rollback(ctx)
		return 0, err
	}

	_, err = tx.Exec(ctx, query, newTargetBalance, transfer.TargetId)
	if err != nil {
		tx.Rollback(ctx)
		return 0, err
	}

	var transactionId int
	transactionQuery := fmt.Sprintf("insert into %s (account_id, amount, transaction_type) values ($1, $2, $3) returning id", constants.TransactionsTable)
	row := tx.QueryRow(ctx, transactionQuery, transfer.Id, transfer.TransferSum, constants.Transfer)
	if err := row.Scan(&transactionId); err != nil {
		tx.Rollback(ctx)
		return 0, err
	}
	if err := tx.Commit(ctx); err != nil {
		return 0, fmt.Errorf("failed to commit transaction: %v", err)
	}

	cacheData := map[string]interface{}{
		"id":     transactionId,
		"amount": transfer.TransferSum,
		"type":   constants.Transfer,
		"date":   time.Now(),
	}

	cacheKey := fmt.Sprintf("transaction:%d:%d", transfer.Id, transactionId)

	err = t.redis.HSet(ctx, cacheKey, cacheData).Err()
	if err != nil {
		return 0, fmt.Errorf("failed to cache transfer data: %v", err)
	}

	err = t.redis.Expire(ctx, cacheKey, 24*time.Hour).Err()
	if err != nil {
		return 0, fmt.Errorf("failed to set cache expiration: %v", err)
	}

	return transactionId, nil
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

func (r *TransactionRepository) GetAll(id int) ([]models.GetTransactions, error) {
	var cursor uint64
	var transactions []models.GetTransactions
	cacheKey := fmt.Sprintf("transaction:%d:*", id)

	for {
		keys, newCursor, err := r.redis.Scan(ctx, cursor, cacheKey, 0).Result()
		if err != nil {
			return nil, fmt.Errorf("failed to scan transaction keys from cache: %s", err)
		}
		for _, key := range keys {

			data, err := r.redis.HGetAll(ctx, key).Result()
			if err != nil {
				return nil, fmt.Errorf("failed to get transaction data from cache: %s", err)
			}

			var transaction models.GetTransactions

			// convert id to int
			if idStr, ok := data["id"]; ok {
				id, err := strconv.Atoi(idStr)
				if err != nil {
					return nil, fmt.Errorf("failed to convert id to int: %s", err)
				}
				transaction.Id = id
			}

			// convert amount to float64
			if amountStr, ok := data["amount"]; ok {
				amount, err := strconv.ParseFloat(amountStr, 64)
				if err != nil {
					return nil, fmt.Errorf("failed to convert amount to float64: %s", err)
				}
				transaction.Amount = amount
			}

			// convert date to time.Time
			if dateStr, ok := data["date"]; ok {
				transaction.CreatedAt, err = time.Parse(time.RFC3339, dateStr)
				if err != nil {
					return nil, fmt.Errorf("failed to parse date: %v", err)
				}
			}

			if typeStr, ok := data["type"]; ok {
				transaction.TransactionType = typeStr
			}
			transactions = append(transactions, transaction)
		}
		if newCursor == 0 {
			break
		}
		cursor = newCursor
	}

	query := fmt.Sprintf("select id, amount, transaction_type, created_at from %s where account_id=$1", constants.TransactionsTable)
	rows, err := r.db.Query(ctx, query, id)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var transaction models.GetTransactions
		if err := rows.Scan(&transaction.Id, &transaction.Amount, &transaction.TransactionType, &transaction.CreatedAt); err != nil {
			return nil, fmt.Errorf("error scanning row: %w", err)
		}
		transactions = append(transactions, transaction)
	}
	return transactions, nil
}
