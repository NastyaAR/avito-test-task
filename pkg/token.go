package pkg

import (
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

func GenerateJWTToken(userId uuid.UUID, role string) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["userID"] = userId
	claims["role"] = role
	tokenString, err := token.SignedString([]byte("some key"))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
