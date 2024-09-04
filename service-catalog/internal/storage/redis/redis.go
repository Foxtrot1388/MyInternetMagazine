package cashe

import (
	"context"
	"encoding/json"
	"v1/internal/entity"
	liberrors "v1/internal/lib/errors"

	"github.com/redis/go-redis/v9"
	"go.opentelemetry.io/otel/trace"
)

type RedisCache struct {
	cli    *redis.Client
	tracer trace.Tracer
}

var ErrRecordNotFound = redis.Nil

func New(redisHost string, tracer trace.Tracer) (*RedisCache, error) {
	const op = "redis.new"

	rdb := redis.NewClient(&redis.Options{
		Addr:     redisHost,
		Password: "", // no password set
		DB:       0,  // use default DB
	})
	_, err := rdb.Ping(context.Background()).Result()
	return &RedisCache{cli: rdb, tracer: tracer}, liberrors.WrapErr(op, err)
}

func (r *RedisCache) Get(ctx context.Context, key string) (*entity.Product, error) {
	const op = "redis.get"

	ctxspan, span := r.tracer.Start(ctx, "redis_get")
	defer span.End()

	result, err := r.cli.Get(ctxspan, key).Result()
	if err != nil {
		return nil, liberrors.WrapErr(op, err)
	}

	var v entity.Product
	err = json.Unmarshal([]byte(result), &v)
	if err != nil {
		return nil, liberrors.WrapErr(op, err)
	}

	return &v, nil

}

func (r *RedisCache) Set(ctx context.Context, key string, v *entity.Product) error {
	const op = "redis.set"

	ctxspan, span := r.tracer.Start(ctx, "redis_set")
	defer span.End()

	res, err := json.Marshal(v)
	if err != nil {
		return liberrors.WrapErr(op, err)
	}

	err = r.cli.Set(ctxspan, key, res, 0).Err()
	if err != nil {
		return liberrors.WrapErr(op, err)
	}

	return nil

}

func (r *RedisCache) Invalidate(ctx context.Context, key string) error {
	const op = "redis.invalidate"

	ctxspan, span := r.tracer.Start(ctx, "redis_invalidate")
	defer span.End()

	_, err := r.cli.Del(ctxspan, key).Result()
	if err != nil {
		return liberrors.WrapErr(op, err)
	}

	return nil

}
