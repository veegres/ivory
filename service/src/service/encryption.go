package service

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
)

type EncryptionService struct{}

func NewEncryptionService() *EncryptionService {
	return &EncryptionService{}
}

// Encrypt method is to encrypt or hide any classified text
func (e *EncryptionService) Encrypt(text string, secret [16]byte) (string, error) {
	block, err := aes.NewCipher(secret[:])
	if err != nil {
		return "", err
	}
	plainText := []byte(text)
	cipherText := make([]byte, aes.BlockSize+len(plainText))
	iv := cipherText[:aes.BlockSize]
	stream := cipher.NewCFBEncrypter(block, iv)
	stream.XORKeyStream(cipherText[aes.BlockSize:], plainText)
	return e.encode(cipherText), nil
}

// Decrypt method is to extract back the encrypted text
func (e *EncryptionService) Decrypt(text string, secret [16]byte) (string, error) {
	block, err := aes.NewCipher(secret[:])
	if err != nil {
		return "", err
	}
	cipherText := e.decode(text)
	iv := cipherText[:aes.BlockSize]
	cipherText = cipherText[aes.BlockSize:]
	stream := cipher.NewCFBDecrypter(block, iv)
	stream.XORKeyStream(cipherText, cipherText)
	return string(cipherText), nil
}

func (e *EncryptionService) encode(b []byte) string {
	return base64.StdEncoding.EncodeToString(b)
}

func (e *EncryptionService) decode(s string) []byte {
	data, _ := base64.StdEncoding.DecodeString(s)
	return data
}
