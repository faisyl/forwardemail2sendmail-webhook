# ForwardEmail Webhook - YunoHost

Receives emails from ForwardEmail.net via webhooks and delivers them to the local Postfix mail server.

## What This Does

This application acts as a bridge between ForwardEmail.net and your YunoHost server's Postfix installation. When ForwardEmail receives an email for your domain, it sends the email as a JSON webhook to this service, which reconstructs the email and delivers it to Postfix for normal processing.

## Features

- **Webhook endpoint** - Receives ForwardEmail.net webhook payloads
- **Full email support** - Handles plain text, HTML, and file attachments
- **MIME compliant** - Properly constructs RFC 5322 compliant emails
- **Postfix integration** - Seamlessly delivers to local mail server
- **Optional authentication** - Webhook secret for security
- **Comprehensive logging** - Full audit trail of all webhook requests

## Use Cases

Perfect for:

- Using ForwardEmail.net as your email provider while hosting on YunoHost
- Receiving emails without exposing your server's SMTP port
- Centralizing email reception through ForwardEmail's infrastructure
- Adding spam filtering and email forwarding features from ForwardEmail

## Technical Details

- **Language**: Go 1.21+
- **Email Format**: RFC 5322 with MIME multipart support
- **Delivery Method**: Sendmail command-line interface
- **Process Management**: Systemd
- **Reverse Proxy**: Nginx with SSL

## Configuration

After installation, configure ForwardEmail.net to send webhooks to:

```
https://yourdomain.com/your-path/webhook/email
```

If you set a webhook secret during installation, add this header to ForwardEmail's webhook configuration:

```
X-Webhook-Secret: your-secret-here
```

## How It Works

1. Email arrives at ForwardEmail.net
2. ForwardEmail sends JSON payload to your webhook endpoint
3. Service parses JSON and reconstructs email with all headers, body, and attachments
4. Email is piped to Postfix via sendmail
5. Postfix delivers email normally to recipient mailboxes

## Getting Started

See the [README.md](file:///home/faisal/Extra/forwardemail-webhook/README.md) for detailed installation instructions, testing procedures, and troubleshooting tips.
