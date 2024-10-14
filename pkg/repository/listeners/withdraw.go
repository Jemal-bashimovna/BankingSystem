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

		// err = d.createDepositTransaction(transaction)
		// if err != nil {
		// 	log.Printf("Failed to create deposit transaction: %v", err)
		// } else {
		// 	log.Println("Deposit transaction successfully created")
		// }

		fmt.Println(transaction)
	}
}
