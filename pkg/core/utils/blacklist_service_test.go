package utils

import (
	"context"
	"testing"
	"time"

	"github.com/alicebob/miniredis/v2"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

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

	service := NewBlacklistService(client, 3600)
	assert.NotNil(t, service)
	assert.Equal(t, int64(3600), service.revokedTTL)
}

func TestBlacklistService_BlacklistToken(t *testing.T) {
	client, mr := setupTestRedis(t)
	defer mr.Close()

	service := NewBlacklistService(client, 3600)
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
			}
		})
	}
}

func TestBlacklistService_IsTokenBlacklisted(t *testing.T) {
	client, mr := setupTestRedis(t)
	defer mr.Close()

	service := NewBlacklistService(client, 3600)
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

	// TTL = 1 second
	service := NewBlacklistService(client, 1)
	ctx := context.Background()

	token := "expiring-token"
	err := service.BlacklistToken(ctx, token)
	require.NoError(t, err)

	// Check token exists
	isBlacklisted := service.IsTokenBlacklisted(ctx, token)
	assert.True(t, isBlacklisted)

	// Kiểm tra TTL đã được set (không cần đợi)
	key := service.cacheKeyBlacklist(token)
	ttl, err := client.TTL(ctx, key).Result()
	assert.NoError(t, err)
	assert.Greater(t, ttl, time.Duration(0))
	assert.LessOrEqual(t, ttl, time.Second*2)
}
