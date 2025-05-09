package security

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestHashPassword(t *testing.T) {
	t.Run("hash is consistent", func(t *testing.T) {
		password := "mypassword"
		hash1 := HashPassword(password)
		hash2 := HashPassword(password)
		assert.Equal(t, hash1, hash2, "Hashes for the same password should be identical")
	})

	t.Run("different passwords produce different hashes", func(t *testing.T) {
		hash1 := HashPassword("password1")
		hash2 := HashPassword("password2")
		assert.NotEqual(t, hash1, hash2, "Different passwords should produce different hashes")
	})

	t.Run("hash is non-empty", func(t *testing.T) {
		hash := HashPassword("test")
		assert.NotEmpty(t, hash, "Hash should not be empty")
	})
}

func TestCheckPasswordHash(t *testing.T) {
	t.Run("valid password", func(t *testing.T) {
		password := "mypassword"
		hash := HashPassword(password)
		assert.True(t, CheckPasswordHash(password, hash), "Valid password should match hash")
	})

	t.Run("invalid password", func(t *testing.T) {
		password := "mypassword"
		hash := HashPassword(password)
		assert.False(t, CheckPasswordHash("wrongpassword", hash), "Invalid password should not match hash")
	})

	t.Run("empty password", func(t *testing.T) {
		hash := HashPassword("")
		assert.True(t, CheckPasswordHash("", hash), "Empty password should match its hash")
	})
}
