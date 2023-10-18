package helper

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"fmt"
)

func EncryptSeedPhrase(seedPhrase, key string) ([]byte, error) {
	block, err := aes.NewCipher([]byte(key))
	if err != nil {
		return nil, err
	}

	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	nonce := make([]byte, aesgcm.NonceSize())
	if _, err := rand.Read(nonce); err != nil {
		return nil, err
	}

	// Encrypt the data
	ciphertext := aesgcm.Seal(nil, nonce, []byte(seedPhrase), nil)

	// Combine the nonce and ciphertext for storage
	encryptedData := append(nonce, ciphertext...)

	return encryptedData, nil
}

func DecryptSeedPhrase(encryptedSeedPhrase []byte, key string) ([]byte, error) {
	block, err := aes.NewCipher([]byte(key))
	if err != nil {
		return nil, err
	}

	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	nonceSize := aesgcm.NonceSize()
	if len(encryptedSeedPhrase) < nonceSize {
		return nil, fmt.Errorf("ciphertext too short")
	}

	nonce := encryptedSeedPhrase[:nonceSize]
	ciphertext := encryptedSeedPhrase[nonceSize:]

	// Decrypt the data
	plaintext, err := aesgcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return nil, err
	}

	return plaintext, nil
}
