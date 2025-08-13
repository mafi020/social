package utils

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"time"
)

// GenerateOpaqueToken creates a URL-safe random string (e.g., for refresh tokens).
func GenerateOpaqueToken(nBytes int) (string, error) {
	b := make([]byte, nBytes)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	// URL-safe, no padding
	return base64.RawURLEncoding.EncodeToString(b), nil
}

// HashToken returns a hex-encoded SHA-256 hash for safe storage.
func HashToken(token string) string {
	sum := sha256.Sum256([]byte(token))
	return hex.EncodeToString(sum[:])
}

// Helper for typical refresh expiries
func RefreshExpiry(days int) time.Time {
	return time.Now().Add(time.Duration(days) * 24 * time.Hour)
}
