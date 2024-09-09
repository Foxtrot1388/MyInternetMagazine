package service

import (
	"context"
	"v1/internal/entity"
)

//go:generate mockery --name CasheRepository
type CasheRepository interface {
	Get(ctx context.Context, key string) (*entity.Product, error)
	Set(ctx context.Context, key string, v *entity.Product) error
	Invalidate(ctx context.Context, key string) error
}

//go:generate mockery --name DBRepository
type DBRepository interface {
	Get(ctx context.Context, id int) (*entity.Product, error)
	Create(ctx context.Context, product *entity.Product) (int, error)
	Delete(ctx context.Context, id int) (bool, error)
	List(ctx context.Context) (*[]entity.ElementOfList, error)
}
