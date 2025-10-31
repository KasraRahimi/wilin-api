package services

import (
	"errors"
	"fmt"
	"log"
	"net/mail"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

const BYCRYPT_COST = 13

var SIGNING_ALG = jwt.SigningMethodHS256

var FAKE_PASSWORD string

var (
	ErrEmptySecret  = errors.New("empty secret")
	ErrInvalidToken = errors.New("invalid token")
)

type MyJWTClaims struct {
	Type string `json:"type,omitempty"`
	jwt.RegisteredClaims
}

func init() {
	passwordHash, err := GeneratePasswordHash("contrasena")
	if err != nil {
		log.Printf("failed to generate fake password hash: %v\n", err)
		FAKE_PASSWORD = "$2a$13$GcHjrPzhZzc4ZYOmaN1Kc.ohBOOT4Z9DtESP3aASZd94xl85NCjee"
		return
	}

	log.Println("generated fake password hash")
	FAKE_PASSWORD = passwordHash
}

func IsValidEmail(email string) bool {
	_, err := mail.ParseAddress(email)
	return err == nil
}

func GeneratePasswordHash(password string) (string, error) {
	hashBytes, err := bcrypt.GenerateFromPassword([]byte(password), BYCRYPT_COST)
	if err != nil {
		return "", err
	}
	return string(hashBytes), nil
}

func IsPasswordAndHashSame(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func FakeHashCompare() {
	_ = IsPasswordAndHashSame("motdepasse", FAKE_PASSWORD)
}

func GenerateToken(tokenType string, userId string, expireMinutes int) (string, error) {
	now := jwt.NewNumericDate(time.Now())
	expiry := jwt.NewNumericDate(time.Now().Add(time.Minute * time.Duration(expireMinutes)))

	claims := MyJWTClaims{
		Type: tokenType,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    "www.wilin.info",
			Subject:   userId,
			IssuedAt:  now,
			ExpiresAt: expiry,
		},
	}

	token := jwt.NewWithClaims(SIGNING_ALG, claims)

	secretKey := os.Getenv("SECRET_KEY")
	if secretKey == "" {
		return "", ErrEmptySecret
	}

	return token.SignedString([]byte(secretKey))
}

func ParseToken(tokenString string) (*MyJWTClaims, error) {
	secretKey := os.Getenv("SECRET_KEY")
	if secretKey == "" {
		return nil, ErrEmptySecret
	}

	claims := &MyJWTClaims{}

	token, err := jwt.ParseWithClaims(
		tokenString,
		claims,
		func(token *jwt.Token) (interface{}, error) {
			return []byte(os.Getenv("SECRET_KEY")), nil
		},
		jwt.WithValidMethods([]string{SIGNING_ALG.Alg()}),
	)
	if err != nil {
		return nil, fmt.Errorf("error parsing token: %v", err)
	}

	claims, ok := token.Claims.(*MyJWTClaims)
	if ok && token.Valid {
		return claims, nil
	}

	return nil, fmt.Errorf("invalid token")
}
