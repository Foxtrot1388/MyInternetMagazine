package cashe

import (
	"context"
	"encoding/json"
	"github.com/redis/go-redis/v9"
	"v1/internal/entity"
	"v1/internal/lib"
)

type RedisCache struct {
	cli *redis.Client
}

func New(redisHost string) (*RedisCache, error) {
	const op = "redis.new"

	rdb := redis.NewClient(&redis.Options{
		Addr:     redisHost,
		Password: "", // no password set
		DB:       0,  // use default DB
	})
	_, err := rdb.Ping(context.Background()).Result()
	return &RedisCache{cli: rdb}, lib.WrapErr(op, err)
}

func (r *RedisCache) Get(ctx context.Context, key string) (*entity.Product, error) {
	const op = "redis.get"

	result, err := r.cli.Get(ctx, key).Result()
	if err != nil {
		return nil, lib.WrapErr(op, err)
	}

	var v entity.Product
	err = json.Unmarshal([]byte(result), &v)
	if err != nil {
		return nil, lib.WrapErr(op, err)
	}

	return &v, nil

}

func (r *RedisCache) Set(ctx context.Context, key string, v *entity.Product) error {
	const op = "redis.set"

	res, err := json.Marshal(v)
	if err != nil {
		return lib.WrapErr(op, err)
	}

	err = r.cli.Set(ctx, key, res, 0).Err()
	if err != nil {
		return lib.WrapErr(op, err)
	}

	return nil

}
