package redis

import (
	"context"
	"fmt"
	"time"

	"github.com/aalexanderkevin/crypto-wallet/config"
	"github.com/aalexanderkevin/crypto-wallet/service"

	"github.com/redis/go-redis/v9"
)

type RedisClient struct {
	Client *redis.Client
}

func NewRedis(cfg config.Redis) service.Cache {
	redisOpt := &redis.Options{
		Addr: fmt.Sprintf("%s:%d", cfg.Host, cfg.Port),
	}

	client := redis.NewClient(redisOpt)
	return &RedisClient{Client: client}
}

func (r *RedisClient) Close() {
	r.Client.Close()
}

func (r *RedisClient) Get(ctx context.Context, key string) (string, error) {
	res, err := r.Client.Get(ctx, key).Result()
	if err != nil && err.Error() == "redis: nil" {
		return "", fmt.Errorf("redis nil")
	} else if err != nil {
		return "", err
	}

	return res, nil
}

func (r *RedisClient) Put(ctx context.Context, key string, data string, ttl time.Duration) error {
	return r.Client.Set(ctx, key, data, ttl).Err()
}

func (r *RedisClient) SetList(ctx context.Context, key string, data string, ttl *time.Duration) (int64, error) {
	res, err := r.Client.RPush(ctx, key, data).Result()
	if err != nil {
		return 0, err
	}

	err = r.Client.Expire(ctx, key, *ttl).Err()
	return res, err
}

func (r *RedisClient) GetList(ctx context.Context, key string) ([]string, error) {
	rangeResult, err := r.Client.LRange(ctx, key, 0, -1).Result()
	if err != nil {
		return nil, err
	}

	return rangeResult, nil
}

func (r *RedisClient) Delete(ctx context.Context, key string) error {
	return r.Client.Del(ctx, key).Err()
}
