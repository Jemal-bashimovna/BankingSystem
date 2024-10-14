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
	db       *repository.Repository // repository.TransactionsRepo
}

func NewDepositConsumer(brokers, groupId string, depostsConsumer *deps.Consumer, db *repository.Repository) *DepositConsumer {
	deps.NewConsumer(brokers, groupId, []string{constants.Deposit})
	return &DepositConsumer{consumer: depostsConsumer, db: db}
}

func (d *DepositConsumer) StartListening() {
	for {
		msg, err := d.consumer.PollMessage()
		if err != nil {
			log.Printf("Failed to read message: %v", err)
			continue
		}
		var transaction models.InputDeposit
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

// func (d *DepositConsumer) createDepositTransaction(deposit models.InputDeposit) error {
// 	ctx := context.Background()
// 	_, err := d.db.Exec(ctx, "INSERT INTO transactions (account_id, amount, transaction_type, created_at) VALUES ($1, $2, 'deposit', NOW())",
// 		deposit.Id, deposit.DepositSum)
// 	if err != nil {
// 		return err
// 	}
// 	return nil
// }
