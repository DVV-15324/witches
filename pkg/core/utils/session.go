// utils/session.go
package utils

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"time"

	"github.com/redis/go-redis/v9"
)

type SessionConfig struct {
	SessionTTL  int64
	IdleTimeout int64 // seconds - thời gian inactive tối đa
}

type SessionCache struct {
	SessionID   string `json:"session_id"`
	UserID      uint32 `json:"user_id"`
	Email       string `json:"email"`
	Role        string `json:"role"`
	DeviceID    string `json:"device_id"`
	IPAddress   string `json:"ip_address"`
	UserAgent   string `json:"user_agent"`
	AccessToken string `json:"access_token"`
	Locale      string `json:"locale"`
	Timezone    string `json:"timezone"`
	LoginAt     int64  `json:"login_at"`
	LastActive  int64  `json:"last_active"`
}

type SessionInfo struct {
	SessionID  string `json:"session_id"`
	DeviceID   string `json:"device_id"`
	IPAddress  string `json:"ip_address"`
	Locale     string `json:"locale"`
	Timezone   string `json:"timezone"`
	LoginAt    int64  `json:"login_at"`
	LastActive int64  `json:"last_active"`
	ExpiresAt  int64  `json:"expires_at"`
	IsActive   bool   `json:"is_active"`
	IsIdle     bool   `json:"is_idle"`
}

type SessionService struct {
	redis       *redis.Client
	SessionTTL  int64
	IdleTimeout int64
}

func NewSessionService(redis *redis.Client, SessionTTL int64,
	IdleTimeout int64) *SessionService {
	return &SessionService{
		redis:       redis,
		IdleTimeout: IdleTimeout,
		SessionTTL:  SessionTTL,
	}
}

func (s *SessionService) cacheKeySession(userID uint32, deviceID string) string {
	return fmt.Sprintf("session:user:%d:device:%s", userID, deviceID)
}

// CreateSession - Tạo session mới
func (s *SessionService) CreateSession(ctx context.Context, session *SessionCache) error {
	if session == nil {
		return fmt.Errorf("session cannot be nil")
	}
	session.SessionID = uuid.New().String()

	// Thêm deviceID vào key
	key := s.cacheKeySession(session.UserID, session.DeviceID)
	data, err := json.Marshal(session)
	if err != nil {
		return err
	}

	ttl := time.Duration(s.SessionTTL) * time.Second
	return s.redis.Set(ctx, key, data, ttl).Err()
}

func (s *SessionService) GetSession(ctx context.Context, userID uint32, deviceID string) (*SessionCache, error) {
	key := s.cacheKeySession(userID, deviceID) // Thêm deviceID
	data, err := s.redis.Get(ctx, key).Result()
	if err != nil {
		if err == redis.Nil {
			return nil, nil
		}
		return nil, err
	}

	var session SessionCache
	if err := json.Unmarshal([]byte(data), &session); err != nil {
		return nil, err
	}

	return &session, nil
}

// DeleteSession - Xóa session theo userID + deviceID
func (s *SessionService) DeleteSession(ctx context.Context, userID uint32, deviceID string) error {
	key := s.cacheKeySession(userID, deviceID) // Thêm deviceID
	return s.redis.Del(ctx, key).Err()
}

// UpdateSession - Update session khi refresh token
func (s *SessionService) UpdateSession(ctx context.Context, userID uint32, deviceID string, accessToken string) error {
	session, err := s.GetSession(ctx, userID, deviceID) // Thêm deviceID
	if err != nil {
		return err
	}
	if session == nil {
		return fmt.Errorf("session not found for user %d, device %s", userID, deviceID)
	}

	session.AccessToken = accessToken
	session.LastActive = time.Now().Unix()

	return s.CreateSession(ctx, session)
}

// UpdateLastActive - Update last active time
func (s *SessionService) UpdateLastActive(ctx context.Context, userID uint32, deviceID string) error {
	session, err := s.GetSession(ctx, userID, deviceID) // Thêm deviceID
	if err != nil {
		return err
	}
	if session == nil {
		return fmt.Errorf("session not found for user %d, device %s", userID, deviceID)
	}

	session.LastActive = time.Now().Unix()
	return s.CreateSession(ctx, session)
}

// IsSessionIdle - Kiểm tra session có bị idle không
func (s *SessionService) IsSessionIdle(ctx context.Context, userID uint32, deviceID string) (bool, error) {
	session, err := s.GetSession(ctx, userID, deviceID) //  Thêm deviceID
	if err != nil {
		return false, err
	}
	if session == nil {
		return true, nil
	}

	now := time.Now().Unix()
	return (now - session.LastActive) > s.IdleTimeout, nil
}

// ValidateSession - Validate session (check tồn tại, token match, idle)
// // Việc Check nên để ở middleware
func (s *SessionService) ValidateSession(ctx context.Context, userID uint32, deviceID string, accessToken string) (*SessionCache, error) {
	session, err := s.GetSession(ctx, userID, deviceID) //  Thêm deviceID
	if err != nil {
		return nil, err
	}
	if session == nil {
		return nil, fmt.Errorf("session not found")
	}

	// Check token match
	if session.AccessToken != accessToken {
		return nil, fmt.Errorf("token mismatch")
	}
	// Check device match
	if session.DeviceID != deviceID {
		return nil, fmt.Errorf("device mismatch")
	}

	// Check idle
	isIdle, err := s.IsSessionIdle(ctx, userID, deviceID)
	if err != nil {
		return nil, err
	}
	if isIdle {
		return nil, fmt.Errorf("session idle timeout")
	}

	return session, nil
}
