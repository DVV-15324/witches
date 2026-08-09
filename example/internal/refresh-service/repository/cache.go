package repository

import (
	"context"
	"encoding/json"
	"errors"
	//entity "example/internal/entity/refresh"
	"time"
	entityRefresh "example/internal/refresh-service/entity"

	"github.com/redis/go-redis/v9"
)

func (r *RefreshTokenRepository) cacheToken(ctx context.Context, token *entityRefresh.RefreshToken) error {
	if token == nil {
		return nil
	}

	// Cache by token
	data, err := json.Marshal(token)
	if err != nil {
		return err
	}
	key := r.cacheKeyToken(token.Token)
	if err := r.redis.Set(ctx, key, data, 5*time.Minute).Err(); err != nil {
		return err
	}

	// Cache by user_id (cho active tokens)
	if !token.Revoked && time.Now().Before(time.Unix(token.ExpiresAt, 0)) {
		r.redis.SAdd(ctx, r.cacheKeyUserTokens(token.UserID), token.Token)
	}

	return nil
}

func (r *RefreshTokenRepository) getCachedToken(ctx context.Context, token string) (*entityRefresh.RefreshToken, error) {
	key := r.cacheKeyToken(token)
	data, err := r.redis.Get(ctx, key).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return nil, nil
		}
		return nil, err
	}

	var rt entityRefresh.RefreshToken
	if err := json.Unmarshal([]byte(data), &rt); err != nil {
		return nil, err
	}
	return &rt, nil
}

func (r *RefreshTokenRepository) deleteCacheToken(ctx context.Context, token string) error {
	key := r.cacheKeyToken(token)
	return r.redis.Del(ctx, key).Err()
}

func (r *RefreshTokenRepository) invalidateUserCache(ctx context.Context, token string) {
	// Lấy user_id từ cache để invalidate
	data, err := r.redis.Get(ctx, r.cacheKeyToken(token)).Result()
	if err == nil {
		var rt entityRefresh.RefreshToken
		if json.Unmarshal([]byte(data), &rt) == nil {
			r.redis.Del(ctx, r.cacheKeyActiveTokens(rt.UserID))
		}
	}
}

func (r *RefreshTokenRepository) cacheNotFound(ctx context.Context, token string) {
	key := r.cacheKeyNotFound(token)
	r.redis.Set(ctx, key, "1", 1*time.Minute)
}

func (r *RefreshTokenRepository) isNotFoundCached(ctx context.Context, token string) bool {
	key := r.cacheKeyNotFound(token)
	exists, _ := r.redis.Exists(ctx, key).Result()
	return exists == 1
}
