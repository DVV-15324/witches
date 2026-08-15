package utils

import (
	"context"
	"fmt"
	"time"

	"github.com/DVV-15324/witches/cmd/utils"
	"github.com/redis/go-redis/v9"
)

type BlacklistService struct {
	redis  *redis.Client
	config *utils.Config
}

func NewBlacklistService(redis *redis.Client, config *utils.Config) *BlacklistService {
	return &BlacklistService{redis: redis, config: config}
}

func (s *BlacklistService) cacheKeyBlacklist(accessToken string) string {
	return fmt.Sprintf("blacklist:%s", accessToken)
}

func (s *BlacklistService) BlacklistToken(ctx context.Context, accessToken string) error {
	if accessToken == "" {
		return fmt.Errorf("Error: accessToken cannot be empty")
	}

	key := s.cacheKeyBlacklist(accessToken)
	return s.redis.Set(ctx, key, "revoked", time.Duration(s.config.RevokedTTL)*time.Second).Err()
}

// Việc Check nên để ở middleware
func (s *BlacklistService) IsTokenBlacklisted(ctx context.Context, token string) bool {
	if token == "" {
		return false
	}

	key := s.cacheKeyBlacklist(token)
	exists, err := s.redis.Exists(ctx, key).Result()
	return err == nil && exists == 1
}
