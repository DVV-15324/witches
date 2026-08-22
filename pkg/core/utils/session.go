package utils

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"time"
)

type SessionCache struct {
	SessionID   string `json:"session_id"`
	UserID      int64  `json:"user_id"`
	DeviceID    string `json:"device_id"`
	AccessToken string `json:"access_token"`
	LastActive  int64  `json:"last_active"`
	Locale      string `json:"locale"`
	Timezone    string `json:"timezone"`
}

type SessionService struct {
	redis           *redis.Client
	RefreshTokenTTL int64
	IdleTimeout     int64
}

func NewSessionService(redis *redis.Client, RefreshTokenTTL int64, IdleTimeout int64) *SessionService {
	return &SessionService{redis: redis, IdleTimeout: IdleTimeout, RefreshTokenTTL: RefreshTokenTTL}
}

func (s *SessionService) cacheKeySession(userID int64, deviceID string) string {
	return fmt.Sprintf("session:userID:%d:deviceID:%s", userID, deviceID)
}

func (s *SessionService) CreateSession(ctx context.Context, session *SessionCache) error {
	if session == nil {
		return fmt.Errorf("session cannot be nil")
	}
	session.SessionID = uuid.New().String()
	key := s.cacheKeySession(session.UserID, session.DeviceID)
	data, err := json.Marshal(session)
	if err != nil {
		return err
	}
	ttl := time.Duration(s.RefreshTokenTTL) * time.Second
	return s.redis.Set(ctx, key, data, ttl).Err()
}

func (s *SessionService) GetSession(ctx context.Context, userID int64, deviceID string) (*SessionCache, error) {
	key := s.cacheKeySession(userID, deviceID)
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

func (s *SessionService) DeleteSession(ctx context.Context, userID int64, deviceID string) error {
	key := s.cacheKeySession(userID, deviceID)
	return s.redis.Del(ctx, key).Err()
}

func (s *SessionService) UpdateSession(ctx context.Context, userID int64, deviceID string, accessToken string) error {
	session, err := s.GetSession(ctx, userID, deviceID)
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

func (s *SessionService) UpdateLastActive(ctx context.Context, userID int64, deviceID string) error {
	session, err := s.GetSession(ctx, userID, deviceID)
	if err != nil {
		return err
	}
	if session == nil {
		return fmt.Errorf("session not found for user %d, device %s", userID, deviceID)
	}
	session.LastActive = time.Now().Unix()
	return s.CreateSession(ctx, session)
}

func (s *SessionService) IsSessionIdle(ctx context.Context, userID int64, deviceID string) (bool, error) {
	session, err := s.GetSession(ctx, userID, deviceID)
	if err != nil {
		return false, err
	}
	if session == nil {
		return true, nil
	}
	now := time.Now().Unix()
	return (now - session.LastActive) > s.IdleTimeout, nil
}

func (s *SessionService) ValidateSession(ctx context.Context, userID int64, deviceID string, accessToken string) (*SessionCache, error) {
	session, err := s.GetSession(ctx, userID, deviceID)
	if err != nil {
		return nil, err
	}
	if session == nil {
		return nil, fmt.Errorf("session not found")
	}

	if session.AccessToken != accessToken {
		return nil, fmt.Errorf("token mismatch")
	}

	if session.DeviceID != deviceID {
		return nil, fmt.Errorf("device mismatch")
	}

	isIdle, err := s.IsSessionIdle(ctx, userID, deviceID)
	if err != nil {
		return nil, err
	}

	if isIdle {
		return nil, fmt.Errorf("session idle timeout")
	}

	return session, nil
}
