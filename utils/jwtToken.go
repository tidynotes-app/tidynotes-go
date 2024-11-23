package utils

import (
	"os"

	"github.com/golang-jwt/jwt"
)

func GenerateToken(UserID uint, UserName string) (string, error) {
	// Generating the JSON WEB TOKEN (JWT) with the "userid" and "username".
	// As I am using the free tier provided by Render,
	// I have removed the expiration time and refresh token concept.
	// This way, I no longer have to worry about refreshing the access token.
	// Even though I have removed these concepts, this application still checks for
	// invalid or tampered JWT tokens.
	jwtkey := os.Getenv("JWTKEY")
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userid":   UserID,
		"username": UserName,
	})
	// returns the Signed Token string.
	return token.SignedString([]byte(jwtkey))
}
