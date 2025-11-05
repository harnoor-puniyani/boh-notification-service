package notifier

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/smtp"
	"net/url"
	"os"

	"github.com/joho/godotenv"
)

var (
	smtpHost     string
	smtpPort     string
	smtpUsername string
	smtpPassword string
	smtpSender   string

	acs_app_id     string
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
	acs_app_secret = os.Getenv("ACS_APP_SECRET")

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
			if acs_app_id == "" || acs_app_secret == "" {
				return fmt.Errorf("ACS WhatsApp parameters not configured")
			}
			err := sendWhatsAppMessage(event.Channels[i].Contact, "abc", event.NotificationMessage)
			if err != nil {
				fmt.Errorf("%v", err)
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

func getOauthToken() (*OauthTokenResponse, error) {

	var response OauthTokenResponse
	tokenUrl := "https://login.microsoftonline.com/0d4e668b-ef12-4a0a-80af-bd78919b7c5a/oauth2/v2.0/token"
	data := map[string]string{
		"grant_type":    "client_credentials",
		"client_id":     acs_app_id,
		"client_secret": acs_app_secret,
		"scope":         "https://communication.azure.com/.default",
	}

	values := url.Values{}

	for key, value := range data {
		values.Add(key, value)
	}

	// jsonData, err := json.Marshal(data)
	// if err != nil {
	// 	return &response, fmt.Errorf("Error occured marshalling the json : %w", err)
	// }
	bodyBuffer := bytes.NewBuffer([]byte(values.Encode()))

	resp, err := http.Post(tokenUrl, "application/x-www-form-urlencoded", bodyBuffer)

	if err != nil {
		return &response, fmt.Errorf("Error Orccured %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		log.Printf("Request Failed \n %v", resp)
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return &response, fmt.Errorf("Error Orccured %v", err)
	}

	log.Printf(string(body))

	err = json.Unmarshal(body, &response)

	if err != nil {
		return &response, fmt.Errorf("Error Orccured %v", err)
	}

	return &response, nil

}

func sendWhatsAppMessage(toNumber, fromNumber, body string) error {

	var event AcsMessage

	if toNumber == "" || fromNumber == "" || body == "" {
		return fmt.Errorf("parameters not configured for whatsApp messaging")
	}

	token, err := getOauthToken()
	if err != nil {
		return fmt.Errorf("Error generating token %v", err)
	}

	event.ChannelRegistrationId = "4eb202d9-3c0d-4594-87a5-1ddef0b9102a"
	event.To = []string{toNumber}
	event.Kind = "text"
	event.Content = body

	data, err := json.Marshal(event)

	if err != nil {
		return err
	}

	bodyBuffer := bytes.NewBuffer(data)

	req, err := http.NewRequest("POST", "https://boh-communication-service.unitedstates.communication.azure.com/messages/notifications:send?api-version=2024-02-01", bodyBuffer)

	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token.AccessToken)

	client := &http.Client{}
	resp, err := client.Do(req)

	if err != nil {
		return err
	}

	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	log.Println("response : ", string(respBody))

	if resp.StatusCode >= http.StatusOK && resp.StatusCode <= http.StatusPermanentRedirect {

		log.Printf("Request Succeeded \n %v", resp)
	} else {
		log.Printf("Request Failed \n %v", resp)
	}

	return nil
}
