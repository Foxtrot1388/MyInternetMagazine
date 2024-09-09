package httpapi

import (
	"context"
	"v1/internal/model"
)

type Service interface {
	Get(ctx context.Context, id int) (*model.Product, error)
	Create(ctx context.Context, name string, description string) (int, error)
	Delete(ctx context.Context, id int) (bool, error)
	List(ctx context.Context) (*[]model.ElementOfList, error)
}
