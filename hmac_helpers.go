package main

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
)

// computeHMAC generates an HMAC SHA-256 signature for the given data
func computeHMAC(data []byte, secret string) []byte {
	h := hmac.New(sha256.New, []byte(secret))
	h.Write(data)
	return h.Sum(nil)
}

// verifySignature performs constant-time comparison of two signatures.
// providedHex is the hex-encoded signature from the header.
// expectedBytes is the raw byte slice of the expected HMAC.
func verifySignature(providedHex string, expectedBytes []byte) bool {
	providedBytes, err := hex.DecodeString(providedHex)
	if err != nil {
		return false
	}

	return hmac.Equal(providedBytes, expectedBytes)
}
