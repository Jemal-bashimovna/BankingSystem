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
	db       *repository.Repository // repository.TransactionsRepo
}

func NewTransferConsumer(brokers, groupId string, depostsConsumer *deps.Consumer, db *repository.Repository) *TransferConsumer {
	deps.NewConsumer(brokers, groupId, []string{constants.Transfer})
	return &TransferConsumer{consumer: depostsConsumer, db: db}
}

func (d *TransferConsumer) StartListening() {
	for {
		msg, err := d.consumer.PollMessage()
		if err != nil {
			log.Printf("Failed to read message: %v", err)
			continue
		}
		var transaction models.InputTransfer
		if err := json.Unmarshal(msg, &transaction); err != nil {
			log.Printf("Failed to unmarshal message: %v", err)
			continue
		}

		// err = d.createDepositTransaction(transaction)
		// if err != nil {
		// 	log.Printf("Failed to create deposit transaction: %v", err)
		// } else {
		// 	log.Println("Deposit transaction successfully created")
		// }
		d.db.Transactions.AddDeposit()
		fmt.Println(transaction)
	}
}
