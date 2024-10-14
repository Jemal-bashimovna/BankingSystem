package models

type CreateAccResponse struct {
	Status string `json:"status"`
	Id     int    `json:"id"`
}

type DeleteAccResponse struct {
	Status string `json:"status"`
}

type GetAccResponse struct {
	Account GetAccount `json:"account"`
}

type GetAllAccResponse struct {
	Accounts []GetAccount `json:"accounts"`
}
