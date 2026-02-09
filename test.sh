#!/bin/bash

# Test script for ForwardEmail webhook with HMAC signature verification
# This script tests the webhook endpoint locally

set -e

echo "=== ForwardEmail Webhook Test Script ==="
echo ""

# Check if binary exists
if [ ! -f "web2mail" ]; then
    echo "❌ Binary not found. Building..."
    go build -o web2mail .
    echo "✅ Build complete"
fi

# Find the binary
BINARY="web2mail"

echo "Using binary: $BINARY"
echo ""

# Start server in background
echo "Starting webhook server..."
PORT=8080 \
DOMAIN=localhost \
PATH_URL=/ \
WEBHOOK_KEY=test-secret \
SENDMAIL_PATH="$(pwd)/mock-sendmail.sh" \
./$BINARY &

SERVER_PID=$!
echo "Server started with PID: $SERVER_PID"

# Wait for server to start
sleep 2

# Function to compute HMAC signature
compute_signature() {
    local payload_file=$1
    local secret=$2
    echo -n "$(cat "$payload_file" | openssl dgst -sha256 -hmac "$secret" | awk '{print $2}')"
}

# Test health endpoint
echo ""
echo "=== Testing health endpoint ==="
curl -s http://localhost:8080/health | jq . || echo "Failed to parse JSON"

# Test webhook endpoint with simple payload
echo ""
echo "=== Testing webhook endpoint (simple) ==="
SIGNATURE=$(compute_signature "test_payload.json" "test-secret")
echo "Computed signature: $SIGNATURE"
curl -s -X POST http://localhost:8080/webhook/email \
  -H "Content-Type: application/json" \
  -H "X-Webhook-Signature: $SIGNATURE" \
  --data-binary @test_payload.json | jq . || echo "Failed"

# Test webhook endpoint with attachment
echo ""
echo "=== Testing webhook endpoint (with attachment) ==="
SIGNATURE=$(compute_signature "test_payload_with_attachment.json" "test-secret")
echo "Computed signature: $SIGNATURE"
curl -s -X POST http://localhost:8080/webhook/email \
  -H "Content-Type: application/json" \
  -H "X-Webhook-Signature: $SIGNATURE" \
  --data-binary @test_payload_with_attachment.json | jq . || echo "Failed"

# Test authentication failure (wrong signature)
echo ""
echo "=== Testing authentication failure (wrong signature) ==="
curl -s -X POST http://localhost:8080/webhook/email \
  -H "Content-Type: application/json" \
  -H "X-Webhook-Signature: invalid-signature-12345" \
  --data-binary @test_payload.json || echo "Expected failure"

# Test missing signature
echo ""
echo "=== Testing missing signature ==="
curl -s -X POST http://localhost:8080/webhook/email \
  -H "Content-Type: application/json" \
  -d @test_payload.json || echo "Expected failure"

# Cleanup
echo ""
echo "=== Cleaning up ==="
kill $SERVER_PID
wait $SERVER_PID 2>/dev/null || true

echo ""
echo "✅ Tests complete!"
