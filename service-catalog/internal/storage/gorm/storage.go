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
	return &Storage{db: db.Table("entity")}, nil
}

func (s *Storage) Get(ctx context.Context, id int) (*entity.Product, error) {

	var product entity.Product
	err := s.db.WithContext(ctx).Take(&product, id).Error
	if err != nil {
		return nil, err
	}

	return &product, nil

}

func (s *Storage) List(ctx context.Context) (*[]entity.ElementOfList, error) {

	var products []entity.ElementOfList
	err := s.db.WithContext(ctx).Find(&products).Error
	if err != nil {
		return nil, err
	}

	return &products, nil

}

func (s *Storage) Create(ctx context.Context, product *entity.Product) (int, error) {

	err := s.db.WithContext(ctx).Create(&product).Error

	if err != nil {
		return 0, err
	} else {
		return product.Id, nil
	}

}

func (s *Storage) Delete(ctx context.Context, id int) (bool, error) {

	type DeleteProductRequest struct {
		Id int
	}

	err := s.db.WithContext(ctx).Delete(&DeleteProductRequest{Id: id}).Error
	if err != nil {
		return false, err
	}

	return true, nil

}
