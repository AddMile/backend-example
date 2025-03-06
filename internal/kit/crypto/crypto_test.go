package crypto_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	cryptokit "github.com/AddMile/backend/internal/kit/crypto"
)

const AESKey = "a1b2c3d4e5f6a7b8c9d0e1f2a3b4c5d6"

func TestEncryptDecrypt(t *testing.T) {
	input := "john"

	encrypted, err := cryptokit.Encrypt(input, AESKey)
	assert.NoError(t, err)

	decrypted, err := cryptokit.Decrypt(encrypted, AESKey)
	assert.NoError(t, err)

	assert.Equal(t, input, decrypted)
}

func TestComputeHMAC(t *testing.T) {
	secret := "secret"
	email := "test@gmail.com"

	expected := "f816c5ca6d01f2949af51978cfd26d74c69436de4d41dd864a24ab58a3962027"
	actual := cryptokit.ComputeHMAC(secret, email)
	assert.Equal(t, expected, actual)
}
