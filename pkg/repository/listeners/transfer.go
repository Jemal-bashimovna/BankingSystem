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

		// err = d.createDepositTransaction(transaction)
		// if err != nil {
		// 	log.Printf("Failed to create deposit transaction: %v", err)
		// } else {
		// 	log.Println("Deposit transaction successfully created")
		// }

		fmt.Println(transaction)
	}
}
