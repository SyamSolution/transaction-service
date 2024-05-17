package config

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/redis/go-redis/v9"
)

type Cache struct {
	db         *redis.Client
	service    string
	defaultExp time.Duration
}

type Cacher interface {
	Set(ctx context.Context, key string, value interface{}, duration time.Duration) error
	SetNX(ctx context.Context, key string, value interface{}, duration time.Duration) error
	Get(ctx context.Context, key string) (string, error)
	Del(ctx context.Context, key string) error
}

func NewCacher(logger Logger) Cacher {
	address := fmt.Sprintf("%s:%s", os.Getenv("CACHER_HOST"), os.Getenv("CACHER_PORT"))
	duration, _ := time.ParseDuration(os.Getenv("CACHER_DEFAULT_EXP"))
	cacher := redis.NewClient(&redis.Options{
		Addr:     address,
		Password: os.Getenv("CACHER_PASSWORD"),
		DB:       0,
	})

	return &Cache{
		db:         cacher,
		service:    os.Getenv("CACHER_SERVICE"),
		defaultExp: duration,
	}
}

func (c *Cache) Set(ctx context.Context, key string, value interface{}, duration time.Duration) error {
	fullKey := fmt.Sprintf("%s:%s", c.service, key)
	if duration == 0*time.Second {
		duration = 0
	}

	if err := c.db.Set(ctx, fullKey, value, duration).Err(); err != nil {
		return err
	}

	return nil
}

func (c *Cache) SetNX(ctx context.Context, key string, value interface{}, duration time.Duration) error {
	fullKey := fmt.Sprintf("%s:%s", c.service, key)
	if duration == 0*time.Second {
		duration = 0
	}

	if err := c.db.SetNX(ctx, fullKey, value, duration).Err(); err != nil {
		return err
	}

	return nil
}

func (c *Cache) Get(ctx context.Context, key string) (string, error) {
	fullKey := fmt.Sprintf("%s:%s", c.service, key)
	value, err := c.db.Get(ctx, fullKey).Result()
	if err != nil {
		return "", err
	}

	return value, nil
}

func (c *Cache) Del(ctx context.Context, key string) error {
	fullKey := fmt.Sprintf("%s:%s", c.service, key)
	_, err := c.db.Del(ctx, fullKey).Result()
	if err != nil {
		return err
	}

	return nil
}
