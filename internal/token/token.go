package token

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
)

func Generate(length int) (string, error) {
	key := make([]byte, length)
	if _, err := rand.Read(key); err != nil {
		return "", fmt.Errorf("failed to generate key: %w", err)
	}

	return base64.URLEncoding.EncodeToString(key), nil
}
