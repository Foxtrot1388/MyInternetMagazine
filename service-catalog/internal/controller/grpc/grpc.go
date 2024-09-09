package grpcapi

import (
	"context"
	catalog "v1/internal/catalog/proto"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Server struct {
	catalog.UnimplementedCatalogApiServer
	s Service
}

func New(s Service) *Server {
	return &Server{s: s}
}

func (s *Server) Ping(ctx context.Context, req *catalog.PingParams) (*catalog.PingResponse, error) {
	return &catalog.PingResponse{Ok: true}, nil
}

func (s *Server) List(ctx context.Context, req *catalog.ListParams) (*catalog.ListResponse, error) {

	result, err := s.s.List(ctx)
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to get list of product")
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

	if req.GetId() == 0 {
		return nil, status.Error(codes.InvalidArgument, "id is empty")
	}

	product, err := s.s.Get(ctx, int(req.GetId()))
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to get a product")
	}

	return &catalog.GetResponse{
		Name:        product.Name,
		Description: product.Description,
	}, nil

}

func (s *Server) Create(ctx context.Context, req *catalog.CreateRequest) (*catalog.CreateResponse, error) {

	if req.GetName() == "" {
		return nil, status.Error(codes.InvalidArgument, "Name is empty")
	}

	if req.GetDescription() == "" {
		return nil, status.Error(codes.InvalidArgument, "GetDescription is empty")
	}

	id, err := s.s.Create(ctx, req.GetName(), req.GetDescription())
	if err != nil {
		return &catalog.CreateResponse{}, status.Error(codes.Internal, "failed to create a product")
	} else {
		return &catalog.CreateResponse{
			Id: int32(id),
		}, nil
	}

}

func (s *Server) Delete(ctx context.Context, req *catalog.GetRequest) (*catalog.DeleteResponse, error) {

	if req.GetId() == 0 {
		return nil, status.Error(codes.InvalidArgument, "Id is empty")
	}

	result, err := s.s.Delete(ctx, int(req.GetId()))
	if err != nil {
		return &catalog.DeleteResponse{}, status.Error(codes.Internal, "failed to delete a product")
	} else {
		return &catalog.DeleteResponse{
			Ok: result,
		}, nil
	}

}
