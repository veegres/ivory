package encryption

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"io"
)

type Service struct{}

func NewService() *Service {
	return &Service{}
}

// Encrypt method encrypts text using AES-GCM (authenticated encryption)
// The IV (nonce) is randomly generated and prepended to the ciphertext
func (e *Service) Encrypt(text string, secret [16]byte) (string, error) {
	if text == "" {
		return "", errors.New("cannot encrypt empty string")
	}
	if secret == [16]byte{} {
		return "", errors.New("secret is empty")
	}

	block, err := aes.NewCipher(secret[:])
	if err != nil {
		return "", err
	}

	// Create GCM mode cipher
	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	// Generate a random nonce (IV)
	// GCM standard nonce size is 12 bytes
	nonce := make([]byte, aesGCM.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}

	plainText := []byte(text)

	// Encrypt and authenticate
	// GCM appends authentication tag automatically
	cipherText := aesGCM.Seal(nonce, nonce, plainText, nil)

	return e.encode(cipherText), nil
}

// Decrypt method decrypts text using AES-GCM
// Verifies authentication tag and returns error if tampered
func (e *Service) Decrypt(text string, secret [16]byte) (string, error) {
	if text == "" {
		return "", errors.New("cannot decrypt empty string")
	}
	if secret == [16]byte{} {
		return "", errors.New("secret is empty")
	}

	block, err := aes.NewCipher(secret[:])
	if err != nil {
		return "", err
	}

	// Create GCM mode cipher
	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	cipherText, err := e.decode(text)
	if err != nil {
		return "", err
	}

	// Check minimum length (nonce + at least some data)
	nonceSize := aesGCM.NonceSize()
	if len(cipherText) < nonceSize {
		return "", errors.New("ciphertext too short")
	}

	// Extract nonce and ciphertext
	nonce, cipherText := cipherText[:nonceSize], cipherText[nonceSize:]

	// Decrypt and verify authentication tag
	plainText, err := aesGCM.Open(nil, nonce, cipherText, nil)
	if err != nil {
		return "", errors.New("decryption failed: authentication check failed or invalid ciphertext")
	}

	return string(plainText), nil
}

func (e *Service) encode(b []byte) string {
	return base64.StdEncoding.EncodeToString(b)
}

func (e *Service) decode(s string) ([]byte, error) {
	data, err := base64.StdEncoding.DecodeString(s)
	if err != nil {
		return nil, errors.New("invalid base64 encoding")
	}
	return data, nil
}
