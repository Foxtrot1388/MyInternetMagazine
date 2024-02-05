package service

import (
	"context"
	"log/slog"
	"strconv"
	"v1/internal/entity"
	liberrors "v1/internal/lib/errors"

	"go.opentelemetry.io/otel"
)

var tracer = otel.Tracer("catalog-server")

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

	ctxspan, span := tracer.Start(ctx, "service_get")
	defer span.End()

	log := s.log.With(
		slog.String("op", op),
	)

	product, err := s.Cashe.Get(ctxspan, strconv.Itoa(id))
	if err != nil {

		product, err = s.DB.Get(ctxspan, id)
		if err != nil {
			log.Error(err.Error())
			return nil, liberrors.WrapErr(op, err)
		}

		err = s.Cashe.Set(ctxspan, strconv.Itoa(id), product)
		if err != nil {
			log.Error(err.Error())
			return nil, liberrors.WrapErr(op, err)
		}

	}

	return product, nil

}

func (s *Service) Create(ctx context.Context, name string, description string) (int, error) {
	const op = "service.create"

	ctxspan, span := tracer.Start(ctx, "service_create")
	defer span.End()

	log := s.log.With(
		slog.String("op", op),
	)

	product := &entity.Product{
		Name:        name,
		Description: description,
	}

	id, err := s.DB.Create(ctxspan, product)
	if err != nil {
		log.Error(err.Error())
		return 0, liberrors.WrapErr(op, err)
	}

	return id, nil

}

func (s *Service) Delete(ctx context.Context, id int) (bool, error) {
	const op = "service.delete"

	ctxspan, span := tracer.Start(ctx, "service_delete")
	defer span.End()

	log := s.log.With(
		slog.String("op", op),
	)

	result, err := s.DB.Delete(ctxspan, id)
	if err != nil {
		log.Error(err.Error())
		return false, liberrors.WrapErr(op, err)
	}

	err = s.Cashe.Invalidate(ctxspan, strconv.Itoa(id))
	if err != nil {
		log.Error(err.Error())
	}

	return result, nil

}

func (s *Service) List(ctx context.Context) (*[]entity.ElementOfList, error) {
	const op = "service.list"

	ctxspan, span := tracer.Start(ctx, "service_list")
	defer span.End()

	log := s.log.With(
		slog.String("op", op),
	)

	result, err := s.DB.List(ctxspan)
	if err != nil {
		log.Error(err.Error())
		return nil, liberrors.WrapErr(op, err)
	}

	return result, nil

}
