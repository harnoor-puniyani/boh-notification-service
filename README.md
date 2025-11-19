# Notification Service

Lightweight Go service that listens to Azure Service Bus queues and dispatches notifications (SMS and email) based on incoming messages.

## Overview

- Connects to Azure Service Bus using a connection string.
- Subscribes to two configurable queues (one for SMS, one for Email) and processes messages reliably.
- Parses message payloads and forwards notifications to provider integrations (SMS gateway, SMTP/email provider).
- Designed for simple integration and easy local development.

## Features

- Queue-driven processing for decoupled notifications.
- Pluggable provider integrations (SMTP, SMS APIs).
- Basic logging and retry-friendly behavior; messages are completed only on success.

## Prerequisites

- Go 1.24+ installed
- Azure Service Bus namespace and required queues created
- Credentials / API keys for chosen SMS and email providers (or SMTP server credentials)

## Configuration (Environment Variables)

You can use a `.env` file for local development or provide environment variables directly.

Required:
- SERVICE_BUS_CONNECTION_STRING - Azure Service Bus connection string
- QUEUE_SMS_NAME - queue name for SMS messages (default: `sms`)
- QUEUE_EMAIL_NAME - queue name for Email messages (default: `email`)

Optional / Provider keys:
- SMS_PROVIDER_API_KEY - API key for SMS gateway (if used)
- SMS_PROVIDER_URL - HTTP endpoint for SMS API
- SMTP_HOST - SMTP server host
- SMTP_PORT - SMTP server port
- SMTP_USERNAME - SMTP username
- SMTP_PASSWORD - SMTP password
- SMTP_SENDER - Sender email address used in From header

Logging and runtime:
- LOG_LEVEL - debug|info|warn|error (optional)
- DOTENV_FILE - optional .env file path for local development

Example `.env`:

```
SERVICE_BUS_CONNECTION_STRING=Endpoint=sb://...;SharedAccessKeyName=...;SharedAccessKey=...
QUEUE_SMS_NAME=sms
QUEUE_EMAIL_NAME=email
SMS_PROVIDER_API_KEY=your_sms_api_key
SMS_PROVIDER_URL=https://api.sms-provider.example/send
SMTP_HOST=smtp.example.com
SMTP_PORT=587
SMTP_USERNAME=user@example.com
SMTP_PASSWORD=secret
SMTP_SENDER=no-reply@example.com
LOG_LEVEL=debug
```

## Message Formats

Messages must be valid JSON. The service expects a simple payload describing the notification and the target channels. Example envelope:

SMS queue (single channel or multiple):

```json
{
  "notificationId": "uuid-or-id",
  "notificationMessage": "Your verification code is 123456",
  "channels": [
    { "type": "sms", "contact": "+11234567890" }
  ]
}
```

Email queue message:

```json
{
  "notificationId": "uuid-or-id",
  "notificationMessage": "<p>Your statement is ready</p>",
  "channels": [
    { "type": "email", "contact": "alice@example.com", "subject": "Monthly Statement" }
  ]
}
```

Fields:
- notificationId (optional): id for tracing
- notificationMessage: string (plain text or HTML for email)
- channels: array of channel objects with at minimum `type` and `contact`. Type values used in this service: `sms`, `email`, `whatsapp` (if implemented).
- email channels may include `subject`.

The service will validate required fields and log & abandon invalid messages so they can be retried or dead-lettered per queue policy.

## How it works

1. Service connects to Service Bus and receives messages from configured queues.
2. For each message: parse JSON -> validate required fields -> call provider integration.
3. On success the message is completed (removed from queue). On transient failure the message is abandoned so Service Bus can retry or dead-letter per queue settings.

## Running locally

1. Copy `.env.example` to `.env` and set required variables (or set env vars directly).
2. Fetch dependencies and build:

   go mod tidy
   go build ./...

3. Run:

   DOTENV_FILE=.env go run ./cmd/notification-service

Or build and run the binary:

   go build -o bin/notification ./cmd/notification-service
   DOTENV_FILE=.env ./bin/notification

## Docker

You can containerize the service. Example Dockerfile (simplified):

```
FROM golang:1.24-alpine AS builder
WORKDIR /src
COPY . .
RUN go build -o /app/bin/notification ./cmd/notification-service

FROM alpine:3.18
COPY --from=builder /app/bin/notification /notification
ENV LOG_LEVEL=info
CMD ["/notification"]
```

Run with env vars supplied (example):

docker build -t boh-notification .
docker run -e SERVICE_BUS_CONNECTION_STRING="..." -e QUEUE_SMS_NAME=sms -e QUEUE_EMAIL_NAME=email boh-notification

## Provider integration

- SMTP: uses net/smtp for simple SMTP sending. Confirm TLS/STARTTLS requirements for your provider. Some providers require explicit TLS or OAuth flows.
- SMS: use HTTP API integration. Implement the adapter in `notifier/` to match your provider's API.

## Error handling and retries

- The service should only complete a Service Bus message when the provider call succeeds.
- On provider errors, return a non-success so Service Bus can retry/abandon according to configured policies.
- Consider dead-lettering messages after repeated failures and implementing a monitoring / alerting mechanism.

## Testing

- Use a local .env to point to a test Service Bus namespace.
- Send test messages via Azure Portal or using `az` CLI.
- Tail logs to observe processing and provider calls.

## Troubleshooting

- "SMTP not configured": ensure SMTP env vars are present.
- Authentication errors from SMS/email providers: verify API keys and provider requirements (TLS, IP allowlists, credentials).
- Service Bus connectivity errors: confirm connection string and network access to the namespace.

## Extensions and improvements

- Add structured logging and metrics (Prometheus) for observability.
- Support batching and rate-limiting to avoid provider throttling.
- Implement dead-letter processing and a retry/alert dashboard.

If you want, I can add example provider adapters (e.g., Twilio, SendGrid) or create a docker-compose example for local testing. Feel free to ask for those next steps.

