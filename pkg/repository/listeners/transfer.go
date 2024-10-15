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

type TransferConsumer struct {
	consumer *deps.Consumer
	db       *repository.TransactionRepository
}

func NewTransferConsumer(brokers, groupId string, db *repository.TransactionRepository) *TransferConsumer {
	transferConsumer := deps.NewConsumer(brokers, groupId, []string{constants.Transfer})
	return &TransferConsumer{
		consumer: transferConsumer,
		db:       db,
	}
}

func (t *TransferConsumer) StartListening() {
	for {
		msg, err := t.consumer.PollMessage()
		if err != nil {
			log.Printf("Failed to read message: %v", err)
			continue
		}

		var transaction models.InputTransfer
		if err := json.Unmarshal(msg.Value, &transaction); err != nil {
			log.Printf("Failed to unmarshal message: %v", err)
			continue
		}

		if err := validateTransfer(transaction); err != nil {
			log.Printf("Validation failed: %s", err)
			continue
		}

		id, err := t.db.Transfer(transaction)
		if err != nil {
			log.Fatalf("Failed to transfer: %v", err)
		}

		log.Printf("The transfer %f from: %d to %d was successfully: %d", transaction.TransferSum, transaction.Id, transaction.TargetId, id)
	}
}

func validateTransfer(transfer models.InputTransfer) error {
	if transfer.Id <= 0 {
		return fmt.Errorf("invalid account id")
	}
	if transfer.TargetId <= 0 {
		return fmt.Errorf("invalid target account id")
	}
	if transfer.TransferSum <= 0 {
		return fmt.Errorf("withdraw amount must be positive")
	}
	return nil
}
