package processor

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
)

type transactionCompletedEvent struct {
	UserID          string  `json:"userId"`
	TransactionType string  `json:"transactionType"`
	Amount          float64 `json:"amount"`
}

func ProcessMessage(ctx context.Context, messageBody []byte) error {
	log.Printf("Processsing Message: %s\n", string(messageBody))

	var event transactionCompletedEvent
	err:= json.Unmarshal(messageBody, &event)
	if err!= nil{
		log.Printf("Error unmarshalling message: %v\n", err)

		return fmt.Errorf("failed to parse message body: %w", err)
	}

	log.Printf("Fetched the contact details for the UserID %s",event.UserID)

	return err
}