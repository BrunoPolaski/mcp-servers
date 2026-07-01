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
		return "", rest_err.NewInternalServerError("%s", err.Error()).WithCause(err)
	}

	nonce := make([]byte, 12)
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", rest_err.NewInternalServerError("%s", err.Error()).WithCause(err)
	}

	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return "", rest_err.NewInternalServerError("%s", err.Error()).WithCause(err)
	}

	ciphertext := aesGCM.Seal(nonce, nonce, []byte(s), nil)

	s = base64.URLEncoding.EncodeToString(ciphertext)

	return s, nil
}

func DecryptSecret(secret string) (string, *rest_err.RestErr) {
	key := os.Getenv("AES_KEY")

	data, err := base64.URLEncoding.DecodeString(secret)
	if err != nil {
		return "", rest_err.NewInternalServerError("%s", err.Error()).WithCause(err)
	}

	block, err := aes.NewCipher([]byte(key))
	if err != nil {
		return "", rest_err.NewInternalServerError("%s", err.Error()).WithCause(err)
	}

	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return "", rest_err.NewInternalServerError("%s", err.Error()).WithCause(err)
	}

	nonceSize := aesGCM.NonceSize()
	if len(data) < nonceSize {
		return "", rest_err.NewInternalServerError("ciphertext too short")
	}

	ciphertext := data[nonceSize:]
	nonce := data[:nonceSize]

	plaintext, err := aesGCM.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return "", rest_err.NewInternalServerError("%s", err.Error()).WithCause(err)
	}

	return string(plaintext), nil
}
