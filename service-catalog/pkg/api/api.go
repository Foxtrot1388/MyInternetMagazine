package api

import (
	"context"
	"strconv"
	entity "v1/internal"
	"v1/internal/catalog/proto"
	"v1/internal/lib"
)

type Server struct {
	catalog.UnimplementedCatalogApiServer
	DB    DBRepository
	Cashe CasheRepository
}

type CasheRepository interface {
	Get(ctx context.Context, key string) (*entity.Product, error)
	Set(ctx context.Context, key string, v *entity.Product) error
}

type DBRepository interface {
	Get(ctx context.Context, id int) (*entity.Product, error)
	Create(ctx context.Context, product *entity.Product) (int, error)
	Delete(ctx context.Context, id int) (bool, error)
	List(ctx context.Context) (*[]entity.ElementOfList, error)
}

func (s *Server) Ping(ctx context.Context, req *catalog.PingParams) (*catalog.PingResponse, error) {
	return &catalog.PingResponse{Ok: true}, nil
}

func (s *Server) List(ctx context.Context, req *catalog.ListParams) (*catalog.ListResponse, error) {
	const op = "api.list"

	result, err := s.DB.List(ctx)
	if err != nil {
		return nil, lib.WrapErr(op, err)
	}

	productslist := make([]*catalog.ElementOfList, len(*result))
	for k, v := range *result {
		productslist[k] = &catalog.ElementOfList{
			Id:   int32(v.Id),
			Name: v.Name,
		}
	}

	return &catalog.ListResponse{List: productslist}, nil

}

func (s *Server) Get(ctx context.Context, req *catalog.GetRequest) (*catalog.GetResponse, error) {
	const op = "api.get"

	product, err := s.Cashe.Get(ctx, strconv.Itoa(int(req.Id)))
	if err != nil {

		product, err = s.DB.Get(ctx, int(req.Id))
		if err != nil {
			return nil, lib.WrapErr(op, err)
		}

		err = s.Cashe.Set(ctx, strconv.Itoa(int(req.Id)), product)
		if err != nil {
			return nil, lib.WrapErr(op, err)
		}

	}

	return &catalog.GetResponse{
		Name:        product.Name,
		Description: product.Description,
	}, nil

}

func (s *Server) Create(ctx context.Context, req *catalog.CreateRequest) (*catalog.CreateResponse, error) {
	const op = "api.create"

	product := entity.Product{
		Name:        req.Name,
		Description: req.Description,
	}

	id, err := s.DB.Create(ctx, &product)

	if err != nil {
		return &catalog.CreateResponse{}, lib.WrapErr(op, err)
	} else {
		return &catalog.CreateResponse{
			Id: int32(id),
		}, nil
	}

}

func (s *Server) Delete(ctx context.Context, req *catalog.GetRequest) (*catalog.DeleteResponse, error) {
	const op = "api.delete"

	result, err := s.DB.Delete(ctx, int(req.Id))

	if err != nil {
		return &catalog.DeleteResponse{}, lib.WrapErr(op, err)
	} else {
		return &catalog.DeleteResponse{
			Ok: result,
		}, nil
	}

}
