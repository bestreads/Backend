package dtos

import (
	"strconv"

	"github.com/golang-jwt/jwt/v5"
)

type CustomTokenClaims struct {
	jwt.RegisteredClaims
	TokenType string `json:"token_type"`
}

func (claims *CustomTokenClaims) GetId() (uint, error) {
	subject := claims.Subject
	id, err := strconv.ParseUint(subject, 10, 0)
	if err != nil {
		return 0, err
	}

	return uint(id), nil
}
