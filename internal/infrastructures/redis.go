package infrastructures

import (
	"context"
	"fmt"

	"github.com/redis/go-redis/v9"
)

func NewRedis(addr, password string, db int) *redis.Client {
	client := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       db,
	})

	ctx := context.Background()
	if err := client.Ping(ctx).Err(); err != nil {
		panic(fmt.Errorf("error connecting to redis: %+v", err.Error()))
	}

	return client
}

func NewRedisFromUrl(url string) *redis.Client {
	opt, err := redis.ParseURL(url)
	if err != nil {
		panic(err)
	}

	return redis.NewClient(opt)
}
