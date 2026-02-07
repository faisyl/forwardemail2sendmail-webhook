# ForwardEmail Webhook for YunoHost

A YunoHost application that receives emails from ForwardEmail.net via webhooks and delivers them to the local Postfix MTA.

## Overview

This application provides a webhook endpoint that ForwardEmail.net can send incoming emails to. The service receives the email as a JSON payload, reconstructs it as a proper RFC 5322 compliant email (with MIME support for HTML, attachments, etc.), and pipes it to the local Postfix server using the `sendmail` command.

## Features

- ✅ **Webhook endpoint** for ForwardEmail.net integration
- ✅ **Full MIME support** - handles plain text, HTML, and attachments
- ✅ **RFC 5322 compliant** email reconstruction
- ✅ **Postfix integration** via sendmail command
- ✅ **Optional webhook authentication** with secret token
- ✅ **Comprehensive logging** for debugging and monitoring
- ✅ **YunoHost v2 packaging** with full backup/restore support
- ✅ **Multi-instance support** - run multiple instances if needed

## How It Works

1. **ForwardEmail receives email** → Incoming email arrives at ForwardEmail.net
2. **Webhook triggered** → ForwardEmail sends JSON payload to your endpoint
3. **Email reconstructed** → Service parses JSON and builds proper email format
4. **Delivered to Postfix** → Email piped to local Postfix via sendmail
5. **Normal delivery** → Postfix handles the email like any other incoming mail

## Installation on YunoHost

### From Command Line

```bash
yunohost app install /path/to/forwardemail-webhook
```

During installation, you'll be prompted for:

- **Domain**: The domain where the webhook will be accessible
- **Path**: The URL path (default: `/goapp`)
- **Admin**: The YunoHost user who will be the admin
- **Public access**: Set to "visitors" to allow ForwardEmail to access the webhook
- **Webhook secret** (optional): A secret token for authenticating webhook requests

### After Installation

Your webhook endpoint will be available at:
```
https://yourdomain.com/your-path/webhook/email
```

## ForwardEmail Configuration

1. Log in to your ForwardEmail.net account
2. Go to your domain settings
3. Configure a webhook for incoming emails
4. Set the webhook URL to: `https://yourdomain.com/your-path/webhook/email`
5. If you configured a webhook secret, add it as a header: `X-Webhook-Secret: your-secret`

## Testing

### Test with curl

```bash
# Simple test without attachment
curl -X POST https://yourdomain.com/your-path/webhook/email \
  -H "Content-Type: application/json" \
  -H "X-Webhook-Secret: your-secret" \
  -d @test_payload.json

# Test with attachment
curl -X POST https://yourdomain.com/your-path/webhook/email \
  -H "Content-Type: application/json" \
  -d @test_payload_with_attachment.json
```

### Check Postfix Queue

```bash
# View mail queue
mailq

# View mail logs
tail -f /var/log/mail.log
```

## Development

### Prerequisites

- Go 1.21 or higher
- Postfix installed and configured

### Building Locally

```bash
cd src
go build -o ../bin/forwardemail-webhook .
```

### Running Locally

```bash
PORT=8080 \
DOMAIN=localhost \
PATH_URL=/ \
WEBHOOK_SECRET=test-secret \
SENDMAIL_PATH=/usr/sbin/sendmail \
./bin/forwardemail-webhook
```

### Project Structure

```
.
├── manifest.toml              # YunoHost app manifest
├── conf/
│   ├── nginx.conf            # Nginx reverse proxy config
│   └── systemd.service       # Systemd service template
├── scripts/
│   ├── install               # Installation script
│   ├── remove                # Removal script
│   ├── backup                # Backup script
│   ├── restore               # Restore script
│   ├── upgrade               # Upgrade script
│   └── change_url            # Domain/path change script
├── src/
│   ├── main.go               # Main application
│   └── go.mod                # Go module file
├── test_payload.json         # Test webhook payload
└── test_payload_with_attachment.json  # Test with attachment
```

## Environment Variables

The application uses these environment variables (configured automatically by YunoHost):

- `PORT` - Port to listen on (managed by YunoHost)
- `DOMAIN` - Domain where app is installed
- `PATH_URL` - URL path where app is accessible
- `WEBHOOK_SECRET` - Optional secret for webhook authentication
- `SENDMAIL_PATH` - Path to sendmail binary (default: `/usr/sbin/sendmail`)

## Webhook Payload Format

ForwardEmail sends emails as JSON with this structure:

```json
{
  "date": "2026-02-06T22:45:00Z",
  "subject": "Email subject",
  "from_address": "sender@example.com",
  "from_name": "Sender Name",
  "to_address": "recipient@yourdomain.com",
  "headers": {
    "message-id": "<unique-id@example.com>",
    "reply-to": "reply@example.com"
  },
  "content": {
    "text": "Plain text body",
    "html": "<html>HTML body</html>"
  },
  "attachments": [
    {
      "filename": "file.pdf",
      "content_type": "application/pdf",
      "content": "base64-encoded-data"
    }
  ]
}
```

## Troubleshooting

### Emails not being delivered

1. **Check the application logs:**
   ```bash
   journalctl -u forwardemail-webhook -f
   ```

2. **Verify Postfix is running:**
   ```bash
   systemctl status postfix
   ```

3. **Check Postfix logs:**
   ```bash
   tail -f /var/log/mail.log
   ```

4. **Test sendmail directly:**
   ```bash
   echo "Test email" | sendmail -t recipient@example.com
   ```

### Webhook authentication failing

- Verify the `X-Webhook-Secret` header matches the configured secret
- Check application logs for authentication errors

### Permission errors

- Ensure the service user has permission to execute sendmail
- Check systemd service logs for permission denied errors

## Security

- **Webhook authentication**: Use the optional webhook secret to prevent unauthorized access
- **HTTPS only**: Always use HTTPS for webhook endpoints (handled by YunoHost)
- **Input validation**: All email fields are validated before processing
- **Sandboxing**: Service runs with systemd security restrictions

## License

MIT License - see LICENSE file

## Resources

- [ForwardEmail.net Documentation](https://forwardemail.net/)
- [YunoHost Documentation](https://yunohost.org/en/packaging_apps)
- [RFC 5322 - Internet Message Format](https://tools.ietf.org/html/rfc5322)
