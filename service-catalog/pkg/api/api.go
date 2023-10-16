package api

import (
	"context"
	"strconv"
	entity "v1/internal"
	"v1/internal/catalog/proto"
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

	result, err := s.DB.List(ctx)
	if err != nil {
		return nil, err
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

	product, err := s.Cashe.Get(ctx, strconv.Itoa(int(req.Id)))
	if err != nil {
		
		product, err = s.DB.Get(ctx, int(req.Id))
		if err != nil {
			return nil, err
		}

		err = s.Cashe.Set(ctx, strconv.Itoa(int(req.Id)), product)
		if err != nil {
			return nil, err
		}

	}

	return &catalog.GetResponse{
		Name:        product.Name,
		Description: product.Description,
	}, nil

}

func (s *Server) Create(ctx context.Context, req *catalog.CreateRequest) (*catalog.CreateResponse, error) {

	product := entity.Product{
		Name:        req.Name,
		Description: req.Description,
	}

	id, err := s.DB.Create(ctx, &product)

	if err != nil {
		return &catalog.CreateResponse{}, err
	} else {
		return &catalog.CreateResponse{
			Id: int32(id),
		}, nil
	}

}

func (s *Server) Delete(ctx context.Context, req *catalog.GetRequest) (*catalog.DeleteResponse, error) {

	result, err := s.DB.Delete(ctx, int(req.Id))

	if err != nil {
		return &catalog.DeleteResponse{}, err
	} else {
		return &catalog.DeleteResponse{
			Ok: result,
		}, nil
	}

}
