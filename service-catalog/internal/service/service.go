package service

import (
	"context"
	"log/slog"
	"strconv"
	"v1/internal/entity"
	liberrors "v1/internal/lib/errors"
	"v1/internal/model"
	"v1/internal/model/converter"

	"go.opentelemetry.io/otel/trace"
)

type Service struct {
	DB     DBRepository
	Cashe  CasheRepository
	log    *slog.Logger
	tracer trace.Tracer
}

func NewService(log *slog.Logger, DB DBRepository, cashe CasheRepository, tracer trace.Tracer) *Service {
	return &Service{DB: DB, log: log, Cashe: cashe, tracer: tracer}
}

func (s *Service) Get(ctx context.Context, id int) (*model.Product, error) {
	const op = "service.get"

	ctxspan, span := s.tracer.Start(ctx, "service_get")
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

	return converter.GetProduct(product), nil

}

func (s *Service) Create(ctx context.Context, name string, description string) (int, error) {
	const op = "service.create"

	ctxspan, span := s.tracer.Start(ctx, "service_create")
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

	ctxspan, span := s.tracer.Start(ctx, "service_delete")
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

func (s *Service) List(ctx context.Context) (*[]model.ElementOfList, error) {
	const op = "service.list"

	ctxspan, span := s.tracer.Start(ctx, "service_list")
	defer span.End()

	log := s.log.With(
		slog.String("op", op),
	)

	resultentity, err := s.DB.List(ctxspan)
	if err != nil {
		log.Error(err.Error())
		return nil, liberrors.WrapErr(op, err)
	}

	result := make([]model.ElementOfList, len(*resultentity))
	for c, element := range *resultentity {
		result[c] = converter.GetElementOfList(element)
	}

	return &result, nil

}
