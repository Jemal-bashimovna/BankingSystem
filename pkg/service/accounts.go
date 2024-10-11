package service

import (
	"bankingsystem/models"
	"bankingsystem/pkg/repository"
	"errors"
	"fmt"
)

type AccountService struct {
	repo repository.Accounts
}

func NewAccountService(repo repository.Accounts) *AccountService {
	return &AccountService{repo: repo}
}

func (s *AccountService) CreateAccount(input models.CreateAccount) (int, error) {

	if input.Balance < 0 {
		return 0, errors.New("negative Balance")
	}

	accountId, err := s.repo.CreateAccount(input)
	if err != nil {
		return 0, err
	}
	if accountId == 0 {
		return 0, fmt.Errorf("id is empty ")
	}
	return accountId, nil
}

func (s *AccountService) DeleteAccount(id int) error {

	return s.repo.DeleteAccount(id)
}

func (s *AccountService) GetAccountById(id int) (models.GetAccount, error) {

	return s.repo.GetAccountById(id)
}

func (s *AccountService) GetAccounts() ([]models.GetAccount, error) {

	return s.repo.GetAccounts()
}
