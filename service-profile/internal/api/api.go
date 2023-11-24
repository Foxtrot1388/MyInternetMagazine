package api

import (
	"context"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"v1/internal/entity"
	"v1/internal/profile/proto"
)

type Server struct {
	profile.UnimplementedProfileApiServer
	S Service
}

type Service interface {
	Login(ctx context.Context, pass, login string) (*entity.LoginUser, error)
	Get(ctx context.Context, id int) (*entity.User, error)
	Create(ctx context.Context, login, pass, fname, sname, lname, email string) (int, error)
	Delete(ctx context.Context, id int) (bool, error)
}

func (s *Server) Ping(ctx context.Context, req *profile.PingParams) (*profile.PingResponse, error) {
	return &profile.PingResponse{Ok: true}, nil
}

func (s *Server) Login(ctx context.Context, req *profile.LoginRequest) (*profile.LoginResponse, error) {

	if req.GetLogin() == "" {
		return nil, status.Error(codes.InvalidArgument, "login is empty")
	}

	if req.GetPass() == "" {
		return nil, status.Error(codes.InvalidArgument, "pass is empty")
	}

	user, err := s.S.Login(ctx, req.GetLogin(), req.GetPass())
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to login")
	}

	return &profile.LoginResponse{
		Token: user.Token,
	}, nil

}

func (s *Server) Get(ctx context.Context, req *profile.GetRequest) (*profile.GetResponse, error) {

	if req.GetId() == 0 {
		return nil, status.Error(codes.InvalidArgument, "id is empty")
	}

	user, err := s.S.Get(ctx, int(req.GetId()))
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to get profile")
	}

	return &profile.GetResponse{
		Login:      user.Login,
		Firstname:  user.Firstname,
		Secondname: user.Secondname,
		Lastname:   user.Lastname,
		Email:      user.Email,
	}, nil

}

func (s *Server) Create(ctx context.Context, req *profile.CreateRequest) (*profile.CreateResponse, error) {

	if req.GetLogin() == "" {
		return nil, status.Error(codes.InvalidArgument, "login is empty")
	}

	if req.GetEmail() == "" {
		return nil, status.Error(codes.InvalidArgument, "email is empty")
	}

	id, err := s.S.Create(ctx, req.GetLogin(), req.GetPass(), req.GetFirstname(), req.GetSecondname(), req.GetLastname(), req.GetEmail())
	if err != nil {
		_, ok := err.(validation.Error)
		if ok {
			return &profile.CreateResponse{}, status.Error(codes.InvalidArgument, err.Error())
		} else {
			return &profile.CreateResponse{}, status.Error(codes.Internal, "failed to create profile")
		}
	} else {
		return &profile.CreateResponse{
			Id: int32(id),
		}, nil
	}

}

func (s *Server) Delete(ctx context.Context, req *profile.GetRequest) (*profile.DeleteResponse, error) {

	if req.GetId() == 0 {
		return nil, status.Error(codes.InvalidArgument, "id is empty")
	}

	result, err := s.S.Delete(ctx, int(req.GetId()))
	if err != nil {
		return &profile.DeleteResponse{}, status.Error(codes.Internal, "failed to delete profile")
	} else {
		return &profile.DeleteResponse{
			Ok: result,
		}, nil
	}

}
