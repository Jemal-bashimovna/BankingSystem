package models

import "time"

type InputDeposit struct {
	AccountId  int     `json:"account_id"`
	DepositSum float64 `json:"sum"`
}

type InputWithdraw struct {
	AccountId   int     `json:"account_id"`
	WithDrawSum float64 `json:"sum"`
}

type InputTransfer struct {
	AccountId   int     `json:"account_id"`
	TargetId    int     `json:"target_id"`
	TransferSum float64 `json:"sum"`
}

type TransactionResponse struct {
	Message string `json:"message"`
}

type GetTransactions struct {
	Id              int       `json:"id"`
	Amount          float64   `json:"amount"`
	TransactionType string    `json:"transaction_type"`
	CreatedAt       time.Time `json:"created_at"`
}

type GetTransactionsResponse struct {
	Transactions []GetTransactions `json:"transactions"`
}
