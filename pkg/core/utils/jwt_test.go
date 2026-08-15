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

func getTestConfig() *utils.Config {
	return &utils.Config{
		JWTSecret:       "test-secret-key",
		AccessTokenTTL:  3600,  // 1 giờ
		RefreshTokenTTL: 86400, // 24 giờ
	}
}

func TestNewJwtService(t *testing.T) {
	cfg := getTestConfig()
	service := NewJwtService(cfg)
	assert.NotNil(t, service)
	assert.Equal(t, cfg, service.config)
}

func TestJwtService_IssueAccessToken(t *testing.T) {
	cfg := getTestConfig()
	service := NewJwtService(cfg)
	ctx := context.Background()

	token, err := service.IssueAccessToken(ctx, "user-123", "trace-456")
	require.NoError(t, err)
	assert.NotNil(t, token)
	assert.NotEmpty(t, token.Token)
	assert.Equal(t, cfg.AccessTokenTTL, token.ExpireAt)
}

func TestJwtService_IssueRefreshToken(t *testing.T) {
	cfg := getTestConfig()
	service := NewJwtService(cfg)
	ctx := context.Background()

	token, err := service.IssueRefreshToken(ctx, "user-123", "trace-456")
	require.NoError(t, err)
	assert.NotNil(t, token)
	assert.NotEmpty(t, token.Token)
	assert.Equal(t, cfg.RefreshTokenTTL, token.ExpireAt)
}

func TestJwtService_IssueTokenPair(t *testing.T) {
	cfg := getTestConfig()
	service := NewJwtService(cfg)
	ctx := context.Background()

	pair, err := service.IssueTokenPair(ctx, "user-123", "trace-456")
	require.NoError(t, err)
	assert.NotNil(t, pair)
	assert.NotEmpty(t, pair.AccessToken.Token)
	assert.NotEmpty(t, pair.RefreshToken.Token)
	assert.NotEqual(t, pair.AccessToken.Token, pair.RefreshToken.Token)
	assert.Equal(t, cfg.AccessTokenTTL, pair.AccessToken.ExpireAt)
	assert.Equal(t, cfg.RefreshTokenTTL, pair.RefreshToken.ExpireAt)
}

func TestJwtService_ParseToken_Valid(t *testing.T) {
	cfg := getTestConfig()
	service := NewJwtService(cfg)
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
	// Create config with negative TTL
	cfg := getTestConfig()
	cfg.AccessTokenTTL = -1 // token hết hạn ngay lập tức
	service := NewJwtService(cfg)
	ctx := context.Background()

	token, err := service.IssueAccessToken(ctx, "user-123", "trace-456")
	require.NoError(t, err)

	// Chờ một chút để token hết hạn
	time.Sleep(10 * time.Millisecond)

	_, err = service.ParseToken(ctx, token.Token)
	assert.ErrorIs(t, err, ErrTokenExpired)
}

func TestJwtService_ParseToken_Invalid(t *testing.T) {
	cfg := getTestConfig()
	service := NewJwtService(cfg)
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
	cfg := getTestConfig()
	service := NewJwtService(cfg)
	ctx := context.Background()

	// Tạo token với algorithm none (không an toàn)
	token := jwt.NewWithClaims(jwt.SigningMethodNone, jwt.MapClaims{
		"sub": "user-123",
		"exp": time.Now().Add(time.Hour).Unix(),
	})
	tokenStr, err := token.SignedString(jwt.UnsafeAllowNoneSignatureType)
	require.NoError(t, err)

	_, err = service.ParseToken(ctx, tokenStr)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "unexpected signing method")
}

func TestJwtService_ParseToken_WithFutureNotBefore(t *testing.T) {
	cfg := getTestConfig()
	service := NewJwtService(cfg)
	ctx := context.Background()

	// Tạo token với NotBefore trong tương lai
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
	// Trường hợp config bị nil (có thể xảy ra nếu không khởi tạo đúng)
	service := &JwtService{config: nil}
	ctx := context.Background()

	// Tạo token hợp lệ nhưng service không có config
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": "user-123",
		"exp": time.Now().Add(time.Hour).Unix(),
	})
	tokenStr, err := token.SignedString([]byte("secret"))
	require.NoError(t, err)

	_, err = service.ParseToken(ctx, tokenStr)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "parse token failed")
}

func TestJwtService_ParseToken_WithMissingClaims(t *testing.T) {
	cfg := getTestConfig()
	service := NewJwtService(cfg)
	ctx := context.Background()

	// Tạo token không có subject (thiếu claim)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"exp": time.Now().Add(time.Hour).Unix(),
	})
	tokenStr, err := token.SignedString([]byte(cfg.JWTSecret))
	require.NoError(t, err)

	claims, err := service.ParseToken(ctx, tokenStr)
	require.NoError(t, err)
	assert.NotNil(t, claims)
	assert.Empty(t, claims.Subject) // subject rỗng
}

func BenchmarkJwtService_IssueAccessToken(b *testing.B) {
	cfg := getTestConfig()
	service := NewJwtService(cfg)
	ctx := context.Background()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = service.IssueAccessToken(ctx, "user-123", "trace-456")
	}
}

func BenchmarkJwtService_ParseToken(b *testing.B) {
	cfg := getTestConfig()
	service := NewJwtService(cfg)
	ctx := context.Background()

	token, err := service.IssueAccessToken(ctx, "user-123", "trace-456")
	if err != nil {
		b.Fatal(err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = service.ParseToken(ctx, token.Token)
	}
}
