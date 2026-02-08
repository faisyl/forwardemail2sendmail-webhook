# forwardingress for YunoHost

A professional and slick YunoHost application that acts as a bridge, relaying webhooks seamlessly to your email inbox.

## Overview

**forwardingress** provides a universal webhook endpoint for services like ForwardEmail.net. It receives email data as a JSON payload, reconstructs it into a proper RFC 5322 compliant email (with full MIME support for HTML and attachments), and relays it via your choice of backend: local Postfix (sendmail) or remote SMTP.

## Features

- ✅ **Slick Landing Page** with real-time status and configuration details
- ✅ **Multiple Backends**: Support for local `sendmail` or remote `SMTP`
- ✅ **Full MIME support**: Handles plain text, HTML, and complex attachments
- ✅ **HMAC Security**: Signature verification for secure webhook processing
- ✅ **YunoHost Integrated**: Standard v2 packaging, systemd service, and portal icon
- ✅ **Self-Contained**: Embedded assets for a zero-dependency frontend

## How It Works

1. **Web Service triggers Webhook** → Data arrives at your `forwardingress` endpoint.
2. **Branding & Verification** → The service verifies the signature and parses the payload.
3. **Email Construction** → `forwardingress` builds a standards-compliant email.
4. **Seamless Relay** → The email is relayed to your inbox via the chosen backend.

## Installation on YunoHost

### From Command Line

```bash
yunohost app install https://github.com/faisyl/forwardemail2sendmail-webhook
```

### Configuration

During installation/upgrade, you can configure:
- **Webhook Key**: For HMAC signature verification (supports ForwardEmail and others).
- **Delivery Backend**: Choose between `sendmail` (default) or `SMTP`.
- **SMTP Settings**: Host, Port, Credentials, and TLS SkipVerify options.

## Development

### Structure
Organized according to the `my_webapp_ynh` standard:
- `sources/`: Go application source code
- `scripts/`: YunoHost packaging scripts
- `logos/`: Standard application icons
- `doc/`: Documentation and internal assets

### Building
```bash
cd sources
go build -o ../bin/web2mail .
```

## License
MIT License - see LICENSE file
