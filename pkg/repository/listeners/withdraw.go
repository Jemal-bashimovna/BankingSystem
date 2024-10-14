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
	db       *repository.Repository // repository.TransactionsRepo
}

func NewWithdrawConsumer(brokers, groupId string, depostsConsumer *deps.Consumer, db *repository.Repository) *WithdrawConsumer {
	deps.NewConsumer(brokers, groupId, []string{constants.Withdraw})
	return &WithdrawConsumer{consumer: depostsConsumer, db: db}
}

func (d *WithdrawConsumer) StartListening() {
	for {
		msg, err := d.consumer.PollMessage()
		if err != nil {
			log.Printf("Failed to read message: %v", err)
			continue
		}
		var transaction models.InputWithdraw
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
