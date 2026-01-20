package services

import (
	"context"
	"fmt"
	"time"

	"github.com/bestreads/Backend/internal/dtos"
	"github.com/bestreads/Backend/internal/middlewares"
	"github.com/bestreads/Backend/internal/types"
	"github.com/golang-jwt/jwt/v5"
)

func GenerateToken(ctx context.Context, userId string, tokenType types.TokenType) (string, error) {
	cfg := middlewares.Config(ctx)

	// Get correct expiry duration
	var duration time.Duration
	switch tokenType {
	case types.AccessToken:
		duration = time.Duration(cfg.AccessTokenDurationMinutes) * time.Minute
	case types.RefreshToken:
		duration = time.Duration(cfg.RefreshTokenDurationDays) * 24 * time.Hour
	default:
		duration = 15 * time.Minute
	}

	// Create the claims for the token
	claims := dtos.CustomTokenClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   userId,
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(duration)),
		},
		TokenType: string(tokenType),
	}

	// Create the token
	token := jwt.NewWithClaims(
		jwt.SigningMethodHS256,
		claims,
	)

	// Get the secret key of the token type
	var secretKey string
	switch tokenType {
	case types.RefreshToken:
		secretKey = cfg.RefreshTokenSecretKey
	case types.AccessToken:
		secretKey = cfg.AccessTokenSecretKey
	default:
		err := fmt.Errorf("invalid token type: %s", tokenType)
		return "", err
	}

	// Sign the token
	tokenString, tokenSigningErr := token.SignedString([]byte(secretKey))
	if tokenSigningErr != nil {
		tokenSigningErr = fmt.Errorf("Failed to sign the token: %w", tokenSigningErr)
		return "", tokenSigningErr
	}

	return tokenString, nil
}
