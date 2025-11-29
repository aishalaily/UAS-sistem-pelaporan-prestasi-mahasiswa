package helper

import (
	"time"
	"os"

	"github.com/golang-jwt/jwt/v5"
)

// GenerateJWT membuat token JWT untuk user
func GenerateJWT(userID string, role string) (string, error) {
	secret := os.Getenv("JWT_SECRET")
	claims := jwt.MapClaims{
		"user_id": userID,
		"role":    role,
		"exp":     time.Now().Add(time.Hour * 24).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}

// ParseJWT memvalidasi token dan mengembalikan Claims
func ParseJWT(tokenString string) (*jwt.Token, jwt.MapClaims, error) {
	secret := os.Getenv("JWT_SECRET")

	claims := jwt.MapClaims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})

	return token, claims, err
}
