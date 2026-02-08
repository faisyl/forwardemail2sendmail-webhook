# forwardingress - YunoHost

A professional and slick bridge that relays webhooks seamlessly to your email inbox.

## What This Does

**forwardingress** acts as a professional bridge between web services (like ForwardEmail.net) and your YunoHost mail infrastructure. When a service sends a JSON webhook to your endpoint, **forwardingress** reconstructs the email and relays it via your chosen backend (local Postfix or remote SMTP).

## Features

- **Universal Webhook Endpoint** - Receives JSON payloads and converts them to email.
- **Slick Management UI** - Minimalist landing page with status and diagnostic info.
- **Multi-Backend Support** - Choose between local `sendmail` or remote `SMTP`.
- **Full Email Support** - Handles plain text, HTML, and complex attachments.
- **HMAC Security** - Secure your endpoint with signature verification.

## Use Cases

- Hosting personal or service-specific email receiving without exposing SMTP ports.
- Integrating external email providers like ForwardEmail.net with your YunoHost instance.
- Creating professional automated workflows that trigger notifications into your inbox.

## How It Works

1. A web service triggers a POST request to your webhook endpoint.
2. **forwardingress** verifies the signature and parses the email data.
3. The service builds a standards-compliant email with all headers and attachments.
4. The email is relayed via your configured backend.

## Getting Started

See the [README.md](file:///home/faisal/Extra/forwardemail-webhook/README.md) for detailed configuration guide and testing procedures.
