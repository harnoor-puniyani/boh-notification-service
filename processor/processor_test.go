package processor

import (
	"context"
	"encoding/json"
	"os"
	"testing"
)

// TestInit_AllEnvVarsSet tests Init when all environment variables are set
func TestInit_AllEnvVarsSet(t *testing.T) {
	// Set up environment variables
	os.Setenv("SMTP_HOST", "smtp.test.com")
	os.Setenv("SMTP_PORT", "587")
	os.Setenv("SMTP_USERNAME", "testuser")
	os.Setenv("SMTP_PASSWORD", "testpass")
	os.Setenv("ACS_SENDER_EMAIL", "sender@test.com")
	os.Setenv("META_API_TOKEN", "test_token")
	os.Setenv("META_API_URL", "https://api.test.com")
	
	defer func() {
		os.Unsetenv("SMTP_HOST")
		os.Unsetenv("SMTP_PORT")
		os.Unsetenv("SMTP_USERNAME")
		os.Unsetenv("SMTP_PASSWORD")
		os.Unsetenv("ACS_SENDER_EMAIL")
		os.Unsetenv("META_API_TOKEN")
		os.Unsetenv("META_API_URL")
	}()
	
	// Call Init
	Init()
	
	// Verify variables are set
	if smtpHost != "smtp.test.com" {
		t.Errorf("Expected smtpHost 'smtp.test.com', got: %s", smtpHost)
	}
	if smtpPort != "587" {
		t.Errorf("Expected smtpPort '587', got: %s", smtpPort)
	}
	if smtpUsername != "testuser" {
		t.Errorf("Expected smtpUsername 'testuser', got: %s", smtpUsername)
	}
	if smtpPassword != "testpass" {
		t.Errorf("Expected smtpPassword 'testpass', got: %s", smtpPassword)
	}
	if smtpSender != "sender@test.com" {
		t.Errorf("Expected smtpSender 'sender@test.com', got: %s", smtpSender)
	}
	if metaApiToken != "test_token" {
		t.Errorf("Expected metaApiToken 'test_token', got: %s", metaApiToken)
	}
	if metaApiUrl != "https://api.test.com" {
		t.Errorf("Expected metaApiUrl 'https://api.test.com', got: %s", metaApiUrl)
	}
}

// TestInit_MissingSMTPVars tests Init when SMTP environment variables are missing
func TestInit_MissingSMTPVars(t *testing.T) {
	// Clear SMTP environment variables
	os.Unsetenv("SMTP_HOST")
	os.Unsetenv("SMTP_PORT")
	os.Unsetenv("SMTP_USERNAME")
	os.Unsetenv("SMTP_PASSWORD")
	os.Unsetenv("ACS_SENDER_EMAIL")
	
	// Set Meta variables
	os.Setenv("META_API_TOKEN", "test_token")
	os.Setenv("META_API_URL", "https://api.test.com")
	
	defer func() {
		os.Unsetenv("META_API_TOKEN")
		os.Unsetenv("META_API_URL")
	}()
	
	// Call Init - should log warning but not fail
	Init()
	
	// Verify SMTP variables are empty
	if smtpHost != "" {
		t.Errorf("Expected empty smtpHost, got: %s", smtpHost)
	}
}

// TestInit_MissingMetaVars tests Init when Meta environment variables are missing
func TestInit_MissingMetaVars(t *testing.T) {
	// Set SMTP variables
	os.Setenv("SMTP_HOST", "smtp.test.com")
	os.Setenv("SMTP_PORT", "587")
	os.Setenv("SMTP_USERNAME", "testuser")
	os.Setenv("SMTP_PASSWORD", "testpass")
	os.Setenv("ACS_SENDER_EMAIL", "sender@test.com")
	
	// Clear Meta variables
	os.Unsetenv("META_API_TOKEN")
	os.Unsetenv("META_API_URL")
	
	defer func() {
		os.Unsetenv("SMTP_HOST")
		os.Unsetenv("SMTP_PORT")
		os.Unsetenv("SMTP_USERNAME")
		os.Unsetenv("SMTP_PASSWORD")
		os.Unsetenv("ACS_SENDER_EMAIL")
	}()
	
	// Call Init - should log warning but not fail
	Init()
	
	// Verify Meta variables are empty
	if metaApiToken != "" {
		t.Errorf("Expected empty metaApiToken, got: %s", metaApiToken)
	}
	if metaApiUrl != "" {
		t.Errorf("Expected empty metaApiUrl, got: %s", metaApiUrl)
	}
}

// TestProcessMessage_ValidJSON tests ProcessMessage with valid JSON
func TestProcessMessage_ValidJSON(t *testing.T) {
	event := NotificationEvent{
		UserID:              "user123",
		NotificationMessage: "Test message",
		Channels: []NotificationChannel{
			{Type: "UNKNOWN", Contact: "test@example.com"},
		},
	}
	
	messageBody, err := json.Marshal(event)
	if err != nil {
		t.Fatalf("Failed to marshal event: %v", err)
	}
	
	ctx := context.Background()
	err = ProcessMessage(ctx, messageBody)
	
	// Should not error on unknown channel type
	if err != nil {
		t.Errorf("Expected no error for unknown channel type, got: %v", err)
	}
}

// TestProcessMessage_InvalidJSON tests ProcessMessage with invalid JSON
func TestProcessMessage_InvalidJSON(t *testing.T) {
	invalidJSON := []byte(`{"invalid json"}`)
	
	ctx := context.Background()
	err := ProcessMessage(ctx, invalidJSON)
	
	if err == nil {
		t.Error("Expected error for invalid JSON, got nil")
	}
}

// TestProcessMessage_EmptyMessage tests ProcessMessage with empty notification message
func TestProcessMessage_EmptyMessage(t *testing.T) {
	event := NotificationEvent{
		UserID:              "user123",
		NotificationMessage: "",
		Channels: []NotificationChannel{
			{Type: "EMAIL", Contact: "test@example.com"},
		},
	}
	
	messageBody, err := json.Marshal(event)
	if err != nil {
		t.Fatalf("Failed to marshal event: %v", err)
	}
	
	ctx := context.Background()
	err = ProcessMessage(ctx, messageBody)
	
	// Should skip processing when message is empty
	if err != nil {
		t.Errorf("Expected no error for empty message, got: %v", err)
	}
}

// TestProcessMessage_NoChannels tests ProcessMessage with no channels
func TestProcessMessage_NoChannels(t *testing.T) {
	event := NotificationEvent{
		UserID:              "user123",
		NotificationMessage: "Test message",
		Channels:            []NotificationChannel{},
	}
	
	messageBody, err := json.Marshal(event)
	if err != nil {
		t.Fatalf("Failed to marshal event: %v", err)
	}
	
	ctx := context.Background()
	err = ProcessMessage(ctx, messageBody)
	
	// Should skip processing when no channels
	if err != nil {
		t.Errorf("Expected no error for empty channels, got: %v", err)
	}
}

// TestProcessMessage_EmailChannel_NotConfigured tests ProcessMessage with EMAIL channel but no SMTP config
func TestProcessMessage_EmailChannel_NotConfigured(t *testing.T) {
	// Clear SMTP configuration
	smtpHost = ""
	smtpPassword = ""
	
	event := NotificationEvent{
		UserID:              "user123",
		NotificationMessage: "Test message",
		Channels: []NotificationChannel{
			{Type: "EMAIL", Contact: "test@example.com"},
		},
	}
	
	messageBody, err := json.Marshal(event)
	if err != nil {
		t.Fatalf("Failed to marshal event: %v", err)
	}
	
	ctx := context.Background()
	err = ProcessMessage(ctx, messageBody)
	
	if err == nil {
		t.Error("Expected error for unconfigured SMTP, got nil")
	}
}

// TestProcessMessage_WhatsAppChannel_NotConfigured tests ProcessMessage with WHATSAPP channel but no Meta config
func TestProcessMessage_WhatsAppChannel_NotConfigured(t *testing.T) {
	// Clear Meta configuration
	metaApiToken = ""
	metaApiUrl = ""
	
	event := NotificationEvent{
		UserID:              "user123",
		NotificationMessage: "Test message",
		Channels: []NotificationChannel{
			{Type: "WHATSAPP", Contact: "+1234567890"},
		},
	}
	
	messageBody, err := json.Marshal(event)
	if err != nil {
		t.Fatalf("Failed to marshal event: %v", err)
	}
	
	ctx := context.Background()
	err = ProcessMessage(ctx, messageBody)
	
	if err == nil {
		t.Error("Expected error for unconfigured Meta API, got nil")
	}
}

// TestProcessMessage_MultipleChannels tests ProcessMessage with multiple channels
func TestProcessMessage_MultipleChannels(t *testing.T) {
	event := NotificationEvent{
		UserID:              "user123",
		NotificationMessage: "Test message",
		Channels: []NotificationChannel{
			{Type: "UNKNOWN", Contact: "test1@example.com"},
			{Type: "INVALID", Contact: "test2@example.com"},
		},
	}
	
	messageBody, err := json.Marshal(event)
	if err != nil {
		t.Fatalf("Failed to marshal event: %v", err)
	}
	
	ctx := context.Background()
	err = ProcessMessage(ctx, messageBody)
	
	// Should process all channels without error (unknown types are ignored)
	if err != nil {
		t.Errorf("Expected no error for unknown channel types, got: %v", err)
	}
}

// TestProcessMessage_MixedChannels tests ProcessMessage with valid and invalid channels
func TestProcessMessage_MixedChannels(t *testing.T) {
	// Clear configurations to ensure errors
	smtpHost = ""
	smtpPassword = ""
	
	event := NotificationEvent{
		UserID:              "user123",
		NotificationMessage: "Test message",
		Channels: []NotificationChannel{
			{Type: "UNKNOWN", Contact: "test@example.com"},
			{Type: "EMAIL", Contact: "test@example.com"},
		},
	}
	
	messageBody, err := json.Marshal(event)
	if err != nil {
		t.Fatalf("Failed to marshal event: %v", err)
	}
	
	ctx := context.Background()
	err = ProcessMessage(ctx, messageBody)
	
	// Should error because EMAIL is not configured
	if err == nil {
		t.Error("Expected error for unconfigured EMAIL channel, got nil")
	}
}

// TestSendEmailViaSMTP_NotInitialized tests sendEmailViaSMTP when SMTP is not initialized
func TestSendEmailViaSMTP_NotInitialized(t *testing.T) {
	// Clear SMTP configuration
	smtpHost = ""
	smtpPassword = ""
	
	ctx := context.Background()
	err := sendEmailViaSMTP(ctx, "test@example.com", "Test Subject", "Test Body")
	
	if err == nil {
		t.Error("Expected error for uninitialized SMTP, got nil")
	}
	if err.Error() != "SMTP client is not initialized" {
		t.Errorf("Expected 'SMTP client is not initialized' error, got: %v", err)
	}
}

// TestSendEmailViaSMTP_MissingHost tests sendEmailViaSMTP when host is missing
func TestSendEmailViaSMTP_MissingHost(t *testing.T) {
	smtpHost = ""
	smtpPassword = "password"
	
	ctx := context.Background()
	err := sendEmailViaSMTP(ctx, "test@example.com", "Test Subject", "Test Body")
	
	if err == nil {
		t.Error("Expected error for missing SMTP host, got nil")
	}
}

// TestSendEmailViaSMTP_MissingPassword tests sendEmailViaSMTP when password is missing
func TestSendEmailViaSMTP_MissingPassword(t *testing.T) {
	smtpHost = "smtp.test.com"
	smtpPassword = ""
	
	ctx := context.Background()
	err := sendEmailViaSMTP(ctx, "test@example.com", "Test Subject", "Test Body")
	
	if err == nil {
		t.Error("Expected error for missing SMTP password, got nil")
	}
}

// TestSendSmsViaMeta_NotConfigured tests sendSmsViaMeta when Meta API is not configured
func TestSendSmsViaMeta_NotConfigured(t *testing.T) {
	// Clear Meta configuration
	metaApiToken = ""
	metaApiUrl = ""
	
	ctx := context.Background()
	err := sendSmsViaMeta(ctx, "+1234567890", "Test message")
	
	if err == nil {
		t.Error("Expected error for unconfigured Meta API, got nil")
	}
	if err.Error() != "Meta API is not configured" {
		t.Errorf("Expected 'Meta API is not configured' error, got: %v", err)
	}
}

// TestSendSmsViaMeta_MissingToken tests sendSmsViaMeta when token is missing
func TestSendSmsViaMeta_MissingToken(t *testing.T) {
	metaApiToken = ""
	metaApiUrl = "https://api.test.com"
	
	ctx := context.Background()
	err := sendSmsViaMeta(ctx, "+1234567890", "Test message")
	
	if err == nil {
		t.Error("Expected error for missing Meta API token, got nil")
	}
}

// TestSendSmsViaMeta_MissingUrl tests sendSmsViaMeta when URL is missing
func TestSendSmsViaMeta_MissingUrl(t *testing.T) {
	metaApiToken = "test_token"
	metaApiUrl = ""
	
	ctx := context.Background()
	err := sendSmsViaMeta(ctx, "+1234567890", "Test message")
	
	if err == nil {
		t.Error("Expected error for missing Meta API URL, got nil")
	}
}

// TestNotificationChannel_Struct tests the NotificationChannel struct
func TestNotificationChannel_Struct(t *testing.T) {
	channel := NotificationChannel{
		Type:    "EMAIL",
		Contact: "test@example.com",
	}
	
	if channel.Type != "EMAIL" {
		t.Errorf("Expected Type 'EMAIL', got: %s", channel.Type)
	}
	if channel.Contact != "test@example.com" {
		t.Errorf("Expected Contact 'test@example.com', got: %s", channel.Contact)
	}
}

// TestNotificationEvent_Struct tests the NotificationEvent struct
func TestNotificationEvent_Struct(t *testing.T) {
	event := NotificationEvent{
		UserID:              "user123",
		NotificationMessage: "Test message",
		Channels: []NotificationChannel{
			{Type: "EMAIL", Contact: "test@example.com"},
		},
	}
	
	if event.UserID != "user123" {
		t.Errorf("Expected UserID 'user123', got: %s", event.UserID)
	}
	if event.NotificationMessage != "Test message" {
		t.Errorf("Expected NotificationMessage 'Test message', got: %s", event.NotificationMessage)
	}
	if len(event.Channels) != 1 {
		t.Errorf("Expected 1 channel, got: %d", len(event.Channels))
	}
}

// TestNotificationEvent_JSONMarshaling tests JSON marshaling/unmarshaling
func TestNotificationEvent_JSONMarshaling(t *testing.T) {
	original := NotificationEvent{
		UserID:              "user456",
		NotificationMessage: "Test notification",
		Channels: []NotificationChannel{
			{Type: "EMAIL", Contact: "user@example.com"},
			{Type: "WHATSAPP", Contact: "+1234567890"},
		},
	}
	
	// Marshal to JSON
	jsonData, err := json.Marshal(original)
	if err != nil {
		t.Fatalf("Failed to marshal NotificationEvent: %v", err)
	}
	
	// Unmarshal back
	var unmarshaled NotificationEvent
	err = json.Unmarshal(jsonData, &unmarshaled)
	if err != nil {
		t.Fatalf("Failed to unmarshal NotificationEvent: %v", err)
	}
	
	// Verify fields
	if unmarshaled.UserID != original.UserID {
		t.Errorf("UserID mismatch: got %s, want %s", unmarshaled.UserID, original.UserID)
	}
	if unmarshaled.NotificationMessage != original.NotificationMessage {
		t.Errorf("NotificationMessage mismatch: got %s, want %s", unmarshaled.NotificationMessage, original.NotificationMessage)
	}
	if len(unmarshaled.Channels) != len(original.Channels) {
		t.Errorf("Channels length mismatch: got %d, want %d", len(unmarshaled.Channels), len(original.Channels))
	}
	for i := range original.Channels {
		if unmarshaled.Channels[i].Type != original.Channels[i].Type {
			t.Errorf("Channel %d Type mismatch: got %s, want %s", i, unmarshaled.Channels[i].Type, original.Channels[i].Type)
		}
		if unmarshaled.Channels[i].Contact != original.Channels[i].Contact {
			t.Errorf("Channel %d Contact mismatch: got %s, want %s", i, unmarshaled.Channels[i].Contact, original.Channels[i].Contact)
		}
	}
}

// TestProcessMessage_ContextCancellation tests ProcessMessage with cancelled context
func TestProcessMessage_ContextCancellation(t *testing.T) {
	event := NotificationEvent{
		UserID:              "user123",
		NotificationMessage: "Test message",
		Channels: []NotificationChannel{
			{Type: "UNKNOWN", Contact: "test@example.com"},
		},
	}
	
	messageBody, err := json.Marshal(event)
	if err != nil {
		t.Fatalf("Failed to marshal event: %v", err)
	}
	
	// Create a cancelled context
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	
	// Should still process successfully since we don't check context in ProcessMessage
	// This documents the current behavior
	err = ProcessMessage(ctx, messageBody)
	
	if err != nil {
		t.Errorf("Got error: %v", err)
	}
}

// TestProcessMessage_EmptyBody tests ProcessMessage with empty message body
func TestProcessMessage_EmptyBody(t *testing.T) {
	ctx := context.Background()
	err := ProcessMessage(ctx, []byte{})
	
	// Should error on empty body
	if err == nil {
		t.Error("Expected error for empty message body, got nil")
	}
}

// TestProcessMessage_NilBody tests ProcessMessage with nil message body
func TestProcessMessage_NilBody(t *testing.T) {
	ctx := context.Background()
	err := ProcessMessage(ctx, nil)
	
	// Should error on nil body
	if err == nil {
		t.Error("Expected error for nil message body, got nil")
	}
}
