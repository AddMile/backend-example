package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
)

var (
	ErrEmptyInput           = errors.New("inputue to encrypt or decrypt cannot be empty")
	ErrEmptyDecryptedOutput = errors.New("empty decrypted inputue")
	ErrCipherTextTooShort   = errors.New("cipher text too short")
)

// Encrypt encrypts the given inputue using AES encryption with the provided AES key.
// It returns the encrypted inputue as a hexadecimal string, or an error if encryption fails.
func Encrypt(input, aesKey string) (string, error) {
	if input == "" {
		return "", ErrEmptyInput
	}

	block, err := aes.NewCipher([]byte(aesKey))
	if err != nil {
		return "", err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}

	ciphertext := gcm.Seal(nonce, nonce, []byte(input), nil)

	return hex.EncodeToString(ciphertext), nil
}

// Decrypt decrypts the given encrypted inputue using AES decryption with the provided AES key.
// It returns the decrypted plain text inputue, or an error if decryption fails.
func Decrypt(input, aesKey string) (string, error) {
	if input == "" {
		return "", ErrEmptyInput
	}

	decodedCipherText, err := hex.DecodeString(input)
	if err != nil {
		return "", fmt.Errorf("failed to decode hex: %w", err)
	}

	block, err := aes.NewCipher([]byte(aesKey))
	if err != nil {
		return "", fmt.Errorf("failed to create cipher: %w", err)
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", fmt.Errorf("failed to create GCM: %w", err)
	}

	nonceSize := gcm.NonceSize()
	if len(decodedCipherText) < nonceSize {
		return "", ErrCipherTextTooShort
	}

	nonce, ciphertext := decodedCipherText[:nonceSize], decodedCipherText[nonceSize:]
	decryptedData, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return "", fmt.Errorf("failed to decrypt: %w", err)
	}

	decrypted := string(decryptedData)
	if decrypted == "" {
		return "", ErrEmptyDecryptedOutput
	}

	return decrypted, nil
}

// ComputeHMAC generates an HMAC-SHA256 hash of the given value using the provided secret key.
// It returns the hexadecimal-encoded string representation of the computed HMAC.
func ComputeHMAC(secret, value string) string {
	h := hmac.New(sha256.New, []byte(secret))

	h.Write([]byte(value))

	return hex.EncodeToString(h.Sum(nil))
}
