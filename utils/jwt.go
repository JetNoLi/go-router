package utils

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type UserJwt struct {
	UserId    string
	ExpiresAt string
}

var secret = []byte("MY_JWT_SECRET") // Replace with your desired secret

func GenerateJwt(userId string) (string, error) {

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userId": userId,
		"iss":    "notion-voice-assistant",         // Replace with App Name
		"exp":    time.Now().Add(time.Hour).Unix(), // Expiration time
		"iat":    time.Now().Unix(),                // Issued at
	})

	return token.SignedString(secret)
}

func DecodeJwt(rawToken string) (string, error) {
	token, err := jwt.Parse(rawToken, func(t *jwt.Token) (interface{}, error) {
		_, ok := t.Method.(*jwt.SigningMethodHMAC)

		if !ok {
			return "", fmt.Errorf("invalid token type")
		}

		return secret, nil
	})

	if err != nil {
		return "", err
	}

	if !token.Valid {
		return "", fmt.Errorf("invalid token")
	}

	claimMap, ok := token.Claims.(jwt.MapClaims)

	if !ok {
		return "", fmt.Errorf("cannot read jwt claims")
	}

	userId := claimMap["userId"].(string)

	return userId, nil

}
