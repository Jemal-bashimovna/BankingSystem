package listeners

import (
	"bankingsystem/constants"
	"bankingsystem/deps"
	"bankingsystem/models"
	"bankingsystem/pkg/repository"
	"encoding/json"
	"fmt"
	"log"
)

type DepositConsumer struct {
	consumer *deps.Consumer
	db       repository.Transactions // repository.TransactionsRepo

}

func NewDepositConsumer(brokers, groupId string, db repository.Transactions) *DepositConsumer {
	depositConsumer := deps.NewConsumer(brokers, groupId, []string{constants.Deposit})
	return &DepositConsumer{
		consumer: depositConsumer,
		db:       db,
	}
}

func (d *DepositConsumer) StartListening() {
	for {
		msg, err := d.consumer.PollMessage()
		if err != nil {
			log.Printf("Failed to read message: %s", err)
			continue
		}

		var transaction models.InputDeposit
		if err := json.Unmarshal(msg.Value, &transaction); err != nil {
			log.Printf("Failed to unmarshal message: %s", err)
			continue
		}
		if err := validateDeposit(transaction); err != nil {
			log.Printf("Validation failed: %s", err)
			continue
		}

		id, err := d.db.AddDeposit(transaction)
		if err != nil {
			log.Printf("Failed to add deposit: %s", err)
			return
		}
		fmt.Println(transaction)
		fmt.Printf("Deposit successfully id: %d", id)
	}
}

func validateDeposit(deposit models.InputDeposit) error {
	if deposit.Id <= 0 {
		return fmt.Errorf("invalid account id")
	}
	if deposit.DepositSum <= 0 {
		return fmt.Errorf("deposit amount must be positive")
	}
	return nil
}
