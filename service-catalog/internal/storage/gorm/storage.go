package storage

import (
	"context"
	"v1/internal/entity"
	liberrors "v1/internal/lib/errors"

	"go.opentelemetry.io/otel"

	postgresgorm "gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var tracer = otel.Tracer("catalog-server")

type Storage struct {
	db *gorm.DB
}

func New(connection string) (*Storage, error) {
	const op = "storage.new"

	db, err := gorm.Open(postgresgorm.Open(connection), &gorm.Config{})
	if err != nil {
		return nil, liberrors.WrapErr(op, err)
	}
	return &Storage{db: db.Table("entity")}, nil
}

func (s *Storage) Get(ctx context.Context, id int) (*entity.Product, error) {
	const op = "storage.get"

	ctxspan, span := tracer.Start(ctx, "gorm_get")
	defer span.End()

	var product entity.Product
	err := s.db.WithContext(ctxspan).Take(&product, id).Error
	if err != nil {
		return nil, liberrors.WrapErr(op, err)
	}

	return &product, nil

}

func (s *Storage) List(ctx context.Context) (*[]entity.ElementOfList, error) {
	const op = "storage.list"

	ctxspan, span := tracer.Start(ctx, "gorm_list")
	defer span.End()

	var products []entity.ElementOfList
	err := s.db.WithContext(ctxspan).Find(&products).Error
	if err != nil {
		return nil, liberrors.WrapErr(op, err)
	}

	return &products, nil

}

func (s *Storage) Create(ctx context.Context, product *entity.Product) (int, error) {
	const op = "storage.create"

	ctxspan, span := tracer.Start(ctx, "gorm_create")
	defer span.End()

	err := s.db.WithContext(ctxspan).Create(&product).Error

	if err != nil {
		return 0, liberrors.WrapErr(op, err)
	} else {
		return product.Id, nil
	}

}

func (s *Storage) Delete(ctx context.Context, id int) (bool, error) {
	const op = "storage.delete"

	ctxspan, span := tracer.Start(ctx, "gorm_delete")
	defer span.End()

	type DeleteProductRequest struct {
		Id int
	}

	err := s.db.WithContext(ctxspan).Delete(&DeleteProductRequest{Id: id}).Error
	if err != nil {
		return false, liberrors.WrapErr(op, err)
	}

	return true, nil

}
