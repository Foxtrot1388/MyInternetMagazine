package grpcapi

import (
	"context"
	"v1/internal/model"
)

type Service interface {
	Login(ctx context.Context, pass, login string) (*model.LoginUser, error)
	Get(ctx context.Context, id int) (*model.User, error)
	Create(ctx context.Context, login, pass, fname, sname, lname, email string) (int, error)
	Delete(ctx context.Context, id int) (bool, error)
}
