package service

import (
	"context"
	"encoding/json"
	"errors"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"log/slog"
	"time"
	"v1/internal/entity"
	"v1/internal/lib"
)

type Service struct {
	DB         Repository
	Bus        KafkaSender
	log        *slog.Logger
	signingkey string
}

type Repository interface {
	Login(ctx context.Context, login string) (*entity.UserDB, error)
	Get(ctx context.Context, id int) (*entity.User, error)
	Create(ctx context.Context, user *entity.NewUser) (int, error)
	Delete(ctx context.Context, id int) (bool, error)
}

type KafkaSender interface {
	Send(v []byte, topic string) error
}

func New(log *slog.Logger, DB Repository, signingkey string, Bus KafkaSender) *Service {
	return &Service{DB: DB, log: log, signingkey: signingkey, Bus: Bus}
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

	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["id"] = userdb.Id
	claims["email"] = userdb.Email
	claims["login"] = userdb.Login
	claims["exp"] = time.Now().Add(2 * time.Hour).Unix()
	tokenString, err := token.SignedString([]byte(s.signingkey))
	if err != nil {
		log.Error(err.Error())
		return nil, lib.WrapErr(op, err)
	}

	return &entity.LoginUser{
		Token: tokenString,
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

	err := validatePass(pass)
	if err != nil {
		return 0, err
	}

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

	if user.Email != "" && s.Bus != nil {
		go s.sendRegistrationEmail(user)
	}

	return id, nil

}

func (s *Service) sendRegistrationEmail(user *entity.NewUser) {
	const op = "service.sendRegistrationEmail"

	log := s.log.With(
		slog.String("op", op),
	)

	kmes := struct {
		MessageType string            `json:"messagetype"`
		Data        map[string]string `json:"data"`
	}{
		MessageType: "Registartion",
		Data:        map[string]string{"email": user.Email, "Login": user.Login},
	}

	value, err := json.Marshal(kmes)
	if err != nil {
		log.Error(err.Error())
		return
	}

	err = s.Bus.Send([]byte(value), "email")
	if err != nil {
		log.Error(err.Error())
		return
	}

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

func validatePass(pass string) error {

	return validation.Validate(pass,
		validation.Required,
		validation.Length(8, 20),
		is.UTFLetterNumeric.Error("Разрешенны только символы и цифры"),
	)

}
