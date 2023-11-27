package service

import (
	"context"
	"log/slog"
	"strconv"
	"v1/internal/entity"
	"v1/internal/lib"
)

type Service struct {
	DB    DBRepository
	Cashe CasheRepository
	log   *slog.Logger
}

type CasheRepository interface {
	Get(ctx context.Context, key string) (*entity.Product, error)
	Set(ctx context.Context, key string, v *entity.Product) error
	Invalidate(ctx context.Context, key string) error
}

type DBRepository interface {
	Get(ctx context.Context, id int) (*entity.Product, error)
	Create(ctx context.Context, product *entity.Product) (int, error)
	Delete(ctx context.Context, id int) (bool, error)
	List(ctx context.Context) (*[]entity.ElementOfList, error)
}

func New(log *slog.Logger, DB DBRepository, cashe CasheRepository) *Service {
	return &Service{DB: DB, log: log, Cashe: cashe}
}

func (s *Service) Get(ctx context.Context, id int) (*entity.Product, error) {
	const op = "service.get"

	log := s.log.With(
		slog.String("op", op),
	)

	product, err := s.Cashe.Get(ctx, strconv.Itoa(id))
	if err != nil {

		product, err = s.DB.Get(ctx, id)
		if err != nil {
			log.Error(err.Error())
			return nil, lib.WrapErr(op, err)
		}

		err = s.Cashe.Set(ctx, strconv.Itoa(id), product)
		if err != nil {
			log.Error(err.Error())
			return nil, lib.WrapErr(op, err)
		}

	}

	return product, nil

}

func (s *Service) Create(ctx context.Context, product *entity.Product) (int, error) {
	const op = "service.create"

	log := s.log.With(
		slog.String("op", op),
	)

	id, err := s.DB.Create(ctx, product)
	if err != nil {
		log.Error(err.Error())
		return 0, lib.WrapErr(op, err)
	}

	return id, nil

}

func (s *Service) Delete(ctx context.Context, id int) (bool, error) {
	const op = "service.delete"

	log := s.log.With(
		slog.String("op", op),
	)

	result, err := s.DB.Delete(ctx, id)
	if err != nil {
		log.Error(err.Error())
		return false, lib.WrapErr(op, err)
	}

	err = s.Cashe.Invalidate(ctx, strconv.Itoa(id))
	if err != nil {
		log.Error(err.Error())
	}

	return result, nil

}

func (s *Service) List(ctx context.Context) (*[]entity.ElementOfList, error) {
	const op = "service.list"

	log := s.log.With(
		slog.String("op", op),
	)

	result, err := s.DB.List(ctx)
	if err != nil {
		log.Error(err.Error())
		return nil, lib.WrapErr(op, err)
	}

	return result, nil

}
