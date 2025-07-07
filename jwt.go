package genh

import (
	"strings"
	"sync"

	"github.com/golang-jwt/jwt/v5"
)

var (
	globalSecret []byte
	once         sync.Once
)

type UserClaims struct {
	UserID   int    `json:"id"`
	Username string `json:"username"`
	jwt.RegisteredClaims
}

func GenTokenWithHS256(claims *UserClaims, secret []byte) (string, error) {
	once.Do(func() {
		globalSecret = secret
	})
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(globalSecret)
}

func ParseTokenWithHS256(tokenStr string) *UserClaims {
	if globalSecret == nil {
		return nil
	}
	claims := &UserClaims{}
	token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
		return globalSecret, nil
	})
	if err != nil || !token.Valid {
		return nil
	}
	return claims
}

func ParseContextWithHS256(h interface{ Header(string) string }) *UserClaims {
	tokenStr := h.Header("Authorization")
	tokenStr = strings.TrimPrefix(tokenStr, "Bearer ")
	if globalSecret == nil {
		return nil
	}
	claims := &UserClaims{}
	token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
		return globalSecret, nil
	})
	if err != nil || !token.Valid {
		return nil
	}
	return claims
}
