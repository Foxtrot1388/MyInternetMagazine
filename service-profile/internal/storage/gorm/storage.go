package storage

import (
	"context"
	"v1/internal/entity"
	liberrors "v1/internal/lib/errors"
	"v1/internal/model"

	"go.opentelemetry.io/otel"
	postgresgorm "gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var tracer = otel.Tracer("profile-server")

type Storage struct {
	db *gorm.DB
}

func New(connection string) (*Storage, error) {
	const op = "gorm.new"

	db, err := gorm.Open(postgresgorm.Open(connection), &gorm.Config{})
	if err != nil {
		return nil, liberrors.WrapErr(op, err)
	}
	return &Storage{db: db.Table("users")}, nil
}

func (s *Storage) Login(ctx context.Context, login string) (*entity.UserDB, error) {
	const op = "gorm.login"

	ctxspan, span := tracer.Start(ctx, "gorm_login")
	defer span.End()

	var user entity.UserDB
	err := s.db.WithContext(ctxspan).Where("login = ?", login).First(&user).Error
	if err != nil {
		return nil, liberrors.WrapErr(op, err)
	}

	return &user, nil

}

func (s *Storage) Get(ctx context.Context, id int) (*entity.User, error) {
	const op = "gorm.get"

	ctxspan, span := tracer.Start(ctx, "gorm_get")
	defer span.End()

	var user entity.User
	err := s.db.WithContext(ctxspan).Take(&user, id).Error
	if err != nil {
		return nil, liberrors.WrapErr(op, err)
	}

	return &user, nil

}

func (s *Storage) Create(ctx context.Context, user *model.NewUser) (int, error) {
	const op = "gorm.create"

	ctxspan, span := tracer.Start(ctx, "gorm_create")
	defer span.End()

	err := s.db.WithContext(ctxspan).Create(&user).Error

	if err != nil {
		return 0, liberrors.WrapErr(op, err)
	} else {
		return user.Id, nil
	}

}

func (s *Storage) Delete(ctx context.Context, id int) (bool, error) {
	const op = "gorm.delete"

	ctxspan, span := tracer.Start(ctx, "gorm_delete")
	defer span.End()

	type DeleteUserRequest struct {
		Id int
	}

	err := s.db.WithContext(ctxspan).Delete(&DeleteUserRequest{Id: id}).Error
	if err != nil {
		return false, liberrors.WrapErr(op, err)
	}

	return true, nil

}
