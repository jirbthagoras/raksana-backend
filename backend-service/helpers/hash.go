package helpers

import (
	"golang.org/x/crypto/bcrypt"
)

// HashPassword hashes a plain password with bcrypt
func HashPassword(password string) (string, error) {
	// bcrypt.DefaultCost = 10, you can increase for more security (but slower)
	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashed), nil
}

// CheckPassword compares a plain password with a hashed one
func CheckPassword(password, hashed string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashed), []byte(password))
	return err == nil
}
