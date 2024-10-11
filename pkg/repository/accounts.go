package repository

import (
	"bankingsystem/models"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
)

type AccountRepository struct {
	db    *pgxpool.Pool
	redis *redis.Client
}

func NewAccountRepository(db *pgxpool.Pool, redis *redis.Client) *AccountRepository {
	return &AccountRepository{
		db:    db,
		redis: redis,
	}
}

func (r *AccountRepository) CreateAccount(input models.CreateAccount) (int, error) {

	var accountId int
	var balance float64
	query := fmt.Sprintf("INSERT INTO %s (balance, currency) VALUES ($1, $2) RETURNING id, balance", accountsTable)
	row := r.db.QueryRow(ctx, query, input.Balance, input.Currency)
	if err := row.Scan(&accountId, &balance); err != nil {
		return 0, err
	}

	// Создание ключа для Redis
	key := fmt.Sprintf("account:%d", accountId)

	// Кэширование данных аккаунта в Redis
	err := r.redis.HSet(ctx, key, map[string]interface{}{
		"id":      accountId,
		"balance": balance,
	}).Err()

	if err != nil {
		return accountId, fmt.Errorf("failed to cache account data in Redis: %s", err)
	}

	err = r.redis.Expire(ctx, key, 24*time.Hour).Err()
	if err != nil {
		return accountId, fmt.Errorf("failed to expiration for cache: %s", err)
	}

	return accountId, nil
}

func (r *AccountRepository) DeleteAccount(id int) error {

	query := fmt.Sprintf(`
    UPDATE %s SET deleted_at=NOW() WHERE id=$1`, accountsTable)
	result, err := r.db.Exec(ctx, query, id)
	if err != nil {
		return err
	}

	updated := result.RowsAffected()

	if updated == 0 {
		return fmt.Errorf("account not found with this id")
	}

	return nil
}

func (r *AccountRepository) GetAccountById(id int) (models.GetAccount, error) {
	var account models.GetAccount

	// // get from redis
	// key := fmt.Sprintf("account:%d", id)
	// acc, err := r.redis.HGetAll(ctx, key).Result()
	// if err == nil && len(acc) > 0 {
	// 	account.Id, _ = strconv.ParseInt(acc["id"], 10, 64)
	// 	account.Balance, _ = strconv.ParseFloat(acc["balance"], 64)
	// 	return account, nil
	// } else if err != nil && err != redis.Nil {
	// 	return account, fmt.Errorf("error getting from cache: %w", err)
	// }

	// get from postgres db
	query := fmt.Sprintf("SELECT id, balance, currency, is_locked, created_at FROM %s WHERE id = $1 AND deleted_at IS NULL", accountsTable)
	row := r.db.QueryRow(ctx, query, id)
	err := row.Scan(&account.Id, &account.Balance, &account.Currency, &account.IsLocked, &account.CreatedAt)
	if err != nil {
		if err.Error() == "no rows in result set" {
			return account, fmt.Errorf("account not found with id: %d", id)
		}
		return account, fmt.Errorf("error getting from database: %s", err)
	}

	// set to redis
	key := fmt.Sprintf("account:%d", id)
	err = r.redis.HSet(ctx, key, map[string]interface{}{
		"id":      account.Id,
		"balance": account.Balance,
	}).Err()
	if err != nil {
		return account, fmt.Errorf("failed to cache account data in Redis: %s", err)
	}

	err = r.redis.Expire(ctx, key, 24*time.Hour).Err()
	if err != nil {
		return account, fmt.Errorf("failed to set expiration for cache: %s", err)
	}

	return account, nil
}

func (r *AccountRepository) GetAccounts() ([]models.GetAccount, error) {
	var accounts []models.GetAccount

	query := fmt.Sprintf("SELECT id, balance, currency, is_locked, created_at FROM %s WHERE deleted_at IS NULL", accountsTable)
	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return accounts, fmt.Errorf("error getting from database: %s", err)
	}

	defer rows.Close()

	for rows.Next() {
		var account models.GetAccount
		if err := rows.Scan(&account.Id, &account.Balance, &account.Currency, &account.IsLocked, &account.CreatedAt); err != nil {
			return nil, fmt.Errorf("error scanning row: %w", err)
		}
		accounts = append(accounts, account)
	}

	return accounts, nil
}
