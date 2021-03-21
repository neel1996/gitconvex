package utils

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"io"
)

type PasswordCipherInterface interface {
	EncryptPassword() string
	DecryptPassword() string
}

type PasswordCipherStruct struct {
	PlainPassword     string
	EncryptedPassword string
	KeyString         string
}

// Encrypts the password and returns the encrypted string
func (x PasswordCipherStruct) EncryptPassword() string {
	plainText := x.PlainPassword
	keyString := x.KeyString

	keyBytes := []byte(keyString + keyString)
	plainBytes := []byte(plainText)

	if keyBytes != nil {
		block, blockErr := aes.NewCipher(keyBytes)
		if blockErr != nil {
			fmt.Println(blockErr.Error())
			return ""
		}

		aesGCM, gcmErr := cipher.NewGCM(block)
		if gcmErr != nil {
			fmt.Println(gcmErr.Error())
			return ""
		}

		nonce := make([]byte, aesGCM.NonceSize())
		if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
			return ""
		}

		encryptedBytes := aesGCM.Seal(nonce, nonce, plainBytes, nil)
		return base64.StdEncoding.EncodeToString(encryptedBytes)
	}
	return ""
}

// Decrypts the AES encrypted password
func (x PasswordCipherStruct) DecryptPassword() string {
	keyString := x.KeyString
	password := x.EncryptedPassword
	encBytes, _ := base64.StdEncoding.DecodeString(password)
	keyBytes := []byte(keyString + keyString)

	if keyBytes == nil && encBytes == nil {
		return ""
	}

	block, blockErr := aes.NewCipher(keyBytes)
	if blockErr != nil {
		fmt.Println(blockErr.Error())
		return ""
	}
	aesGCM, gcmErr := cipher.NewGCM(block)
	if gcmErr != nil {
		fmt.Println(gcmErr.Error())
		return ""
	}

	nonceSize := aesGCM.NonceSize()
	nonce, cipherText := encBytes[:nonceSize], encBytes[nonceSize:]
	plainText, decryptErr := aesGCM.Open(nil, nonce, cipherText, nil)
	if decryptErr != nil {
		fmt.Println(decryptErr.Error())
		return ""
	}
	return string(plainText)
}
