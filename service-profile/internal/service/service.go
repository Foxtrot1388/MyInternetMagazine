package service

import (
	"context"
	"errors"
	"golang.org/x/crypto/bcrypt"
	"log/slog"
	"v1/internal/entity"
	"v1/internal/lib"
)

type Service struct {
	DB  Repository
	log *slog.Logger
}

type Repository interface {
	Login(ctx context.Context, login string) (*entity.UserDB, error)
	Get(ctx context.Context, id int) (*entity.User, error)
	Create(ctx context.Context, user *entity.NewUser) (int, error)
	Delete(ctx context.Context, id int) (bool, error)
}

func New(log *slog.Logger, DB Repository) *Service {
	return &Service{DB: DB, log: log}
}

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
)

func (s *Service) Login(ctx context.Context, login, pass string) (*entity.LoginUser, error) {
	const op = "service.login"

	log := s.log.With(
		slog.String("op", op),
	)

	userdb, err := s.DB.Login(ctx, login)
	if err != nil {
		log.Error(err.Error())
		return nil, lib.WrapErr(op, ErrInvalidCredentials)
	}

	if err := bcrypt.CompareHashAndPassword(userdb.Pass, []byte(pass)); err != nil {
		log.Error(err.Error())
		return nil, lib.WrapErr(op, ErrInvalidCredentials)
	}

	return &entity.LoginUser{
		Id:         userdb.Id,
		Firstname:  userdb.Firstname,
		Secondname: userdb.Secondname,
		Lastname:   userdb.Lastname,
		Email:      userdb.Email,
	}, nil

}

func (s *Service) Get(ctx context.Context, id int) (*entity.User, error) {
	const op = "service.get"

	log := s.log.With(
		slog.String("op", op),
	)

	user, err := s.DB.Get(ctx, id)
	if err != nil {
		log.Error(err.Error())
		return nil, lib.WrapErr(op, err)
	}

	return user, nil

}

func (s *Service) Create(ctx context.Context, login, pass, fname, sname, lname, email string) (int, error) {
	const op = "service.create"

	log := s.log.With(
		slog.String("op", op),
	)

	passHash, err := bcrypt.GenerateFromPassword([]byte(pass), bcrypt.DefaultCost)
	if err != nil {
		log.Error(err.Error())
		return 0, lib.WrapErr(op, err)
	}

	user := &entity.NewUser{
		Pass:       passHash,
		Login:      login,
		Firstname:  fname,
		Secondname: sname,
		Lastname:   lname,
		Email:      email,
	}

	id, err := s.DB.Create(ctx, user)
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

	ok, err := s.DB.Delete(ctx, id)
	if err != nil {
		log.Error(err.Error())
		return false, lib.WrapErr(op, err)
	}

	return ok, nil

}
