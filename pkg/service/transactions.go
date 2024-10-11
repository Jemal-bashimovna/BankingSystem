package service

import (
	"bankingsystem/pkg/repository"
	kafkaproducer "bankingsystem/pkg/repository/kafka-producer"
)

type TransactionService struct {
	repo     repository.Transactions
	producer *kafkaproducer.Producer
}

func NewTransactionService(repo repository.Transactions, producer *kafkaproducer.Producer) *TransactionService {
	return &TransactionService{
		repo:     repo,
		producer: producer,
	}
}

func (s *TransactionService) Deposit(id int) error {
	return nil
}
