package api

import (
	"context"
	"gorm.io/gorm"
	"v1/pkg/profile/proto"
)

type Server struct {
	profile.UnimplementedProfileApiServer
	DB *gorm.DB
}

func (s *Server) Ping(ctx context.Context, req *profile.PingParams) (*profile.PingResponse, error) {
	return &profile.PingResponse{Ok: true}, nil
}

func (s *Server) Login(ctx context.Context, req *profile.LoginRequest) (*profile.LoginResponse, error) {

	var user LoginUserRequest
	err := s.DB.WithContext(ctx).Where("pass = ? AND login = ?", req.Pass, req.Login).First(&user).Error
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

	var user GetUserRequest
	err := s.DB.WithContext(ctx).Take(&user, req.Id).Error
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

	user := CreateUserRequest{
		Pass:       req.Pass,
		Login:      req.Login,
		Firstname:  req.Firstname,
		Secondname: req.Secondname,
		Lastname:   req.Lastname,
		Email:      req.Email,
	}

	err := s.DB.WithContext(ctx).Create(&user).Error

	if err != nil {
		return &profile.CreateResponse{}, err
	} else {
		return &profile.CreateResponse{
			Id: int32(user.Id),
		}, nil
	}

}

func (s *Server) Delete(ctx context.Context, req *profile.GetRequest) (*profile.DeleteResponse, error) {

	err := s.DB.WithContext(ctx).Delete(&DeleteUserRequest{Id: int(req.Id)}).Error
	if err != nil {
		return nil, err
	}

	if err != nil {
		return &profile.DeleteResponse{}, err
	} else {
		return &profile.DeleteResponse{
			Ok: true,
		}, nil
	}

}
