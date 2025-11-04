package notifier

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/smtp"
	"os"
	"github.com/joho/godotenv"
)

var (
	smtpHost     string
	smtpPort     string
	smtpUsername string
	smtpPassword string
	smtpSender   string

	acs_app_id string
	acs_app_secret string
)

type loginAuth struct {
	username, password string
}

func LoginAuth(username, password string) smtp.Auth {
	return &loginAuth{username, password}
}

func (a *loginAuth) Start(server *smtp.ServerInfo) (string, []byte, error) {
	return "LOGIN", []byte{}, nil
}

func (a *loginAuth) Next(fromServer []byte, more bool) ([]byte, error) {
	if more {
		switch string(fromServer) {
		case "Username:":
			return []byte(a.username), nil
		case "Password:":
			return []byte(a.password), nil
		default:
			return nil, errors.New("unknown fromserver")
		}
	}

	return nil, nil
}

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
	err = godotenv.Load()
	if err != nil {
		log.Printf(" Could not load .env file")
	}
	smtpHost = os.Getenv("SMTP_HOST")
	smtpPort = os.Getenv("SMTP_PORT")
	smtpUsername = os.Getenv("SMTP_USERNAME")
	smtpPassword = os.Getenv("SMTP_PASSWORD")
	smtpSender = os.Getenv("SMTP_SENDER")

	acs_app_id = os.Getenv("ACS_APP_ID")
	acs_app_secret =os.Getenv("ACS_APP_SECRET")

	log.Println("Processing")
	i := 0
	for i < len(event.Channels) {

		switch event.Channels[i].Type {
		case "email":
			if smtpHost == "" || smtpPassword == "" || smtpPort == "" || smtpUsername == "" {
				log.Println(smtpHost, smtpPort, "\n", smtpUsername, "\n", smtpPassword)
				return fmt.Errorf("SMTP not configured")
			}

			err = SendEmailSMTP(event.Channels[i].Contact, smtpSender, "Notification", event.NotificationMessage)

		case "whatsapp":
			if acs_app_id == "" || acs_app_secret == ""{
				return fmt.Errorf("ACS WhatsApp parameters not configured")
			}


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

	address := smtpHost + ":" + smtpPort
	// auth := smtp.PlainAuth("", smtpUsername, smtpPassword, smtpHost)
	auth := LoginAuth(smtpUsername, smtpPassword)

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


func sendWhatsAppMessage(toNumber,fromNumber, body string){
	
}