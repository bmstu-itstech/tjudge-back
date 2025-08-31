package services_test

import (
	"testing"
	"time"

	"github.com/bmstu-itstech/tjudge-back/internal/domain/tjudge"
	"github.com/bmstu-itstech/tjudge-back/internal/services"
	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)


func TestTokenContent(t *testing.T) {
	secretKey := "test-secret-key"
	tokenDuration := 30 * time.Minute
	userID := tjudge.UserID(12)

	generator := services.NewJWTTokenGenerator(secretKey, tokenDuration)
	token, err := generator.Generate(userID)
	require.NoError(t, err)
	assert.NotEmpty(t, token)

	parsedToken, err := jwt.ParseWithClaims(
		string(token), 
		&services.Claims{}, 
		func(token *jwt.Token) (interface{}, error) {
			return []byte(secretKey), nil
		},
	)
	require.NoError(t, err)
	assert.True(t, parsedToken.Valid)

	claims, ok := parsedToken.Claims.(*services.Claims)
	require.True(t, ok)
	assert.Equal(t, userID, claims.UserID)
}

func TestTokenSecret(t *testing.T) {
	secretKey := "test-secret-key"
	invalidSecretKey := "invalid-secret-key"
	tokenDuration := 30 * time.Minute
	userID := tjudge.UserID(12)

	generator := services.NewJWTTokenGenerator(secretKey, tokenDuration)
	token, err := generator.Generate(userID)
	require.NoError(t, err)

	parsedToken, err := jwt.ParseWithClaims(
		string(token), 
		&services.Claims{}, 
		func(token *jwt.Token) (interface{}, error) {
			return []byte(invalidSecretKey), nil
		},
	)
	require.Error(t, err)
	assert.False(t, parsedToken.Valid)
}




