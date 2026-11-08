package repository

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"e-wallet-system-api/domains/JWTVerification"
	"e-wallet-system-api/domains/JWTVerification/entities"
	"e-wallet-system-api/infrastructure/cache"

	"github.com/redis/go-redis/v9"
)

type serverJWTRepository struct {
	cache cache.Cache
}

func NewServerJWTRepository(cache cache.Cache) JWTVerification.ServerJWTRepository {
	return serverJWTRepository{cache: cache}
}

// CreateServerJWT implements [JWTVerification.ServerJWTRepository].
func (s serverJWTRepository) CreateServerJWT(ctx context.Context, accessToken string, userId string) (*entities.ServerJWT, error) {
	user := &entities.ServerJWT{
		Id:        accessToken,
		UserId:    userId,
		Revoked:   false,
		ExpiresAt: time.Now().Add(1 * time.Hour),
	}

	data, err := json.Marshal(user)
	if err != nil {
		return nil, err
	}

	result := s.cache.GetInstance().Set(ctx, "UserId:"+userId, data, time.Hour*1)
	if result.Err() != nil {
		return nil, result.Err()
	}

	return user, nil
}

// DeleteServerJWT implements [JWTVerification.ServerJWTRepository].
func (s serverJWTRepository) DeleteServerJWT(ctx context.Context, userId string) error {
	result := s.cache.GetInstance().Del(ctx, "UserId:"+userId)
	if result.Err() != nil {
		return result.Err()
	}
	return nil
}

// FindServerJWTbyUserID implements [JWTVerification.ServerJWTRepository].
func (s serverJWTRepository) FindServerJWTbyUserID(ctx context.Context, userId string) (*entities.ServerJWT, error) {
	var user *entities.ServerJWT
	result, err := s.cache.GetInstance().Get(ctx, "UserId:"+userId).Result()
	if errors.Is(err, redis.Nil) {
		return nil, err
	}
	err = json.Unmarshal([]byte(result), &user)
	return user, nil
}

// UpdateServerJWT implements [JWTVerification.ServerJWTRepository].
func (s serverJWTRepository) UpdateServerJWT(ctx context.Context, ServerJWT *entities.ServerJWT) error {
	var user *entities.ServerJWT
	result, err := s.cache.GetInstance().Get(ctx, "UserId:"+ServerJWT.UserId).Result()
	if errors.Is(err, redis.Nil) {
		return err
	}
	err = json.Unmarshal([]byte(result), &user)

	user.Revoked = true

	data, _ := json.Marshal(user)
	s.cache.GetInstance().Set(ctx, "UserId:"+ServerJWT.UserId, data, redis.KeepTTL)

	return nil
}
