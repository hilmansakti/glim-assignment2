package utils

import (
	"bytes"
	"crypto/rand"
	"encoding/base64"
	"golang.org/x/crypto/scrypt"
	"strings"
)

func GenerateRandomKey(length int) []byte {
	key := make([]byte, length)
	return key
}

func HashPassword(password string) (string, error) {
	salt := GenerateRandomKey(32)
	if _, err := rand.Read(salt); err != nil {
		return "", err
	}

	hashedPassword, err := scrypt.Key([]byte(password), salt, 16384, 8, 1, 32)
	if err != nil {
		return "", err
	}

	encodedSalt := base64.StdEncoding.EncodeToString(salt)
	encodedPassword := base64.StdEncoding.EncodeToString(hashedPassword)

	return encodedSalt + "-" + encodedPassword, nil
}

func CheckPasswordHash(password, hash string) bool {
	parts := strings.Split(hash, "-")
	if len(parts) != 2 {
		return false
	}

	encodedSalt := parts[0]
	encodedHashedPassword := parts[1]

	salt, err := base64.StdEncoding.DecodeString(encodedSalt)
	if err != nil {
		return false
	}

	hashedPassword, err := base64.StdEncoding.DecodeString(encodedHashedPassword)
	if err != nil {
		return false
	}

	testHash, err := scrypt.Key([]byte(password), salt, 16384, 8, 1, 32)
	if err != nil {
		return false
	}

	return bytes.Equal(testHash, hashedPassword)
}
