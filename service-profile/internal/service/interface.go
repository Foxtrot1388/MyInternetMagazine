package service

import (
	"context"
	"v1/internal/entity"
	"v1/internal/model"
)

type Repository interface {
	Login(ctx context.Context, login string) (*entity.UserDB, error)
	Get(ctx context.Context, id int) (*entity.User, error)
	Create(ctx context.Context, user *model.NewUser) (int, error)
	Delete(ctx context.Context, id int) (bool, error)
}

type KafkaSender interface {
	Send(v []byte, topic string) error
}
