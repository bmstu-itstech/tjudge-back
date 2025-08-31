package services

import (
	"fmt"
	"time"

	"github.com/bmstu-itstech/tjudge-back/internal/domain/tjudge"
	"github.com/golang-jwt/jwt/v5"
)

type JWTTokenGenerator struct {
	secretKey     []byte
	tokenDuration time.Duration
}

func NewJWTTokenGenerator(secretKey string, tokenDuration time.Duration) *JWTTokenGenerator {
	return &JWTTokenGenerator{
		secretKey:     []byte(secretKey),
		tokenDuration: tokenDuration,
	}
}

type Claims struct {
	jwt.RegisteredClaims
	UserID tjudge.UserID `json:"user_id"`
}

func (g *JWTTokenGenerator) Generate(id tjudge.UserID) (tjudge.Token, error) {
	expirationTime := time.Now().Add(g.tokenDuration)
	
	claims := &Claims{
		UserID: id,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(g.secretKey)
	if err != nil {
		return "", fmt.Errorf("failed to generate token: %w", err)
	}

	return tjudge.Token(tokenString), nil
}