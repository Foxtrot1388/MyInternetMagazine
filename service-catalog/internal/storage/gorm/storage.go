package storage

import (
	"context"
	postgresgorm "gorm.io/driver/postgres"
	"gorm.io/gorm"
	"v1/internal/entity"
	"v1/internal/lib"
)

type Storage struct {
	db *gorm.DB
}

func New(connection string) (*Storage, error) {
	const op = "storage.new"

	db, err := gorm.Open(postgresgorm.Open(connection), &gorm.Config{})
	if err != nil {
		return nil, lib.WrapErr(op, err)
	}
	return &Storage{db: db.Table("entity")}, nil
}

func (s *Storage) Get(ctx context.Context, id int) (*entity.Product, error) {
	const op = "storage.get"

	var product entity.Product
	err := s.db.WithContext(ctx).Take(&product, id).Error
	if err != nil {
		return nil, lib.WrapErr(op, err)
	}

	return &product, nil

}

func (s *Storage) List(ctx context.Context) (*[]entity.ElementOfList, error) {
	const op = "storage.list"

	var products []entity.ElementOfList
	err := s.db.WithContext(ctx).Find(&products).Error
	if err != nil {
		return nil, lib.WrapErr(op, err)
	}

	return &products, nil

}

func (s *Storage) Create(ctx context.Context, product *entity.Product) (int, error) {
	const op = "storage.create"

	err := s.db.WithContext(ctx).Create(&product).Error

	if err != nil {
		return 0, lib.WrapErr(op, err)
	} else {
		return product.Id, nil
	}

}

func (s *Storage) Delete(ctx context.Context, id int) (bool, error) {
	const op = "storage.delete"

	type DeleteProductRequest struct {
		Id int
	}

	err := s.db.WithContext(ctx).Delete(&DeleteProductRequest{Id: id}).Error
	if err != nil {
		return false, lib.WrapErr(op, err)
	}

	return true, nil

}
