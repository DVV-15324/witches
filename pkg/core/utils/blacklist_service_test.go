package utils

import (
	"context"
	"github.com/alicebob/miniredis/v2"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func setupTestRedis(t *testing.T) (*redis.Client, *miniredis.Miniredis) {
	mr, err := miniredis.Run()
	require.NoError(t, err)

	client := redis.NewClient(&redis.Options{
		Addr: mr.Addr(),
	})

	return client, mr
}

// Blacklist Tests

func TestNewBlacklistService(t *testing.T) {
	client, mr := setupTestRedis(t)
	defer mr.Close()

	cfg := getTestConfig()
	service := NewBlacklistService(client, cfg)
	assert.NotNil(t, service)
	assert.Equal(t, cfg.RevokedTTL, service.config.RevokedTTL)
}

func TestBlacklistService_BlacklistToken(t *testing.T) {
	client, mr := setupTestRedis(t)
	defer mr.Close()

	cfg := getTestConfig()
	service := NewBlacklistService(client, cfg)
	ctx := context.Background()

	tests := []struct {
		name        string
		token       string
		expectError bool
	}{
		{
			name:        "valid token",
			token:       "valid-token-123",
			expectError: false,
		},
		{
			name:        "empty token",
			token:       "",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := service.BlacklistToken(ctx, tt.token)
			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)

				// Verify token was blacklisted
				key := service.cacheKeyBlacklist(tt.token)
				val, err := client.Get(ctx, key).Result()
				assert.NoError(t, err)
				assert.Equal(t, "revoked", val)

				// Check TTL
				ttl, err := client.TTL(ctx, key).Result()
				assert.NoError(t, err)
				assert.Greater(t, ttl, time.Duration(0))
				assert.LessOrEqual(t, ttl, time.Duration(cfg.RevokedTTL+1)*time.Second)
			}
		})
	}
}

func TestBlacklistService_IsTokenBlacklisted(t *testing.T) {
	client, mr := setupTestRedis(t)
	defer mr.Close()

	cfg := getTestConfig()
	service := NewBlacklistService(client, cfg)
	ctx := context.Background()

	t.Run("token is blacklisted", func(t *testing.T) {
		token := "blacklisted-token"
		err := service.BlacklistToken(ctx, token)
		require.NoError(t, err)

		isBlacklisted := service.IsTokenBlacklisted(ctx, token)
		assert.True(t, isBlacklisted)
	})

	t.Run("token is not blacklisted", func(t *testing.T) {
		token := "valid-token"
		isBlacklisted := service.IsTokenBlacklisted(ctx, token)
		assert.False(t, isBlacklisted)
	})

	t.Run("empty token", func(t *testing.T) {
		isBlacklisted := service.IsTokenBlacklisted(ctx, "")
		assert.False(t, isBlacklisted)
	})
}

func TestBlacklistService_BlacklistToken_WithTTL(t *testing.T) {
	client, mr := setupTestRedis(t)
	defer mr.Close()

	cfg := getTestConfig()
	cfg.RevokedTTL = 1 // 1 giây
	service := NewBlacklistService(client, cfg)
	ctx := context.Background()

	token := "expiring-token"
	err := service.BlacklistToken(ctx, token)
	require.NoError(t, err)

	// Check token exists
	isBlacklisted := service.IsTokenBlacklisted(ctx, token)
	assert.True(t, isBlacklisted)

	// Kiểm tra TTL đã được set
	key := service.cacheKeyBlacklist(token)
	ttl, err := client.TTL(ctx, key).Result()
	assert.NoError(t, err)
	assert.Greater(t, ttl, time.Duration(0))
	assert.LessOrEqual(t, ttl, time.Second*2)

	// Đợi TTL hết
	mr.FastForward(2 * time.Second)

	// Token không còn trong blacklist
	isBlacklisted = service.IsTokenBlacklisted(ctx, token)
	assert.False(t, isBlacklisted)
}

func TestBlacklistService_CacheKeyBlacklist(t *testing.T) {
	client, mr := setupTestRedis(t)
	defer mr.Close()

	cfg := getTestConfig()
	service := NewBlacklistService(client, cfg)

	token := "test-token"
	key := service.cacheKeyBlacklist(token)
	assert.Equal(t, "blacklist:test-token", key)
}

// Integration Test

func TestJWTAndBlacklist_Integration(t *testing.T) {
	// 1. Setup Redis
	client, mr := setupTestRedis(t)
	defer mr.Close()

	// 2. Setup config
	cfg := getTestConfig()
	cfg.AccessTokenTTL = 3600
	cfg.RevokedTTL = 300

	// 3. Create services
	jwtService := NewJwtService(cfg)
	blacklistService := NewBlacklistService(client, cfg)
	ctx := context.Background()

	// 4. Issue token
	tokenPair, err := jwtService.IssueTokenPair(ctx, "user-123", "trace-456")
	require.NoError(t, err)
	accessToken := tokenPair.AccessToken.Token

	// 5. Parse token (valid)
	claims, err := jwtService.ParseToken(ctx, accessToken)
	require.NoError(t, err)
	assert.Equal(t, "user-123", claims.Subject)

	// 6. Blacklist token
	err = blacklistService.BlacklistToken(ctx, accessToken)
	require.NoError(t, err)

	// 7. Check blacklist
	isBlacklisted := blacklistService.IsTokenBlacklisted(ctx, accessToken)
	assert.True(t, isBlacklisted)

	// 8. Parse token again (should still be valid structurally, but middleware will check blacklist)
	claims2, err := jwtService.ParseToken(ctx, accessToken)
	require.NoError(t, err)
	assert.Equal(t, "user-123", claims2.Subject)
}
