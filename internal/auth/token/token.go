package token

import (
	"errors"
	"fmt"
	"net/http"
	"strings"
	"tech-challenge-order/internal/config"

	jwt "github.com/dgrijalva/jwt-go"
)

func ValidateToken(r *http.Request) error {
	tokenString := getToken(r)

	token, err := jwt.Parse(tokenString, returnSecretKey)
	if err != nil {
		return err
	}

	if _, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return nil
	}

	return errors.New("invalid token")
}

func getToken(r *http.Request) string {
	token := r.Header.Get("Authorization")

	if len(strings.Split(token, " ")) == 2 {
		return strings.Split(token, " ")[1]
	}

	return ""
}

func returnSecretKey(token *jwt.Token) (interface{}, error) {
	if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
		return nil, fmt.Errorf("unexpected signature method %v", token.Header["alg"])
	}

	return []byte(config.Get().Token.Key), nil
}

func ExtractCustomerId(request *http.Request) (string, error) {
	token := request.Header.Get("Authorization")

	tokenS, err := jwt.Parse(token, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("invalid jwt")
		}

		return []byte(config.Get().Token.Key), nil
	})
	if err != nil {
		return "", err
	}

	return tokenS.Claims.(jwt.MapClaims)["userId"].(string), nil
}
