package security

import (
	"crypto/sha256"
	"encoding/hex"
)

// HashPassword returns password hash
func HashPassword(password string) string {
	hash := sha256.Sum256([]byte(password))
	return hex.EncodeToString(hash[:])
}

// CheckPasswordHash checks password hash and request password hash
func CheckPasswordHash(password, hash string) bool {
	return HashPassword(password) == hash
}
