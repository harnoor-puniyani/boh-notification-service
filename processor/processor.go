package processor

import (
	"bytes"
	"context"
	// We no longer need: crypto/hmac, crypto/sha256, encoding/base64, net/url, time
	"encoding/json"
	"fmt"
	"log"
	"net/http" // Used for Meta
	"net/smtp" // <-- ADDED: For sending email via SMTP
	"os"
	"strings"
)

// --- Structs to match the JSON message contract ---

type NotificationChannel struct {
	Type    string `json:"type"`
	Contact string `json:"contact"`
}

type NotificationEvent struct {
	UserID              string                `json:"userId"`
	NotificationMessage string                `json:"notificationMessage"`
	Channels            []NotificationChannel `json:"channels"`
}

// --- Client variables ---
var (
	// ACS (Email) variables - Now for SMTP
	smtpHost     string
	smtpPort     string
	smtpUsername string
	smtpPassword string
	smtpSender   string // The "From" email address

	// Meta (WhatsApp) variables
	metaApiToken string
	metaApiUrl   string
)

// Init sets up the clients (call this from main.go)
func Init() {
	// 1. Setup Azure Communication Services (Email via SMTP)
	smtpHost = os.Getenv("SMTP_HOST")         // e.g., "smtp.communication.azure.com"
	smtpPort = os.Getenv("SMTP_PORT")         // e.g., "587"
	smtpUsername = os.Getenv("SMTP_USERNAME") // Your Client ID or full username
	smtpPassword = os.Getenv("SMTP_PASSWORD") // Your Client Secret
	smtpSender = os.Getenv("ACS_SENDER_EMAIL")   // e.g., "donotreply@your-domain.com"

	if smtpHost == "" || smtpPort == "" || smtpUsername == "" || smtpPassword == "" || smtpSender == "" {
		log.Println("Warning: SMTP variables not fully set. Email will be disabled.")
	} else {
		log.Println("Azure Communication Services (SMTP) client configured.")
	}

	// 2. Setup Meta (WhatsApp)
	metaApiToken = os.Getenv("META_API_TOKEN")
	metaApiUrl = os.Getenv("META_API_URL")
	if metaApiToken == "" || metaApiUrl == "" {
		log.Println("Warning: META_API_TOKEN or META_API_URL not set. WhatsApp will be disabled.")
	} else {
		log.Println("Meta (WhatsApp) API configured.")
	}
}

// ProcessMessage handles a single message received from Service Bus
func ProcessMessage(ctx context.Context, messageBody []byte) error {
	log.Printf("Processing message: %s\n", string(messageBody))

	var event NotificationEvent
	if err := json.Unmarshal(messageBody, &event); err != nil {
		log.Printf("Error unmarshalling message: %v\n", err)
		return fmt.Errorf("failed to parse message body: %w", err)
	}

	if event.NotificationMessage == "" || len(event.Channels) == 0 {
		log.Println("Message has no message body or channels. Skipping.")
		return nil
	}

	var hasError bool

	for _, channel := range event.Channels {
		switch channel.Type {
		case "EMAIL":
			// Call the new, simpler SMTP function
			if err := sendEmailViaSMTP(ctx, channel.Contact, "Transaction Notification", event.NotificationMessage); err != nil {
				log.Printf("Failed to send EMAIL to %s: %v\n", channel.Contact, err)
				hasError = true
			}
		case "WHATSAPP":
			if err := sendSmsViaMeta(ctx, channel.Contact, event.NotificationMessage); err != nil {
				log.Printf("Failed to send WHATSAPP to %s: %v\n", channel.Contact, err)
				hasError = true
			}
		default:
			log.Printf("Warning: Unknown notification type '%s'\n", channel.Type)
		}
	}

	if hasError {
		return fmt.Errorf("failed to send one or more notifications")
	}

	log.Println("Successfully processed message and sent all notifications.")
	return nil
}

// --- Notification Sending Functions ---

// sendEmailViaSMTP uses Go's built-in SMTP client.
func sendEmailViaSMTP(ctx context.Context, toEmail, subject, body string) error {
	if smtpHost == "" || smtpPassword == "" {
		return fmt.Errorf("SMTP client is not initialized")
	}

	// 1. Set up authentication information.
	// We use PlainAuth, which is what most SMTP servers (including ACS) expect.
	auth := smtp.PlainAuth("", smtpUsername, smtpPassword, smtpHost)

	// 2. Format the email message (MIME format)
	// Note: We need to manually add headers.
	msg := []byte(
		"To: " + toEmail + "\r\n" +
			"From: " + smtpSender + "\r\n" +
			"Subject: " + subject + "\r\n" +
			"Content-Type: text/plain; charset=UTF-8\r\n" +
			"\r\n" + // Empty line separates headers from body
			body + "\r\n")

	// 3. Send the email
	// We combine the host and port for the address.
	addr := fmt.Sprintf("%s:%s", smtpHost, smtpPort)

	err := smtp.SendMail(addr, auth, smtpSender, []string{toEmail}, msg)
	if err != nil {
		return fmt.Errorf("SMTP SendMail failed: %w", err)
	}

	log.Printf("Successfully sent EMAIL to %s\n", toEmail)
	return nil
}

// sendSmsViaMeta... (This function remains unchanged)
func sendSmsViaMeta(ctx context.Context, toPhone, body string) error {
	if metaApiToken == "" || metaApiUrl == "" {
		return fmt.Errorf("Meta API is not configured")
	}

	// NOTE: This assumes a pre-approved template named 'transaction_update'
	payload := fmt.Sprintf(`{
        "messaging_product": "whatsapp",
        "to": "%s",
        "type": "template",
        "template": {
            "name": "transaction_update", 
            "language": { "code": "en_US" },
            "components": [
                {
                    "type": "body",
                    "parameters": [
                        { "type": "text", "text": "%s" }
                    ]
                }
            ]
        }
    }`, toPhone, body)

	req, err := http.NewRequestWithContext(ctx, "POST", metaApiUrl, bytes.NewBufferString(payload))
	if err != nil {
		return fmt.Errorf("failed to create Meta request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+metaApiToken)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send Meta request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 300 {
		var errResp bytes.Buffer
		errResp.ReadFrom(resp.Body)
		return fmt.Errorf("Meta API returned error status %d: %s", resp.StatusCode, errResp.String())
	}

	log.Printf("Successfully sent WHATSAPP to %s\n", toPhone)
	return nil
}

