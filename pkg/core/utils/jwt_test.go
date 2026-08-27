package utils

import (
	"context"
	"testing"
	"time"

	"github.com/DVV-15324/witches/cmd/utils"
	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func getTestConfigJWT() *utils.Config {
	return &utils.Config{
		JWTSecret:       "test-secret-key",
		AccessTokenTTL:  3600,
		RefreshTokenTTL: 86400,
	}
}

func TestNewJwtService(t *testing.T) {
	cfg := getTestConfigJWT()
	service := NewJwtService(cfg)
	assert.NotNil(t, service)
	assert.Equal(t, cfg, service.config)
}

func TestJwtService_IssueAccessToken(t *testing.T) {
	cfg := getTestConfigJWT()
	service := NewJwtService(cfg)
	ctx := context.Background()

	token, err := service.IssueAccessToken(ctx, "user-123", "trace-456", "jkt")
	require.NoError(t, err)
	assert.NotNil(t, token)
	assert.NotEmpty(t, token.Token)
	assert.Equal(t, cfg.AccessTokenTTL, token.ExpireAt)
}

func TestJwtService_IssueRefreshToken(t *testing.T) {
	cfg := getTestConfigJWT()
	service := NewJwtService(cfg)
	ctx := context.Background()

	token, err := service.IssueRefreshToken(ctx, "user-123", "trace-456")
	require.NoError(t, err)
	assert.NotNil(t, token)
	assert.NotEmpty(t, token.Token)
	assert.Equal(t, cfg.RefreshTokenTTL, token.ExpireAt)
}

func TestJwtService_IssueTokenPair(t *testing.T) {
	cfg := getTestConfigJWT()
	service := NewJwtService(cfg)
	ctx := context.Background()

	pair, err := service.IssueTokenPair(ctx, "user-123", "trace-456", "jkt")
	require.NoError(t, err)
	assert.NotNil(t, pair)
	assert.NotEmpty(t, pair.AccessToken.Token)
	assert.NotEmpty(t, pair.RefreshToken.Token)
	assert.NotEqual(t, pair.AccessToken.Token, pair.RefreshToken.Token)
}

func TestJwtService_ParseToken_Valid(t *testing.T) {
	cfg := getTestConfigJWT()
	service := NewJwtService(cfg)
	ctx := context.Background()

	token, err := service.IssueAccessToken(ctx, "user-123", "trace-456", "jkt")
	require.NoError(t, err)

	claims, err := service.ParseToken(ctx, token.Token)
	require.NoError(t, err)
	assert.Equal(t, "user-123", claims.Subject)
	assert.Equal(t, "trace-456", claims.ID)
}

func TestJwtService_ParseToken_Expired(t *testing.T) {
	cfg := getTestConfigJWT()
	cfg.AccessTokenTTL = -1
	service := NewJwtService(cfg)
	ctx := context.Background()

	token, err := service.IssueAccessToken(ctx, "user-123", "trace-456", "jkt")
	require.NoError(t, err)
	time.Sleep(10 * time.Millisecond)

	_, err = service.ParseToken(ctx, token.Token)
	assert.ErrorIs(t, err, ErrTokenExpired)
}

func TestJwtService_ParseToken_Invalid(t *testing.T) {
	cfg := getTestConfigJWT()
	service := NewJwtService(cfg)
	ctx := context.Background()

	tests := []struct {
		name  string
		token string
		err   error
	}{
		{"malformed token", "invalid.token.here", ErrMalformedToken},
		{"empty token", "", ErrMalformedToken},
		{"wrong signature", "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ.SflKxwRJSMeKKF2QT4fwpMeJf36POk6yJV_adQssw5c", ErrInvalidSignature},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := service.ParseToken(ctx, tt.token)
			assert.ErrorIs(t, err, tt.err)
		})
	}
}

func TestJwtService_ParseToken_WithInvalidAlgorithm(t *testing.T) {
	cfg := getTestConfigJWT()
	service := NewJwtService(cfg)
	ctx := context.Background()

	// Tạo token với algorithm none
	token := jwt.NewWithClaims(jwt.SigningMethodNone, jwt.MapClaims{
		"sub": "user-123",
		"exp": time.Now().Add(time.Hour).Unix(),
	})
	tokenStr, err := token.SignedString(jwt.UnsafeAllowNoneSignatureType)
	require.NoError(t, err)

	_, err = service.ParseToken(ctx, tokenStr)
	assert.ErrorIs(t, err, ErrInvalidAlgorithm)
}

func TestJwtService_ParseToken_WithFutureNotBefore(t *testing.T) {
	cfg := getTestConfigJWT()
	service := NewJwtService(cfg)
	ctx := context.Background()

	claims := JwtClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   "user-123",
			NotBefore: jwt.NewNumericDate(time.Now().Add(time.Hour)),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(2 * time.Hour)),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenStr, err := token.SignedString([]byte(cfg.JWTSecret))
	require.NoError(t, err)

	_, err = service.ParseToken(ctx, tokenStr)
	assert.ErrorIs(t, err, ErrTokenNotYetValid)
}

func TestJwtService_ParseToken_WithConfigNil(t *testing.T) {
	service := &JwtService{config: nil}
	ctx := context.Background()
	_, err := service.ParseToken(ctx, "any.token")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "jwt service config is nil")
}

func TestJwtService_IssueAccessToken_WithNilConfig(t *testing.T) {
	service := &JwtService{config: nil}
	ctx := context.Background()
	_, err := service.IssueAccessToken(ctx, "sub", "tid", "jkt")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "jwt service config is nil")
}
