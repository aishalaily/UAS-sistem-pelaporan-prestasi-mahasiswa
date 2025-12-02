package utils

import (
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"uas-go/app/model"
)

func GenerateJWT(user model.User, roles model.Role) (string, error) {
	secret := os.Getenv("JWT_SECRET")

	claims := model.JWTClaims{
		UserID:   user.ID,
		Username: user.Username,
		RoleName: roles.RoleName,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(3 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString([]byte(secret))
}

func ParseJWT(tokenStr string) (*model.JWTClaims, error) {
	secret := os.Getenv("JWT_SECRET")

	token, err := jwt.ParseWithClaims(tokenStr, &model.JWTClaims{}, func(t *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})

	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*model.JWTClaims)
	if !ok || !token.Valid {
		return nil, err
	}

	return claims, nil
}
