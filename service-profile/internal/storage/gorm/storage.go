package storage

import (
	"context"
	postgresgorm "gorm.io/driver/postgres"
	"gorm.io/gorm"
	entity "v1/internal"
)

type Storage struct {
	db *gorm.DB
}

func New(connection string) (*Storage, error) {
	db, err := gorm.Open(postgresgorm.Open(connection), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	return &Storage{db: db.Table("users")}, nil
}

func (s *Storage) Login(ctx context.Context, pass, login string) (*entity.LoginUser, error) {

	var user entity.LoginUser
	err := s.db.WithContext(ctx).Where("pass = ? AND login = ?", pass, login).First(&user).Error
	if err != nil {
		return nil, err
	}

	return &user, nil

}

func (s *Storage) Get(ctx context.Context, id int) (*entity.User, error) {

	var user entity.User
	err := s.db.WithContext(ctx).Take(&user, id).Error
	if err != nil {
		return nil, err
	}

	return &user, nil

}

func (s *Storage) Create(ctx context.Context, user *entity.NewUser) (int, error) {

	err := s.db.WithContext(ctx).Create(&user).Error

	if err != nil {
		return 0, err
	} else {
		return user.Id, nil
	}

}

func (s *Storage) Delete(ctx context.Context, id int) (bool, error) {

	type DeleteUserRequest struct {
		Id int
	}

	err := s.db.WithContext(ctx).Delete(&DeleteUserRequest{Id: id}).Error
	if err != nil {
		return false, err
	}

	return true, nil

}
