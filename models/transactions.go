package models

type InputDeposit struct {
	Id         int     `json:"id"`
	DepositSum float64 `json:"sum"`
}

type InputWithdraw struct {
	Id          int     `json:"id"`
	WithDrawSum float64 `json:"sum"`
}

type InputTransfer struct {
	Id          int     `json:"id"`
	TargetId    int     `json:"target_id"`
	TransferSum float64 `json:"sum"`
}
