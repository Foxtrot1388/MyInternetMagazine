package storage

import (
	"context"
	postgresgorm "gorm.io/driver/postgres"
	"gorm.io/gorm"
	entity "v1/internal"
	"v1/internal/lib"
)

type Storage struct {
	db *gorm.DB
}

func New(connection string) (*Storage, error) {
	const op = "gorm.new"

	db, err := gorm.Open(postgresgorm.Open(connection), &gorm.Config{})
	if err != nil {
		return nil, lib.WrapErr(op, err)
	}
	return &Storage{db: db.Table("users")}, nil
}

func (s *Storage) Login(ctx context.Context, pass, login string) (*entity.LoginUser, error) {
	const op = "gorm.login"

	var user entity.LoginUser
	err := s.db.WithContext(ctx).Where("pass = ? AND login = ?", pass, login).First(&user).Error
	if err != nil {
		return nil, lib.WrapErr(op, err)
	}

	return &user, nil

}

func (s *Storage) Get(ctx context.Context, id int) (*entity.User, error) {
	const op = "gorm.get"

	var user entity.User
	err := s.db.WithContext(ctx).Take(&user, id).Error
	if err != nil {
		return nil, lib.WrapErr(op, err)
	}

	return &user, nil

}

func (s *Storage) Create(ctx context.Context, user *entity.NewUser) (int, error) {
	const op = "gorm.create"

	err := s.db.WithContext(ctx).Create(&user).Error

	if err != nil {
		return 0, lib.WrapErr(op, err)
	} else {
		return user.Id, nil
	}

}

func (s *Storage) Delete(ctx context.Context, id int) (bool, error) {
	const op = "gorm.delete"

	type DeleteUserRequest struct {
		Id int
	}

	err := s.db.WithContext(ctx).Delete(&DeleteUserRequest{Id: id}).Error
	if err != nil {
		return false, lib.WrapErr(op, err)
	}

	return true, nil

}
