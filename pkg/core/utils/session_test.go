// pkg/core/utils/session_test.go
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

func setupSessionTest(t *testing.T) (*SessionService, *miniredis.Miniredis) {
	mr, err := miniredis.Run()
	require.NoError(t, err)

	client := redis.NewClient(&redis.Options{
		Addr: mr.Addr(),
	})

	service := NewSessionService(client, 3600, 300)
	return service, mr
}

func TestNewSessionService(t *testing.T) {
	client, mr := setupTestRedis(t)
	defer mr.Close()

	service := NewSessionService(client, 3600, 300)
	assert.NotNil(t, service)
	assert.Equal(t, int64(3600), service.SessionTTL)
	assert.Equal(t, int64(300), service.IdleTimeout)
}

func TestSessionService_CreateAndGetSession(t *testing.T) {
	service, mr := setupSessionTest(t)
	defer mr.Close()

	ctx := context.Background()
	session := &SessionCache{
		UserID:    1,
		Email:     "test@example.com",
		Role:      "user",
		DeviceID:  "device-123",
		IPAddress: "192.168.1.1",
		UserAgent: "Mozilla/5.0",
		Locale:    "vi-VN",
		Timezone:  "Asia/Ho_Chi_Minh",
	}

	// Create session
	err := service.CreateSession(ctx, session)
	require.NoError(t, err)
	assert.NotEmpty(t, session.SessionID)

	// Get session
	retrieved, err := service.GetSession(ctx, session.UserID, session.DeviceID)
	require.NoError(t, err)
	assert.NotNil(t, retrieved)
	assert.Equal(t, session.UserID, retrieved.UserID)
	assert.Equal(t, session.Email, retrieved.Email)
	assert.Equal(t, session.SessionID, retrieved.SessionID)
}

func TestSessionService_DeleteSession(t *testing.T) {
	service, mr := setupSessionTest(t)
	defer mr.Close()

	ctx := context.Background()
	session := &SessionCache{
		UserID:   1,
		Email:    "test@example.com",
		DeviceID: "device-123",
	}

	err := service.CreateSession(ctx, session)
	require.NoError(t, err)

	// Delete session
	err = service.DeleteSession(ctx, session.UserID, session.DeviceID)
	require.NoError(t, err)

	// Verify deleted
	retrieved, err := service.GetSession(ctx, session.UserID, session.DeviceID)
	require.NoError(t, err)
	assert.Nil(t, retrieved)
}

func TestSessionService_UpdateSession(t *testing.T) {
	service, mr := setupSessionTest(t)
	defer mr.Close()

	ctx := context.Background()
	session := &SessionCache{
		UserID:      1,
		Email:       "test@example.com",
		DeviceID:    "device-123",
		AccessToken: "old-token",
	}

	err := service.CreateSession(ctx, session)
	require.NoError(t, err)

	// Update session
	newToken := "new-token"
	err = service.UpdateSession(ctx, session.UserID, session.DeviceID, newToken)
	require.NoError(t, err)

	// Verify update
	retrieved, err := service.GetSession(ctx, session.UserID, session.DeviceID)
	require.NoError(t, err)
	assert.Equal(t, newToken, retrieved.AccessToken)
}

func TestSessionService_IsSessionIdle(t *testing.T) {
	service, mr := setupSessionTest(t)
	defer mr.Close()

	ctx := context.Background()
	session := &SessionCache{
		UserID:     1,
		Email:      "test@example.com",
		DeviceID:   "device-123",
		LastActive: time.Now().Unix(),
	}

	err := service.CreateSession(ctx, session)
	require.NoError(t, err)

	isIdle, err := service.IsSessionIdle(ctx, session.UserID, session.DeviceID)
	require.NoError(t, err)
	assert.False(t, isIdle)

	// Update last active to old time
	session.LastActive = time.Now().Unix() - 400 // 400 seconds ago
	err = service.CreateSession(ctx, session)
	require.NoError(t, err)

	isIdle, err = service.IsSessionIdle(ctx, session.UserID, session.DeviceID)
	require.NoError(t, err)
	assert.True(t, isIdle)
}

func TestSessionService_ValidateSession(t *testing.T) {
	service, mr := setupSessionTest(t)
	defer mr.Close()

	ctx := context.Background()
	session := &SessionCache{
		UserID:      1,
		Email:       "test@example.com",
		DeviceID:    "device-123",
		AccessToken: "valid-token",
		LastActive:  time.Now().Unix(),
	}

	err := service.CreateSession(ctx, session)
	require.NoError(t, err)

	// Valid session
	validated, err := service.ValidateSession(ctx, session.UserID, session.DeviceID, "valid-token")
	require.NoError(t, err)
	assert.NotNil(t, validated)

	// Invalid token
	_, err = service.ValidateSession(ctx, session.UserID, session.DeviceID, "invalid-token")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "token mismatch")

	// Invalid device
	_, err = service.ValidateSession(ctx, session.UserID, "wrong-device", "valid-token")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "session not found")
}

// pkg/core/utils/session_test.go
func TestSessionService_UpdateLastActive(t *testing.T) {
	service, mr := setupSessionTest(t)
	defer mr.Close()

	ctx := context.Background()
	session := &SessionCache{
		UserID:     1,
		Email:      "test@example.com",
		DeviceID:   "device-123",
		LastActive: time.Now().Unix() - 100,
	}

	err := service.CreateSession(ctx, session)
	require.NoError(t, err)

	// Update last active
	err = service.UpdateLastActive(ctx, session.UserID, session.DeviceID)
	require.NoError(t, err)

	// Verify update
	retrieved, err := service.GetSession(ctx, session.UserID, session.DeviceID)
	require.NoError(t, err)
	assert.Greater(t, retrieved.LastActive, session.LastActive)
}

func TestSessionService_UpdateLastActive_SessionNotFound(t *testing.T) {
	service, mr := setupSessionTest(t)
	defer mr.Close()

	ctx := context.Background()
	err := service.UpdateLastActive(ctx, 999, "non-existent-device")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "session not found")
}

func TestSessionService_CreateSession_WithNilSession(t *testing.T) {
	service, mr := setupSessionTest(t)
	defer mr.Close()

	ctx := context.Background()
	err := service.CreateSession(ctx, nil)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "session cannot be nil")
}

func TestSessionService_GetSession_NotFound(t *testing.T) {
	service, mr := setupSessionTest(t)
	defer mr.Close()

	ctx := context.Background()
	session, err := service.GetSession(ctx, 999, "non-existent-device")
	assert.NoError(t, err)
	assert.Nil(t, session)
}
