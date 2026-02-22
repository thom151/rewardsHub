package auth

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/alexedwards/argon2id"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

func HashPassword(password string) (string, error) {
	hashed, err := argon2id.CreateHash(password, argon2id.DefaultParams)
	if err != nil {
		return "", err
	}

	return hashed, nil
}

func CheckPasswordHash(password, hash string) (bool, error) {
	match, err := argon2id.ComparePasswordAndHash(password, hash)
	if err != nil {
		return false, err
	}

	return match, nil
}

func MakeJWT(userID uuid.UUID, tokenSecret string, expiresIn time.Duration) (string, error) {

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		Issuer:    "leadway",
		IssuedAt:  jwt.NewNumericDate(time.Now().UTC()),
		ExpiresAt: jwt.NewNumericDate(time.Now().UTC().Add(expiresIn)),
		Subject:   userID.String(),
	})

	strToken, err := token.SignedString([]byte(tokenSecret))
	if err != nil {
		return "", err
	}

	return strToken, nil
}

func ValidateJWT(tokenString, tokenSecret string) (uuid.UUID, error) {
	claimsStruct := jwt.RegisteredClaims{}

	token, err := jwt.ParseWithClaims(
		tokenString,
		&claimsStruct,
		func(token *jwt.Token) (interface{}, error) { return []byte(tokenSecret), nil },
	)

	if err != nil {
		return uuid.Nil, err
	}

	userIDString, err := token.Claims.GetSubject()
	if err != nil {
		return uuid.Nil, err
	}

	issuer, err := token.Claims.GetIssuer()
	if err != nil {
		return uuid.Nil, err
	}
	if issuer != "leadway" {
		return uuid.Nil, errors.New("invalid issuer")
	}

	id, err := uuid.Parse(userIDString)
	if err != nil {
		return uuid.Nil, fmt.Errorf("invalid user ID: %w", err)
	}
	return id, nil
}

func GetBearerToken(headers http.Header, cookies []*http.Cookie) (string, error) {

	for _, cookie := range cookies {
		log.Printf("cookie name: %v\n", cookie.Name)
		if cookie.Name == "acc_token" {
			return cookie.Value, nil
		}
	}

	authHeader := headers.Get("Authorization")
	if authHeader == "" {
		return "", fmt.Errorf("no auth header")
	}
	splitAuth := strings.Split(authHeader, " ")
	if len(splitAuth) == 2 || splitAuth[0] == "Bearer" {
		return splitAuth[1], nil
	}

	return "", fmt.Errorf("malformed bearer")
}

func MakeRefreshToken() string {
	token := make([]byte, 32)
	rand.Read(token)
	return hex.EncodeToString(token)
}

// RUNES FOR AUTH IDENTITY PROVIDER
type AuthProvider string

const (
	ProviderEmail  AuthProvider = "email"
	ProviderGoogle AuthProvider = "google"
)
