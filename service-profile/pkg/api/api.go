package api

import (
	"context"
	entity "v1/internal"
	"v1/internal/profile/proto"
)

type Server struct {
	profile.UnimplementedProfileApiServer
	DB Repository
}

type Repository interface {
	Login(ctx context.Context, pass, login string) (*entity.LoginUser, error)
	Get(ctx context.Context, id int) (*entity.User, error)
	Create(ctx context.Context, user *entity.NewUser) (int, error)
	Delete(ctx context.Context, id int) (bool, error)
}

func (s *Server) Ping(ctx context.Context, req *profile.PingParams) (*profile.PingResponse, error) {
	return &profile.PingResponse{Ok: true}, nil
}

func (s *Server) Login(ctx context.Context, req *profile.LoginRequest) (*profile.LoginResponse, error) {

	user, err := s.DB.Login(ctx, req.Pass, req.Login)
	if err != nil {
		return nil, err
	}

	return &profile.LoginResponse{
		Id:         int32(user.Id),
		Firstname:  user.Firstname,
		Secondname: user.Secondname,
		Lastname:   user.Lastname,
		Email:      user.Email,
	}, nil

}

func (s *Server) Get(ctx context.Context, req *profile.GetRequest) (*profile.GetResponse, error) {

	user, err := s.DB.Get(ctx, int(req.Id))
	if err != nil {
		return nil, err
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

	user := entity.NewUser{
		Pass:       req.Pass,
		Login:      req.Login,
		Firstname:  req.Firstname,
		Secondname: req.Secondname,
		Lastname:   req.Lastname,
		Email:      req.Email,
	}

	id, err := s.DB.Create(ctx, &user)

	if err != nil {
		return &profile.CreateResponse{}, err
	} else {
		return &profile.CreateResponse{
			Id: int32(id),
		}, nil
	}

}

func (s *Server) Delete(ctx context.Context, req *profile.GetRequest) (*profile.DeleteResponse, error) {

	result, err := s.DB.Delete(ctx, int(req.Id))

	if err != nil {
		return &profile.DeleteResponse{}, err
	} else {
		return &profile.DeleteResponse{
			Ok: result,
		}, nil
	}

}
