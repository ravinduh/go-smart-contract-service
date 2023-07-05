package service

import (
	"encoding/base64"
	"github.com/go-kit/kit/log"
)

type EncryptionHelper struct {
	logger log.Logger
}

func (eh EncryptionHelper) EncryptBase64(requestId, plaintext string) (string, error) {
	// Convert the plaintext to a byte array
	plaintextBytes := []byte(plaintext)

	// Encrypt the plaintext using Base64 encoding
	ciphertext := base64.StdEncoding.EncodeToString(plaintextBytes)

	return ciphertext, nil
}

func (eh EncryptionHelper) DecryptBase64(requestId, ciphertext string) (string, error) {
	// Decode the ciphertext from Base64 encoding
	ciphertextBytes, err := base64.StdEncoding.DecodeString(ciphertext)
	if err != nil {
		return "", err
	}

	// Convert the ciphertext byte array to a string
	plaintext := string(ciphertextBytes)

	return plaintext, nil
}

func NewEncryptionHelper(logger log.Logger) EncryptionHelper {
	return EncryptionHelper{
		logger: logger,
	}
}
