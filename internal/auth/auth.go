/* Package auth
* Package auth is for handling authenctication
* */
package auth

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"net/http"
	"strings"

	"github.com/alexedwards/argon2id"
)

// ErrorNoAuthHeaderIncluded -
var ErrorNoAuthHeaderIncluded = errors.New("no auth header included in the request")

// HashPassword -
func HashPassword(password string) (string, error) {
	hash, err := argon2id.CreateHash(password, argon2id.DefaultParams)
	if err != nil {
		return "", err
	}

	return hash, nil
}

// CheckPasswordHash -
func CheckPasswordHash(password, hash string) (bool, error) {
	match, err := argon2id.ComparePasswordAndHash(password, hash)
	if err != nil {
		return false, err
	}

	return match, nil
}

// GetBearerToken -
func GetBearerToken(headers http.Header) (string, error) {
	authHeader := headers.Get("Authorization")
	if authHeader == "" {
		return "", ErrorNoAuthHeaderIncluded
	}

	authSlice := strings.Split(authHeader, " ")
	if authSlice[0] != "Bearer" || len(authSlice) != 2 {
		return "", errors.New("malformed auth header")
	}
	return authSlice[1], nil
}

// MakeRefreshToken -
func MakeRefreshToken() string {
	key := make([]byte, 32)
	rand.Read(key)
	return hex.EncodeToString(key)
}

// GetAPIKey -
func GetAPIKey(headers http.Header) (string, error) {
	authHeader := headers.Get("Authorization")
	if authHeader == "" {
		return "", ErrorNoAuthHeaderIncluded
	}

	authSlice := strings.Split(authHeader, " ")
	if authSlice[0] != "ApiKey" || len(authSlice) != 2 {
		return "", errors.New("malformed auth header")
	}
	return authSlice[1], nil
}
