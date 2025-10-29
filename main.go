package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"
	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azservicebus"
	"github.com/joho/godotenv"
)

type TransactionCompletedEvent struct{
	UserID string `json:"userId"`
	TransactionID int64 `json: "transactionId"`
	TransactionType string `json: "transactionType"`
	Amount float64 `json: "amount"`
	Timestamp string `json: "timestamp"`
}

func main()  {
	err := godotenv.Load()
	if err != nil {
		log.Printf(" Could not load .env file")
	}
	connectionstring := os.Getenv("SERVICEBUS_CONNECTION_STRING")
	queueName := os.Getenv("SERVICEBUS_QUEUE_NAME")

	if connectionstring == "" || queueName == "" {
		log.Fatal("service bus env variables must be set")
	}

	client, err := azservicebus.NewClientFromConnectionString(connectionstring,nil)

	if err != nil {
		log.Fatal("Failed to create service bus client %v",err)
	}

	defer client.Close(context.Background())

	receiver, err := client.NewReceiverForQueue(queueName,&azservicebus.ReceiverOptions{
		ReceiveMode: azservicebus.ReceiveModePeekLock})

	if err != nil {
		log.Fatalf("Failed to create receiver for queue %s:%v",queueName,nil)
	}

	defer receiver.Close(context.Background())

	fmt.Printf("Notification service started")

	for {
		ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
		messages,err := receiver.ReceiveMessages(ctx,1,nil)
		cancel()

		if err != nil {
			log.Printf("Error receiving message: %v. Retrying...\n", err)
			time.Sleep(5 * time.Second) // Wait before retrying
			continue
		}

		if len(messages) == 0 {
			continue
		}

		msg := messages[0]
		log.Printf("Received message ID: %s\n", msg.MessageID)
		log.Printf("Message body: %s\n", string(msg.Body))

		err = receiver.CompleteMessage(context.Background(),msg,nil)

		if err!=nil{
			log.Printf("Error completing message %s: %v\n", msg.MessageID, err)
		} else {
			log.Printf("Message %s completed successfully.\n", msg.MessageID)
		}
	}



}