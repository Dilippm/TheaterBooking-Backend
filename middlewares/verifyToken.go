package middlewares

import (
	"net/http"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/dilippm92/bookingapplication/config"

	"github.com/gin-gonic/gin"
)
var jwtSecretKey = []byte(config.JWT_SECRET_KEY) 
func JwtTokenVerify()gin.HandlerFunc  {
	return func (c *gin.Context)  {
		authHeader:= c.GetHeader("Authorization")
		if authHeader == "" {
            c.JSON(http.StatusUnauthorized, gin.H{"error": "Missing Authorization Header"})
            c.Abort()
            return
        }
		// split the token string to get token part
		parts:= strings.Split(authHeader,"Bearer ")
		if len(parts) !=2{
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid Authorization Header Format"})
            c.Abort()
            return
		}
		tokenString := parts[1]
		   // Parse and validate the token
		   token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
            // Check if the token method is HMAC and return the secret key
            if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
                return nil, http.ErrNotSupported
            }
            return jwtSecretKey, nil
        })
        if err != nil || !token.Valid {
            c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid Token"})
            c.Abort()
            return
        }

        // Extract claims from the token
        claims, ok := token.Claims.(jwt.MapClaims)
        if !ok {
            c.JSON(http.StatusUnauthorized, gin.H{"errors": "Invalid Token Claims"})
            c.Abort()
            return
        }

        // Extract user ID from the claims and set it in the context
        userId, ok := claims["sub"].(string)
        if !ok {
            c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid Token Claims"})
            c.Abort()
            return
        }

        // Set the user ID in the context
        c.Set("userId", userId)

        // Call the next handler
        c.Next()
	}
}