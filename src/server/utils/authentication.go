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

func GenerateToken(tokenType string, userId string, expireMinutes int) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"token_type": tokenType,
		"iss":        "www.wilin.info",
		"sub":        userId,
		"iat":        jwt.NewNumericDate(time.Now()),
		"exp":        jwt.NewNumericDate(time.Now().Add(time.Minute * time.Duration(expireMinutes))),
	})

	secretKey := os.Getenv("SECRET_KEY")
	if secretKey == "" {
		return "", ErrEmptySecret
	}

	return token.SignedString([]byte(secretKey))
}

type TokenStruct struct {
	TokenType string
	Iss       string
	Sub       string
	Iat       float64
	Exp       float64
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

	tokenStruct.TokenType, ok = claims["token_type"].(string)
	if !ok {
		return tokenStruct, ErrInvalidToken
	}

	tokenStruct.Iss, ok = claims["iss"].(string)
	if !ok {
		return tokenStruct, ErrInvalidToken
	}

	tokenStruct.Sub, ok = claims["sub"].(string)
	if !ok {
		return tokenStruct, ErrInvalidToken
	}

	tokenStruct.Iat, ok = claims["iat"].(float64)
	if !ok {
		return tokenStruct, ErrInvalidToken
	}

	tokenStruct.Exp, ok = claims["exp"].(float64)
	if !ok {
		return tokenStruct, ErrInvalidToken
	}

	return tokenStruct, nil
}

func IsExpired(exp float64) bool {
	return time.Now().After(time.Unix(int64(exp), 0))
}
