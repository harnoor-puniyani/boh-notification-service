# Notification Service

Lightweight Go service that listens to Azure Service Bus queues and dispatches notifications (SMS and email) based on incoming messages.

## Overview

- Connects to Azure Service Bus using a connection string.
- Subscribes to two queues (or configurable queue names): one for SMS and one for Email.
- On message receipt, parses JSON payload and triggers the corresponding provider integration (SMS gateway or email provider).

## Features

- Queue-driven processing for reliable, decoupled notifications.
- Pluggable SMS / Email provider integration points.
- Basic logging and retry-friendly behavior.

## Prerequisites

- Go 1.24+ installed
- Azure Service Bus namespace with queues created
- Credentials / API keys for chosen SMS and email providers

## Configuration (Environment Variables)

- SERVICE_BUS_CONNECTION_STRING - Azure Service Bus connection string
- QUEUE_SMS_NAME - queue name for SMS messages (default: `sms`) 
- QUEUE_EMAIL_NAME - queue name for Email messages (default: `email`)
- SMS_PROVIDER_API_KEY - API key for SMS gateway (if used)
- EMAIL_PROVIDER_API_KEY - API key for email provider (if used)
- LOG_LEVEL - debug|info|warn|error (optional)
- DOTENV_FILE - optional .env file path for local development

## Message Formats

SMS queue message (JSON):
{
  "to": "+11234567890",
  "body": "Your verification code is 123456"
}

Email queue message (JSON):
{
  "to": "alice@example.com",
  "subject": "Welcome",
  "body": "Hello Alice, welcome to our service.`"
}

Messages must be valid JSON. The service will validate required fields and log/reject invalid messages.

## How it works

1. Service connects to Azure Service Bus and receives messages from configured queues.
2. For each message: parse JSON -> validate required fields -> call provider integration.
3. On success the message is completed (removed from queue). On transient failure the message is abandoned so Service Bus can retry or dead-letter per queue policy.

## Running locally

1. Copy `.env.example` to `.env` and set required variables (or set env vars directly).
2. Install dependencies and build:

   go mod tidy
   go build ./...

3. Run:

   go run ./cmd/notification-service

(Or run the built binary.)

## Deployment

- Build the binary and deploy to your target environment.
- Ensure environment variables are provided securely (Key Vault / environment configs).

## Troubleshooting

- Check logs for JSON parse errors or provider errors.
- Verify Service Bus connection string and queue names.
- Ensure provider API keys are valid.

## Extensions

- Add dead-letter handling for poison messages.
- Implement batching, rate-limiting, and metrics export (Prometheus).

