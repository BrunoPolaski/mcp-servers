package aes

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"io"
	"os"

	"github.com/BrunoPolaski/go-rest-err/rest_err"
)

func EncryptSecret(s string) (string, *rest_err.RestErr) {
	key := os.Getenv("AES_KEY")

	block, err := aes.NewCipher([]byte(key))
	if err != nil {
		return "", rest_err.NewInternalServerError(err.Error())
	}

	nonce := make([]byte, 12)
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", rest_err.NewInternalServerError(err.Error())
	}

	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return "", rest_err.NewInternalServerError(err.Error())
	}

	ciphertext := aesGCM.Seal(nonce, nonce, []byte(s), nil)

	s = base64.URLEncoding.EncodeToString(ciphertext)

	return s, nil
}

func DecryptSecret(secret string) (string, *rest_err.RestErr) {
	key := os.Getenv("AES_KEY")

	data, err := base64.URLEncoding.DecodeString(secret)
	if err != nil {
		return "", rest_err.NewInternalServerError(err.Error())
	}

	block, err := aes.NewCipher([]byte(key))
	if err != nil {
		return "", rest_err.NewInternalServerError(err.Error())
	}

	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return "", rest_err.NewInternalServerError(err.Error())
	}

	nonceSize := aesGCM.NonceSize()
	if len(data) < nonceSize {
		return "", rest_err.NewInternalServerError("ciphertext too short")
	}

	ciphertext := data[nonceSize:]
	nonce := data[:nonceSize]

	plaintext, err := aesGCM.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return "", rest_err.NewInternalServerError(err.Error())
	}

	return string(plaintext), nil
}
