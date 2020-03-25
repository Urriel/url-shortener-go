package utils

import "crypto/rand"

// GenerateToken generate a random alphanumerical token
func GenerateToken() (string, error) {
	const letters = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
	buf := make([]byte, 8)

	_, err := rand.Read(buf)
	if err != nil {
		return "", err
	}

	// Translate random byte into letter
	for i, b := range buf {
		buf[i] = letters[b%byte(len(letters))]
	}

	return string(buf), nil
}
