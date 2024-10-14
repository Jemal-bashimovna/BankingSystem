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
	db       *repository.TransactionRepository // repository.TransactionsRepo

}

func NewDepositConsumer(brokers, groupId string, db *repository.TransactionRepository) *DepositConsumer {
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

		_, err = d.db.AddDeposit(transaction)
		if err != nil {
			log.Printf("Failed to add deposit: %s", err)
		}
		fmt.Println(transaction)
	}
}

// err = d.createDepositTransaction(transaction)
// if err != nil {
// 	log.Printf("Failed to create deposit transaction: %v", err)
// } else {
// 	log.Println("Deposit transaction successfully created")
// }
