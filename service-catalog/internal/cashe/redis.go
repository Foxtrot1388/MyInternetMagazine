package cashe

import (
	"context"
	"encoding/json"
	"github.com/redis/go-redis/v9"
	entity "v1/internal"
)

type RedisCache struct {
	cli *redis.Client
}

func New(redisHost string) (*RedisCache, error) {
	rdb := redis.NewClient(&redis.Options{
		Addr:     redisHost,
		Password: "", // no password set
		DB:       0,  // use default DB
	})
	_, err := rdb.Ping(context.Background()).Result()
	return &RedisCache{cli: rdb}, err
}

func (r *RedisCache) Get(ctx context.Context, key string) (*entity.Product, error) {

	result, err := r.cli.Get(ctx, key).Result()
	if err != nil {
		return nil, err
	}

	var v entity.Product
	err = json.Unmarshal([]byte(result), &v)
	if err != nil {
		return nil, err
	}

	return &v, nil

}

func (r *RedisCache) Set(ctx context.Context, key string, v *entity.Product) error {

	res, err := json.Marshal(v)
	if err != nil {
		return err
	}

	err = r.cli.Set(ctx, key, res, 0).Err()
	if err != nil {
		return err
	}

	return nil

}
