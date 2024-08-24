package helpers

import(
	"golang.org/x/crypto/bcrypt"
)
func HashPassword(password string) (string, error) {
	// bcrypt generates a hashed version of the password with a cost of 14
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}