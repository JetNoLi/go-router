package utils

import (
	"fmt"
	"net/http"
)

// More Generic: Update to Take in a value
func GenerateAuthCookie(userId int) (*http.Cookie, error) {
	userIdStr := fmt.Sprintf("%d", userId)
	token, err := GenerateJwt(userIdStr)

	if err != nil {
		return &http.Cookie{}, err
	}

	return &http.Cookie{
		Name:     "Authorization",
		Value:    token,
		HttpOnly: true,
		Secure:   true,
		MaxAge:   60 * 60 * 24 * 7, // 1 Week
		Path:     "/",
	}, nil
}
