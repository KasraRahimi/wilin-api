package utils

import (
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"net/mail"
	"os"
	"time"
)

const BYCRYPT_COST = 13

var (
	ErrEmptySecret  = errors.New("empty secret")
	ErrInvalidToken = errors.New("invalid token")
)

func IsValidEmail(email string) bool {
	_, err := mail.ParseAddress(email)
	return err == nil
}

func GeneratePasswordHash(password string) (string, error) {
	hashBytes, err := bcrypt.GenerateFromPassword([]byte(password), BYCRYPT_COST)
	if err != nil {
		return "", fmt.Errorf("GeneratePasswordHash, generating password hash: %w", err)
	}
	return string(hashBytes), nil
}

func IsPasswordAndHashSame(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func GenerateToken(userId string, timeToExpireMinutes int) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":  userId,
		"ttl": time.Now().Add(time.Minute * time.Duration(timeToExpireMinutes)).Unix(),
	})

	secretKey := os.Getenv("SECRET_KEY")
	if secretKey == "" {
		return "", ErrEmptySecret
	}

	return token.SignedString([]byte(secretKey))
}

type TokenStruct struct {
	Id  string
	Ttl time.Time
}

func ParseToken(tokenString string) (TokenStruct, error) {
	var tokenStruct TokenStruct
	secretKey := os.Getenv("SECRET_KEY")
	if secretKey == "" {
		return tokenStruct, ErrEmptySecret
	}
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("SECRET_KEY")), nil
	})
	if err != nil {
		return tokenStruct, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return tokenStruct, ErrInvalidToken
	}

	tokenStruct.Id, ok = claims["id"].(string)
	if !ok {
		return tokenStruct, ErrInvalidToken
	}

	ttlFloat, ok := claims["ttl"].(float64) // this conversion is required due to reasons
	if !ok {
		return tokenStruct, ErrInvalidToken
	}

	tokenStruct.Ttl = time.Unix(int64(ttlFloat), 0)

	return tokenStruct, nil
}

func IsTokenTTLExpired(ttl time.Time) bool {
	return time.Now().After(ttl)
}
