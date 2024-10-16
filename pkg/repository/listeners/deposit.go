package listeners

import (
	"bankingsystem/deps"
	"bankingsystem/models"
	"bankingsystem/pkg/repository"
	"encoding/json"
	"fmt"
	"log"
)

type DepositConsumer struct {
	consumer *deps.Consumer
	repo     repository.Transactions // repository.TransactionsRepo

}

func NewDepositConsumer(depositConsumer *deps.Consumer, groupId string, repo repository.Transactions) *DepositConsumer {
	return &DepositConsumer{
		consumer: depositConsumer,
		repo:     repo,
	}
}

func (d *DepositConsumer) StartListening() {
	for {
		msg, err := d.consumer.PollMessage()
		if err != nil {
			log.Printf("Failed to read message: %s", err)
			break
		}

		var transaction models.InputDeposit
		if err := json.Unmarshal(msg.Value, &transaction); err != nil {
			log.Printf("Failed to unmarshal message: %s", err)
			break
		}
		if err := validateDeposit(transaction); err != nil {
			log.Printf("Validation failed: %s", err)
			continue
		}

		id, err := d.repo.AddDeposit(transaction)
		if err != nil {
			log.Fatalf("Failed to add deposit: %s", err)
		}

		log.Printf("Deposit to account: %d successfully id: %d", transaction.AccountId, id)
	}
}

func validateDeposit(deposit models.InputDeposit) error {
	if deposit.AccountId <= 0 {
		return fmt.Errorf("invalid account id")
	}
	if deposit.DepositSum <= 0 {
		return fmt.Errorf("deposit amount must be positive")
	}
	return nil
}
