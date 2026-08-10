package utils

import (
	"context"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewJwtService(t *testing.T) {
	service := NewJwtService("secret", 3600, 86400)
	assert.NotNil(t, service)
	assert.Equal(t, int64(3600), service.accessTokenExpiry)
	assert.Equal(t, int64(86400), service.refreshTokenExpiry)
}

func TestJwtService_IssueAccessToken(t *testing.T) {
	service := NewJwtService("secret", 3600, 86400)
	ctx := context.Background()

	token, err := service.IssueAccessToken(ctx, "user-123", "trace-456")
	require.NoError(t, err)
	assert.NotNil(t, token)
	assert.NotEmpty(t, token.Token)
	assert.Equal(t, int64(3600), token.ExpireAt)
}

func TestJwtService_IssueRefreshToken(t *testing.T) {
	service := NewJwtService("secret", 3600, 86400)
	ctx := context.Background()

	token, err := service.IssueRefreshToken(ctx, "user-123", "trace-456")
	require.NoError(t, err)
	assert.NotNil(t, token)
	assert.NotEmpty(t, token.Token)
	assert.Equal(t, int64(86400), token.ExpireAt)
}

func TestJwtService_IssueTokenPair(t *testing.T) {
	service := NewJwtService("secret", 3600, 86400)
	ctx := context.Background()

	pair, err := service.IssueTokenPair(ctx, "user-123", "trace-456")
	require.NoError(t, err)
	assert.NotNil(t, pair)
	assert.NotEmpty(t, pair.AccessToken.Token)
	assert.NotEmpty(t, pair.RefreshToken.Token)
	assert.NotEqual(t, pair.AccessToken.Token, pair.RefreshToken.Token)
}

func TestJwtService_ParseToken_Valid(t *testing.T) {
	service := NewJwtService("secret", 3600, 86400)
	ctx := context.Background()

	sub := "user-123"
	tid := "trace-456"

	token, err := service.IssueAccessToken(ctx, sub, tid)
	require.NoError(t, err)

	claims, err := service.ParseToken(ctx, token.Token)
	require.NoError(t, err)
	assert.NotNil(t, claims)
	assert.Equal(t, sub, claims.Subject)
	assert.Equal(t, tid, claims.ID)
}

func TestJwtService_ParseToken_Expired(t *testing.T) {
	// Create service with negative expiry
	service := NewJwtService("secret", -1, 86400)
	ctx := context.Background()

	token, err := service.IssueAccessToken(ctx, "user-123", "trace-456")
	require.NoError(t, err)

	_, err = service.ParseToken(ctx, token.Token)
	assert.ErrorIs(t, err, ErrTokenExpired)
}

func TestJwtService_ParseToken_Invalid(t *testing.T) {
	service := NewJwtService("secret", 3600, 86400)
	ctx := context.Background()

	tests := []struct {
		name  string
		token string
		err   error
	}{
		{
			name:  "malformed token",
			token: "invalid.token.here",
			err:   ErrMalformedToken,
		},
		{
			name:  "empty token",
			token: "",
			err:   ErrMalformedToken,
		},
		{
			name:  "wrong signature",
			token: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ.SflKxwRJSMeKKF2QT4fwpMeJf36POk6yJV_adQssw5c",
			err:   ErrInvalidSignature,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := service.ParseToken(ctx, tt.token)
			assert.ErrorIs(t, err, tt.err)
		})
	}
}
func TestJwtService_ParseToken_WithInvalidAlgorithm(t *testing.T) {
	service := NewJwtService("secret", 3600, 86400)
	ctx := context.Background()

	// Tạo token với algorithm không hợp lệ
	token := jwt.NewWithClaims(jwt.SigningMethodNone, jwt.MapClaims{
		"sub": "user-123",
		"exp": time.Now().Add(time.Hour).Unix(),
	})
	tokenStr, err := token.SignedString(jwt.UnsafeAllowNoneSignatureType)
	require.NoError(t, err)

	_, err = service.ParseToken(ctx, tokenStr)
	assert.Error(t, err)
}

func TestJwtService_ParseToken_WithFutureNotBefore(t *testing.T) {
	service := NewJwtService("secret", 3600, 86400)
	ctx := context.Background()

	claims := JwtClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   "user-123",
			NotBefore: jwt.NewNumericDate(time.Now().Add(time.Hour)),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(2 * time.Hour)),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenStr, err := token.SignedString([]byte("secret"))
	require.NoError(t, err)

	_, err = service.ParseToken(ctx, tokenStr)
	assert.ErrorIs(t, err, ErrTokenNotYetValid)
}

func TestJwtService_IssueTokenPair_WithError(t *testing.T) {
	// Test khi IssueAccessToken bị lỗi (khó mock, có thể skip)
	t.Skip("Skipping error test - requires complex mock")
}
