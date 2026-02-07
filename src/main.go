package main

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/textproto"
	"os"
	"os/exec"
	"strings"
	"time"
)

// WebhookPayload represents the incoming email from ForwardEmail (mailparser output)
type WebhookPayload struct {
	Date        string            `json:"date"`
	Subject     string            `json:"subject"`
	From        AddressGroup      `json:"from"`
	To          AddressGroup      `json:"to"`
	Recipients  []string          `json:"recipients"`
	Text        string            `json:"text"`
	HTML        string            `json:"html"`
	Headers     interface{}       `json:"headers"`
	Attachments []EmailAttachment `json:"attachments"`
}

type AddressGroup struct {
	Value []AddressEntry `json:"value"`
	Text  string         `json:"text"`
}

type AddressEntry struct {
	Address string `json:"address"`
	Name    string `json:"name"`
}

// EmailAttachment represents an email attachment
type EmailAttachment struct {
	Filename    string            `json:"filename"`
	ContentType string            `json:"contentType"`
	Content     AttachmentContent `json:"content"`
}

type AttachmentContent struct {
	Type string `json:"type"` // Should be "Buffer"
	Data []int  `json:"data"` // Array of bytes as integers
}

func main() {
	// Get configuration from environment variables
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	domain := os.Getenv("DOMAIN")
	pathURL := os.Getenv("PATH_URL")
	webhookKey := os.Getenv("WEBHOOK_KEY")
	sendmailPath := os.Getenv("SENDMAIL_PATH")
	if sendmailPath == "" {
		sendmailPath = "/usr/sbin/sendmail"
	}

	// Log startup information
	log.Printf("Starting ForwardEmail Webhook Handler on port %s", port)
	log.Printf("Domain: %s, Path: %s", domain, pathURL)
	log.Printf("Sendmail path: %s", sendmailPath)
	if webhookKey != "" {
		log.Printf("Webhook key authentication enabled")
	} else {
		log.Printf("Webhook key authentication disabled (optional)")
	}

	// Normalize pathURL: ensure it starts with / and has no trailing slash
	if pathURL == "" || pathURL == "/" {
		pathURL = ""
	} else {
		if pathURL[0] != '/' {
			pathURL = "/" + pathURL
		}
		pathURL = strings.TrimSuffix(pathURL, "/")
	}

	// Set up routes with path prefix support
	http.HandleFunc(pathURL+"/", handleHome)
	http.HandleFunc(pathURL+"/health", handleHealth)
	http.HandleFunc(pathURL+"/webhook/email", makeWebhookHandler(webhookKey, sendmailPath))

	// Create server with timeouts
	server := &http.Server{
		Addr:           ":" + port,
		Handler:        nil,
		ReadTimeout:    30 * time.Second,
		WriteTimeout:   30 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	// Start server
	log.Printf("Server listening on :%s", port)
	if err := server.ListenAndServe(); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}

// makeWebhookHandler creates the webhook handler with configuration
func makeWebhookHandler(webhookKey, sendmailPath string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Received %s request at %s from %s", r.Method, r.URL.Path, r.RemoteAddr)

		// Only accept POST requests
		if r.Method != http.MethodPost {
			log.Printf("Method not allowed: %s", r.Method)
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		// Read request body first (we need it for signature verification)
		body, err := io.ReadAll(r.Body)
		if err != nil {
			log.Printf("Error reading request body: %v", err)
			http.Error(w, "Error reading request body", http.StatusInternalServerError)
			return
		}

		// Log request headers (useful for debugging signature/content-type)
		for name, values := range r.Header {
			for _, value := range values {
				log.Printf("Header: %s: %s", name, value)
			}
		}

		// Verify webhook signature if key is configured
		if webhookKey != "" {
			providedSignature := r.Header.Get("X-Webhook-Signature")
			if providedSignature == "" {
				log.Printf("Webhook authentication failed: missing signature header")
				http.Error(w, "Unauthorized: missing signature", http.StatusUnauthorized)
				return
			}

			// Compute HMAC signature of the request body
			expectedSignatureBytes := computeHMAC(body, webhookKey)

			// Compare signatures using constant-time comparison
			if !verifySignature(providedSignature, expectedSignatureBytes) {
				log.Printf("Webhook authentication failed: invalid signature")
				http.Error(w, "Unauthorized: invalid signature", http.StatusUnauthorized)
				return
			}
		}

		// Parse JSON payload (body already read above for signature verification)

		var payload WebhookPayload
		if err := json.Unmarshal(body, &payload); err != nil {
			log.Printf("Error parsing JSON payload: %v", err)
			log.Printf("Raw body was: %s", string(body))
			http.Error(w, "Invalid JSON payload", http.StatusBadRequest)
			return
		}

		// Validate required fields
		fromAddress := ""
		if len(payload.From.Value) > 0 {
			fromAddress = payload.From.Value[0].Address
		}

		toAddress := ""
		if len(payload.Recipients) > 0 {
			toAddress = payload.Recipients[0]
		} else if len(payload.To.Value) > 0 {
			toAddress = payload.To.Value[0].Address
		}

		if fromAddress == "" || toAddress == "" {
			log.Printf("Missing required fields: from=%s, recipients=%v",
				fromAddress, payload.Recipients)
			http.Error(w, "Missing required fields", http.StatusBadRequest)
			return
		}

		log.Printf("Processing email from %s to %s, subject: %s",
			fromAddress, toAddress, payload.Subject)

		// Create a buffer to construct the email in RFC822 format
		var emailBuffer bytes.Buffer

		// Write headers
		fmt.Fprintf(&emailBuffer, "From: %s\r\n", payload.From.Text)
		fmt.Fprintf(&emailBuffer, "To: %s\r\n", toAddress)
		fmt.Fprintf(&emailBuffer, "Subject: %s\r\n", payload.Subject)
		fmt.Fprintf(&emailBuffer, "Date: %s\r\n", payload.Date)
		fmt.Fprintf(&emailBuffer, "X-Forwarded-By: ForwardEmail Webhook\r\n")

		// Determine MIME structure
		hasHTML := payload.HTML != ""
		hasText := payload.Text != ""
		hasAttachments := len(payload.Attachments) > 0

		if !hasAttachments && !hasHTML {
			// Simple plain text email
			fmt.Fprintf(&emailBuffer, "Content-Type: text/plain; charset=utf-8\r\n")
			fmt.Fprintf(&emailBuffer, "\r\n")
			fmt.Fprintf(&emailBuffer, "%s\r\n", payload.Text)
		} else {
			// Multipart email
			boundary := generateBoundary()

			if hasAttachments {
				fmt.Fprintf(&emailBuffer, "MIME-Version: 1.0\r\n")
				fmt.Fprintf(&emailBuffer, "Content-Type: multipart/mixed; boundary=\"%s\"\r\n", boundary)
				fmt.Fprintf(&emailBuffer, "\r\n")

				// Write body part
				if hasHTML && hasText {
					// Nested multipart/alternative for text and HTML
					altBoundary := generateBoundary()
					fmt.Fprintf(&emailBuffer, "--%s\r\n", boundary)
					fmt.Fprintf(&emailBuffer, "Content-Type: multipart/alternative; boundary=\"%s\"\r\n", altBoundary)
					fmt.Fprintf(&emailBuffer, "\r\n")

					writeTextPart(&emailBuffer, altBoundary, payload.Text)
					writeHTMLPart(&emailBuffer, altBoundary, payload.HTML)

					fmt.Fprintf(&emailBuffer, "--%s--\r\n", altBoundary)
				} else if hasText {
					fmt.Fprintf(&emailBuffer, "--%s\r\n", boundary)
					writeTextPart(&emailBuffer, "", payload.Text)
				} else if hasHTML {
					fmt.Fprintf(&emailBuffer, "--%s\r\n", boundary)
					writeHTMLPart(&emailBuffer, "", payload.HTML)
				}

				// Write attachments
				for _, att := range payload.Attachments {
					if err := writeAttachment(&emailBuffer, boundary, att); err != nil {
						log.Printf("Warning: failed to write attachment %s: %v", att.Filename, err)
					}
				}

				fmt.Fprintf(&emailBuffer, "--%s--\r\n", boundary)
			} else {
				// multipart/alternative for text and HTML only
				fmt.Fprintf(&emailBuffer, "MIME-Version: 1.0\r\n")
				fmt.Fprintf(&emailBuffer, "Content-Type: multipart/alternative; boundary=\"%s\"\r\n", boundary)
				fmt.Fprintf(&emailBuffer, "\r\n")

				if hasText {
					writeTextPart(&emailBuffer, boundary, payload.Text)
				}
				if hasHTML {
					writeHTMLPart(&emailBuffer, boundary, payload.HTML)
				}

				fmt.Fprintf(&emailBuffer, "--%s--\r\n", boundary)
			}
		}

		// Pipe the email to sendmail
		cmd := exec.Command(sendmailPath, "-t", "-f", fromAddress)
		cmd.Stdin = &emailBuffer

		var stderr bytes.Buffer
		cmd.Stderr = &stderr

		if err := cmd.Run(); err != nil {
			log.Printf("Error sending email: %v, stderr: %s", err, stderr.String())
			http.Error(w, "Error processing email", http.StatusInternalServerError)
			return
		}

		log.Printf("Email successfully delivered to Postfix")
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, `{"status":"success","message":"Email delivered"}`)
	}
}

// writeTextPart writes a plain text MIME part
func writeTextPart(w io.Writer, boundary, text string) {
	if boundary != "" {
		fmt.Fprintf(w, "--%s\r\n", boundary)
	}
	fmt.Fprintf(w, "Content-Type: text/plain; charset=utf-8\r\n")
	fmt.Fprintf(w, "Content-Transfer-Encoding: 8bit\r\n")
	fmt.Fprintf(w, "\r\n")
	fmt.Fprintf(w, "%s\r\n", text)
}

// writeHTMLPart writes an HTML MIME part
func writeHTMLPart(w io.Writer, boundary, html string) {
	if boundary != "" {
		fmt.Fprintf(w, "--%s\r\n", boundary)
	}
	fmt.Fprintf(w, "Content-Type: text/html; charset=utf-8\r\n")
	fmt.Fprintf(w, "Content-Transfer-Encoding: 8bit\r\n")
	fmt.Fprintf(w, "\r\n")
	fmt.Fprintf(w, "%s\r\n", html)
}

// writeAttachment writes an attachment MIME part
func writeAttachment(w io.Writer, boundary string, att EmailAttachment) error {
	// Extract content bytes from the integer array
	content := make([]byte, len(att.Content.Data))
	for i, v := range att.Content.Data {
		content[i] = byte(v)
	}

	fmt.Fprintf(w, "--%s\r\n", boundary)

	// Create MIME headers for attachment
	mimeHeader := make(textproto.MIMEHeader)
	mimeHeader.Set("Content-Type", att.ContentType)
	mimeHeader.Set("Content-Transfer-Encoding", "base64")
	mimeHeader.Set("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", att.Filename))

	// Write headers
	for key, values := range mimeHeader {
		for _, value := range values {
			fmt.Fprintf(w, "%s: %s\r\n", key, value)
		}
	}
	fmt.Fprintf(w, "\r\n")

	// Write base64 encoded content (re-encode for proper line wrapping)
	encoder := base64.NewEncoder(base64.StdEncoding, &lineWrapper{w: w, lineLength: 76})
	encoder.Write(content)
	encoder.Close()
	fmt.Fprintf(w, "\r\n")

	return nil
}

// lineWrapper wraps base64 output to 76 characters per line
type lineWrapper struct {
	w           io.Writer
	lineLength  int
	currentLine int
}

func (lw *lineWrapper) Write(p []byte) (n int, err error) {
	for i, b := range p {
		if lw.currentLine >= lw.lineLength {
			if _, err := lw.w.Write([]byte("\r\n")); err != nil {
				return i, err
			}
			lw.currentLine = 0
		}
		if _, err := lw.w.Write([]byte{b}); err != nil {
			return i, err
		}
		lw.currentLine++
	}
	return len(p), nil
}

// generateBoundary creates a MIME boundary string
func generateBoundary() string {
	return fmt.Sprintf("----=_Part_%d_%d", time.Now().Unix(), time.Now().Nanosecond())
}

// handleHome serves the home page
func handleHome(w http.ResponseWriter, r *http.Request) {
	domain := os.Getenv("DOMAIN")
	pathURL := os.Getenv("PATH_URL")

	html := fmt.Sprintf(`<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>ForwardEmail Webhook - YunoHost</title>
    <style>
        body {
            font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, Oxygen, Ubuntu, Cantarell, sans-serif;
            max-width: 800px;
            margin: 50px auto;
            padding: 20px;
            background: linear-gradient(135deg, #667eea 0%%, #764ba2 100%%);
            color: #333;
        }
        .container {
            background: white;
            border-radius: 10px;
            padding: 40px;
            box-shadow: 0 10px 40px rgba(0,0,0,0.1);
        }
        h1 {
            color: #667eea;
            margin-top: 0;
        }
        .info {
            background: #f5f7fa;
            padding: 20px;
            border-radius: 5px;
            margin: 20px 0;
        }
        .info-item {
            margin: 10px 0;
        }
        .label {
            font-weight: bold;
            color: #667eea;
        }
        code {
            background: #f5f7fa;
            padding: 2px 6px;
            border-radius: 3px;
            font-family: 'Courier New', monospace;
        }
        .endpoint {
            background: #e8f5e9;
            padding: 15px;
            border-radius: 5px;
            margin: 15px 0;
            border-left: 4px solid #4caf50;
        }
    </style>
</head>
<body>
    <div class="container">
        <h1>�� ForwardEmail Webhook Handler/h1>
        <p>This service receives emails from ForwardEmail.net and delivers them to the local Postfix MTA.</p>
        
        <div class="info">
            <div class="info-item">
                <span class="label">Domain:</span> %s
            </div>
            <div class="info-item">
                <span class="label">Path:</span> %s
            </div>
            <div class="info-item">
                <span class="label">Server Time:</span> %s
            </div>
        </div>
        
        <h2>Webhook Endpoint</h2>
        <div class="endpoint">
            <strong>POST</strong> <code>https://%s%s/webhook/email</code>
        </div>
        <p>Configure this URL in your ForwardEmail.net settings to receive incoming emails.</p>
        
        <h2>Available Endpoints</h2>
        <ul>
            <li><code>/health</code> - Health check endpoint</li>
            <li><code>/webhook/email</code> - Email webhook receiver (POST only)</li>
        </ul>
        
        <h2>Status</h2>
        <p>✅ Service is running and ready to receive webhooks</p>
    </div>
</body>
</html>`, domain, pathURL, time.Now().Format(time.RFC1123), domain, pathURL)

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	fmt.Fprint(w, html)
}

// handleHealth provides a health check endpoint
func handleHealth(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(w, `{"status":"healthy","timestamp":"%s","service":"forwardemail-webhook"}`, time.Now().Format(time.RFC3339))
}
