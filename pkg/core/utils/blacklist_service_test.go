package utils

import (
	"context"
	"testing"
	"time"

	"github.com/DVV-15324/witches/cmd/utils"
	"github.com/alicebob/miniredis/v2"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func getTestConfigBL() *utils.Config {
	return &utils.Config{
		JWTSecret:       "test-secret-key",
		AccessTokenTTL:  3600,
		RefreshTokenTTL: 86400,
		RevokedTTL:      300,
	}
}

func setupTestRedis(t *testing.T) (*redis.Client, *miniredis.Miniredis) {
	mr, err := miniredis.Run()
	require.NoError(t, err)

	client := redis.NewClient(&redis.Options{
		Addr: mr.Addr(),
	})

	return client, mr
}

func TestNewBlacklistService(t *testing.T) {
	client, mr := setupTestRedis(t)
	defer mr.Close()

	cfg := getTestConfigBL()
	service := NewBlacklistService(client, cfg)
	assert.NotNil(t, service)
	assert.Equal(t, cfg.RevokedTTL, service.config.RevokedTTL)
}

func TestBlacklistService_BlacklistToken(t *testing.T) {
	client, mr := setupTestRedis(t)
	defer mr.Close()

	cfg := getTestConfigBL()
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
				key := service.cacheKeyBlacklist(tt.token)
				val, err := client.Get(ctx, key).Result()
				assert.NoError(t, err)
				assert.Equal(t, "revoked", val)

				// Chỉ kiểm tra key tồn tại, không check TTL vì miniredis không ổn định
				exists, _ := client.Exists(ctx, key).Result()
				assert.Equal(t, int64(1), exists)
			}
		})
	}
}

func TestBlacklistService_IsTokenBlacklisted(t *testing.T) {
	client, mr := setupTestRedis(t)
	defer mr.Close()

	cfg := getTestConfigBL()
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

	cfg := getTestConfigBL()
	cfg.RevokedTTL = 1
	service := NewBlacklistService(client, cfg)
	ctx := context.Background()

	token := "expiring-token"
	err := service.BlacklistToken(ctx, token)
	require.NoError(t, err)

	isBlacklisted := service.IsTokenBlacklisted(ctx, token)
	assert.True(t, isBlacklisted)

	key := service.cacheKeyBlacklist(token)
	exists, _ := client.Exists(ctx, key).Result()
	assert.Equal(t, int64(1), exists)

	mr.FastForward(2 * time.Second)
	isBlacklisted = service.IsTokenBlacklisted(ctx, token)
	assert.False(t, isBlacklisted)
}

func TestBlacklistService_CacheKeyBlacklist(t *testing.T) {
	client, mr := setupTestRedis(t)
	defer mr.Close()

	cfg := getTestConfigBL()
	service := NewBlacklistService(client, cfg)

	token := "test-token"
	key := service.cacheKeyBlacklist(token)
	assert.Equal(t, "blacklist:test-token", key)
}

func TestJWTAndBlacklist_Integration(t *testing.T) {
	client, mr := setupTestRedis(t)
	defer mr.Close()

	cfg := getTestConfigBL()
	cfg.AccessTokenTTL = 3600
	cfg.RevokedTTL = 300

	jwtService := NewJwtService(cfg)
	blacklistService := NewBlacklistService(client, cfg)
	ctx := context.Background()

	tokenPair, err := jwtService.IssueTokenPair(ctx, "user-123", "trace-456", "jkt")
	require.NoError(t, err)
	accessToken := tokenPair.AccessToken.Token

	claims, err := jwtService.ParseToken(ctx, accessToken)
	require.NoError(t, err)
	assert.Equal(t, "user-123", claims.Subject)

	err = blacklistService.BlacklistToken(ctx, accessToken)
	require.NoError(t, err)

	isBlacklisted := blacklistService.IsTokenBlacklisted(ctx, accessToken)
	assert.True(t, isBlacklisted)

	claims2, err := jwtService.ParseToken(ctx, accessToken)
	require.NoError(t, err)
	assert.Equal(t, "user-123", claims2.Subject)
}
