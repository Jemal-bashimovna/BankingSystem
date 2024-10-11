package models

import "time"

type CreateAccount struct {
	Balance  float64 `json:"balance" binding:"required"`
	Currency string  `json:"currency" binding:"required"`
}

type GetAccount struct {
	Id        int64     `json:"id" db:"id" `
	Balance   float64   `json:"balance" db:"balance"`
	Currency  string    `json:"currency" db:"currency"`
	IsLocked  bool      `json:"is_locked" db:"is_locked"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}
