package auth

import (
	"errors"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type TokenType string

const (
	AccessJWTIssuer  TokenType = "chirpy-access"
	RefreshJWTIssuer TokenType = "chirpy-refresh"
)

var (
	AccessTokenDuration  time.Duration = time.Hour
	RefreshTokenDuration time.Duration = time.Hour * 24 * 60
)

func GenerateToken(issuer TokenType, secret string, expiry time.Duration, id int) (string, error) {
	signingKey := []byte(secret)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		Issuer:    string(issuer),
		IssuedAt:  jwt.NewNumericDate(time.Now().UTC()),
		ExpiresAt: jwt.NewNumericDate(time.Now().UTC().Add(expiry)),
		Subject:   strconv.Itoa(id),
	})

	return token.SignedString(signingKey)
}

func GeneratePasswordHash(password string) ([]byte, error) {
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return []byte{}, err
	}
	return passwordHash, nil
}

func AuthorizePasswordHash(storedPassword, requestPassword string) error {
	err := bcrypt.CompareHashAndPassword([]byte(storedPassword), []byte(requestPassword))
	if err != nil {
		return err
	}
	return nil
}

func GetBearerToken(headers http.Header) (string, error) {
	authHeader := headers.Get("Authorization")
	if authHeader == "" {
		return "", errors.New("no authorization included")
	}
	splitAuth := strings.Split(authHeader, " ")
	if len(splitAuth) < 2 || splitAuth[0] != "Bearer" {
		return "", errors.New("malformed authorization header")
	}
	return splitAuth[1], nil
}

func GetApiKeyToken(headers http.Header) (string, error) {
	authHeader := headers.Get("Authorization")
	if authHeader == "" {
		return "", errors.New("no authorization included")
	}
	splitAuth := strings.Split(authHeader, " ")
	if len(splitAuth) < 2 || splitAuth[0] != "ApiKey" {
		return "", errors.New("malformed authorization header")
	}
	return splitAuth[1], nil
}

func ValidateJWT(tokenString, tokenSecret string, JWTIssuer TokenType) (string, error) {
	claimsStruct := jwt.RegisteredClaims{}
	token, err := jwt.ParseWithClaims(
		tokenString,
		&claimsStruct,
		func(token *jwt.Token) (interface{}, error) { return []byte(tokenSecret), nil })
	if err != nil {
		return "", err
	}

	userIDString, err := token.Claims.GetSubject()
	if err != nil {
		return "", err
	}

	issuer, err := token.Claims.GetIssuer()
	if err != nil {
		return "", err
	}
	if issuer != string(JWTIssuer) {
		return "", errors.New("invalid issuer")
	}

	return userIDString, nil
}

func ValidateRefreshJWT(tokenString, tokenSecret string) (string, error) {
	claimsStruct := jwt.RegisteredClaims{}
	token, err := jwt.ParseWithClaims(
		tokenString,
		&claimsStruct,
		func(token *jwt.Token) (interface{}, error) { return []byte(tokenSecret), nil })
	if err != nil {
		return "", err
	}

	userIDString, err := token.Claims.GetSubject()
	if err != nil {
		return "", err
	}

	issuer, err := token.Claims.GetIssuer()
	if err != nil {
		return "", err
	}
	if issuer != string(RefreshJWTIssuer) {
		return "", errors.New("invalid issuer")
	}

	return userIDString, nil
}

func ValidateUserAccess(r *http.Request, secret string) (string, error) {
	tokenString, err := GetBearerToken(r.Header)
	if err != nil {
		return "", err
	}

	// requires an access token
	subject, err := ValidateJWT(tokenString, secret, AccessJWTIssuer)
	if err != nil {
		return "", err
	}
	return subject, nil
}
