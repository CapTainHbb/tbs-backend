package usecase

import (
	"golang.org/x/crypto/bcrypt"
)

// HashPassword returns the bcrypt hash of the plain password.
// cost is the work factor. bcrypt.DefaultCost is 10; 12 is common in production.
func HashPassword(password string, ) (string, error) {
	cost := bcrypt.DefaultCost
	hash, err := bcrypt.GenerateFromPassword([]byte(password), cost)
	if err != nil {
		return "", err
	}
	return string(hash), nil // store this string in DB (it includes the salt & cost)
}

// CheckPassword checks password against a bcrypt hash.
func CheckPassword(hash, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
}