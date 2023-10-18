package service

import (
	"context"
	"time"
)

type Cache interface {
	Close()
	Get(ctx context.Context, key string) (string, error)
	Put(context.Context, string, string, time.Duration) error
	Delete(context.Context, string) error
	SetList(ctx context.Context, key string, data string, ttl *time.Duration) (int64, error)
	GetList(ctx context.Context, key string) ([]string, error)
}
