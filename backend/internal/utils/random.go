package utils

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"math/big"
)

func GenerateNumericCode(length int) (string, error) {
	if length < 1 {
		return "", fmt.Errorf("invalid code length")
	}

	result := make([]byte, length)
	for i := 0; i < length; i++ {
		n, err := rand.Int(rand.Reader, big.NewInt(10))
		if err != nil {
			return "", err
		}
		result[i] = byte('0' + n.Int64())
	}

	return string(result), nil
}

func GenerateSecureToken(length int) (string, error) {
	if length < 1 {
		return "", fmt.Errorf("invalid token length")
	}

	buffer := make([]byte, length)
	if _, err := rand.Read(buffer); err != nil {
		return "", err
	}

	return base64.RawURLEncoding.EncodeToString(buffer), nil
}
