package notifier

import (
	"encoding/json"
	"fmt"
	"log"
	"net/smtp"
	"os"
)

var (
	smtpHost     string = os.Getenv("SMTP_HOST")
	smtpPort     string = os.Getenv("SMTP_PORT")
	smtpUsername string = os.Getenv("SMTP_USERNAME")
	smtpPassword string = os.Getenv("SMTP_PASSWORD")
	smtpSender   string = os.Getenv("SMTP_SENDER")
)

func MessageUnmarshal(messageBody []byte) error {
	var event NotificationEvent
	err := json.Unmarshal(messageBody, &event)
	if err != nil {
		log.Print("Error unmarhalling message : %v", err)
		return err
	}
	err = ProcessMessage(event)
	if err != nil {
		log.Print("Error  message : %v", err)
		return err
	}
	return nil
}

func ProcessMessage(event NotificationEvent) error {
	var err error
	log.Println(event)

	log.Println("Processing")
	i := 0
	for i < len(event.Channels) {

		switch event.Channels[i].Type {
		case "email":
			if smtpHost == "" || smtpPassword == "" || smtpPort == "" || smtpUsername == "" {
				return fmt.Errorf("SMTP not configured")
			}

			err = SendEmailSMTP(event.Channels[i].Contact, smtpSender, "Notification", event.NotificationMessage)

		case "whatsapp":

		default:

		}
		if err != nil {
			log.Println("Error occured %v", err)
			return err
		}

		i++
	}

	return nil
}

func SendEmailSMTP(toEmail, smtpSender, subject, body string) error {

	auth := smtp.PlainAuth("", smtpUsername, smtpPassword, smtpHost)
	address := smtpHost + ":" + smtpPort
	msg := []byte("To: " + toEmail + "\r\n" +
		"From: " + smtpSender + "\r\n" +
		"Subject: " + subject + "\r\n" +
		"Content-Type: text/html; charset=UTF-8 \r\n" +
		"\r\n" +
		body + "\r\n")
	err := smtp.SendMail(address, auth, smtpSender, []string{toEmail}, msg)

	if err != nil {
		return fmt.Errorf("SMTP send mail failed: %w", err)
	}

	log.Println("Email triggered successfully to %s", toEmail)

	return nil
}
