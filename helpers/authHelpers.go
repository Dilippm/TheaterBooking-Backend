package helpers

import (
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/dilippm92/bookingapplication/config"
	"golang.org/x/crypto/bcrypt"
)
var jwtSecretKey =[]byte(config.JWT_SECRET_KEY)
// hash input password
func HashPassword(password string) (string, error) {
	// bcrypt generates a hashed version of the password with a cost of 14
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}
// compare passwords
func ComparePasswords(hashedPassword, password string) error {

	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}
// generate jwt token
func GenerateJWTToken(id string, email string) (string, error) {
	// Define the expiration time of the token (e.g., 24 hours)
	expirationTime := time.Now().Add(24 * time.Hour)
	// Create a new token object
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub":   id,
		"email": email,                 // Subject (user ID)
		"exp":   expirationTime.Unix(), // Expiration time
	})

	// Sign the token with the secret key
	tokenString, err := token.SignedString(jwtSecretKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}