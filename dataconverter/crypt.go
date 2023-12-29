package dataconverter

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"fmt"
	"io"
)

func encrypt(plainData []byte, key []byte) ([]byte, error) {
	c, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(c)
	if err != nil {
		return nil, err
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, err
	}

	return gcm.Seal(nonce, nonce, plainData, nil), nil
}

func decrypt(encryptedData []byte, key []byte) ([]byte, error) {
	c, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(c)
	if err != nil {
		return nil, err
	}

	nonceSize := gcm.NonceSize()
	if len(encryptedData) < nonceSize {
		return nil, fmt.Errorf("ciphertext too short: %v", encryptedData)
	}

	nonce, encryptedData := encryptedData[:nonceSize], encryptedData[nonceSize:]
	return gcm.Open(nil, nonce, encryptedData, nil)
}

func Decrypt(encryptedData []byte, key []byte) ([]byte, error) {
	c, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	fmt.Println("Decypt: c:", c)

	gcm, err := cipher.NewGCM(c)
	if err != nil {
		return nil, err
	}
	fmt.Println("Decypt: gcm:", gcm)

	nonceSize := gcm.NonceSize()
	if len(encryptedData) < nonceSize {
		return nil, fmt.Errorf("ciphertext too short: %v", encryptedData)
	}
	fmt.Println("Decypt: nonceSize:", nonceSize)

	fmt.Println("Decypt: encryptedData[:nonceSize]:", encryptedData[:nonceSize])
	fmt.Println("Decypt: encryptedData[nonceSize:]:", encryptedData[nonceSize:])

	nonce, encryptedData := encryptedData[:nonceSize], encryptedData[nonceSize:]

	fmt.Println("Decypt: nonce:", nonce)
	fmt.Println("Decypt: encryptedData:", encryptedData)
	return gcm.Open(nil, nonce, encryptedData, nil)
}
