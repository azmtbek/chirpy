/* Package auth
* Package auth is for handling authenctication
* */
package auth

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/alexedwards/argon2id"
)

func HashPassword(password string) (string, error) {
	hash, err := argon2id.CreateHash(password, argon2id.DefaultParams)
	if err != nil {
		return "", err
	}

	return hash, nil
}

func CheckPasswordHash(password, hash string) (bool, error) {
	match, err := argon2id.ComparePasswordAndHash(password, hash)
	if err != nil {
		return false, err
	}

	return match, nil
}

func GetBearerToken(headers http.Header) (string, error) {
	authHeader := headers.Get("Authorization")

	if authHeader == "" {
		return "", fmt.Errorf("no authorization header")
	}

	authSlice := strings.Split(authHeader, " ")
	if authSlice[0] != "Bearer" || len(authSlice) != 2 {
		return "", fmt.Errorf("invalid auth header")
	}
	return authSlice[1], nil
}
