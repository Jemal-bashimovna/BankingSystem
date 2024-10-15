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

type WithdrawConsumer struct {
	consumer *deps.Consumer
	db       *repository.TransactionRepository
}

func NewWithdrawConsumer(brokers, groupId string, db *repository.TransactionRepository) *WithdrawConsumer {
	withdrawConsumer := deps.NewConsumer(brokers, groupId, []string{constants.Withdraw})
	return &WithdrawConsumer{
		consumer: withdrawConsumer,
		db:       db,
	}
}

func (w *WithdrawConsumer) StartListening() {
	for {
		msg, err := w.consumer.PollMessage()
		if err != nil {
			log.Printf("Failed to read message: %v", err)
			continue
		}
		var transaction models.InputWithdraw
		if err := json.Unmarshal(msg.Value, &transaction); err != nil {
			log.Printf("Failed to unmarshal message: %v", err)
			continue
		}

		if err := validateWithdraw(transaction); err != nil {
			log.Printf("Validation failed: %s", err)
			continue
		}

		_, err = w.db.Withdraw(transaction)
		if err != nil {
			log.Fatalf("error withdrawing money from account: %s", err)
		}

		log.Printf("withdrawing money (%f) successfully from account: %d\n", transaction.WithDrawSum, transaction.Id)

	}
}

func validateWithdraw(withdraw models.InputWithdraw) error {
	if withdraw.Id <= 0 {
		return fmt.Errorf("invalid account id")
	}
	if withdraw.WithDrawSum <= 0 {
		return fmt.Errorf("withdraw amount must be positive")
	}
	return nil
}
