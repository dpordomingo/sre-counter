package counter

import (
	"context"

	redis "github.com/go-redis/redis/v8"
)

var ctx = context.Background()

// CacheClient modelates a Cache Service
type CacheClient interface {
	Get(key string) (string, error)
	Set(key, value string) error
}

type redisClient struct {
	client *redis.Client
}

// NewRedisClient returns a CacheClient backed in Redis
func NewRedisClient(clientAddr string) CacheClient {
	return &redisClient{
		client: redis.NewClient(&redis.Options{
			Addr:     clientAddr,
			Password: "", // no password set
			DB:       0,  // use default DB
		}),
	}
}

func (r *redisClient) Set(key, value string) error {
	return r.client.Set(ctx, key, value, 0).Err()
}

func (r *redisClient) Get(key string) (string, error) {
	val, err := r.client.Get(ctx, key).Result()
	if err == redis.Nil {
		return "", nil
	} else if err != nil {
		return "", err
	} else {
		return val, nil
	}
}
