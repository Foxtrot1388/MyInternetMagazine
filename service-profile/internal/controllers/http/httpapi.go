package httpapi

import (
	"context"
	"github.com/gin-gonic/gin"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"net/http"
	"strconv"
	_ "v1/docs"
	"v1/internal/entity"
)

type Server struct {
	s Service
	R *gin.Engine
}

type Service interface {
	Login(ctx context.Context, pass, login string) (*entity.LoginUser, error)
	Get(ctx context.Context, id int) (*entity.User, error)
	Create(ctx context.Context, login, pass, fname, sname, lname, email string) (int, error)
	Delete(ctx context.Context, id int) (bool, error)
}

type newUser struct {
	Login      string `json:"login"`
	Pass       string `json:"pass"`
	Firstname  string `json:"firstname"`
	Secondname string `json:"secondname"`
	Lastname   string `json:"lastname"`
	Email      string `json:"email"`
}

func (a newUser) Validate() error {
	return validation.ValidateStruct(&a,
		validation.Field(&a.Email, validation.Required),
		validation.Field(&a.Login, validation.Required),
	)
}

type loginUser struct {
	Login string `json:"login"`
	Pass  string `json:"pass"`
}

func (a loginUser) Validate() error {
	return validation.ValidateStruct(&a,
		validation.Field(&a.Login, validation.Required),
		validation.Field(&a.Pass, validation.Required),
	)
}

type response struct {
	Message string
}

type createResponse struct {
	Id int
}

type deleteResponse struct {
	OK bool
}

type tokenResponse struct {
	Token string
}

type getResponse struct {
	Login      string `json:"login"`
	Firstname  string `json:"firstname"`
	Secondname string `json:"secondname"`
	Lastname   string `json:"lastname"`
	Email      string `json:"email"`
}

func New(s Service) *Server {

	server := Server{
		s: s,
		R: gin.Default(),
	}

	server.R.GET("/ping", server.ping)
	server.R.POST("/profile/new", server.create)
	server.R.POST("/profile/login", server.login)
	server.R.DELETE("/profile/:id", server.delete)
	server.R.GET("/profile/:id", server.get)

	server.R.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	return &server
}

// @Summary Ping
// @Tags other
// @Description ping
// @ID ping
// @Produce json
// @Success 200 {object} response
// @Router /ping [get]
func (s *Server) ping(c *gin.Context) {
	c.JSON(http.StatusOK, response{
		Message: "OK",
	})
}

// @Summary Create
// @Tags profile
// @Description create
// @ID create
// @Accept json
// @Produce json
// @Param input body newUser true "new"
// @Success 200 {object} createResponse
// @Failure 500 {object} response
// @Router /profile/new [post]
func (s *Server) create(c *gin.Context) {

	var newUserRequest newUser
	if err := c.BindJSON(&newUserRequest); err != nil {
		c.JSON(http.StatusInternalServerError, response{
			Message: "bad request",
		})
		return
	}

	err := newUserRequest.Validate()
	if err != nil {
		c.JSON(http.StatusInternalServerError, response{
			Message: err.Error(),
		})
		return
	}

	id, err := s.s.Create(c, newUserRequest.Login, newUserRequest.Pass, newUserRequest.Firstname,
		newUserRequest.Secondname, newUserRequest.Lastname, newUserRequest.Email)
	if err != nil {
		_, ok := err.(validation.Error)
		if ok {
			c.JSON(http.StatusInternalServerError, response{
				Message: err.Error(),
			})
			return
		} else {
			c.JSON(http.StatusInternalServerError, response{
				Message: "failed to create profile",
			})
			return
		}
	} else {
		c.JSON(http.StatusOK, createResponse{
			Id: id,
		})
	}

}

// @Summary Delete
// @Tags profile
// @Description delete
// @ID delete
// @Produce json
// @Param user_id path int true "User ID"
// @Success 200 {object} deleteResponse
// @Failure 500 {object} response
// @Router /profile/{user_id} [delete]
func (s *Server) delete(c *gin.Context) {

	idparam := c.Param("id")
	if idparam == "" {
		c.JSON(http.StatusInternalServerError, response{
			Message: "id is empty",
		})
		return
	}

	id, err := strconv.Atoi(idparam)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response{
			Message: "failed to delete profile",
		})
		return
	}

	result, err := s.s.Delete(c, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response{
			Message: "failed to delete profile",
		})
		return
	} else {
		c.JSON(http.StatusOK, deleteResponse{
			OK: result,
		})
	}

}

// @Summary Get
// @Tags profile
// @Description get
// @ID get
// @Produce json
// @Param user_id path int true "User ID"
// @Success 200 {object} getResponse
// @Failure 500 {object} response
// @Router /profile/{user_id} [get]
func (s *Server) get(c *gin.Context) {

	idparam := c.Param("id")
	if idparam == "" {
		c.JSON(http.StatusInternalServerError, response{
			Message: "id is empty",
		})
		return
	}

	id, err := strconv.Atoi(idparam)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response{
			Message: "failed to get profile",
		})
		return
	}

	user, err := s.s.Get(c, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response{
			Message: "failed to get profile",
		})
		return
	}

	c.JSON(http.StatusOK, getResponse{
		Login:      user.Login,
		Firstname:  user.Firstname,
		Secondname: user.Secondname,
		Lastname:   user.Lastname,
		Email:      user.Email,
	})

}

// @Summary Login
// @Tags auth
// @Description login
// @ID login
// @Accept json
// @Produce json
// @Param input body loginUser true "credentials"
// @Success 200 {object} tokenResponse
// @Failure 500 {object} response
// @Router /profile/login [post]
func (s *Server) login(c *gin.Context) {

	var UserRequest loginUser
	if err := c.BindJSON(&UserRequest); err != nil {
		c.JSON(http.StatusInternalServerError, response{
			Message: "bad request",
		})
		return
	}

	err := UserRequest.Validate()
	if err != nil {
		c.JSON(http.StatusInternalServerError, response{
			Message: err.Error(),
		})
		return
	}

	user, err := s.s.Login(c, UserRequest.Login, UserRequest.Pass)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response{
			Message: "failed to login",
		})
		return
	}

	c.JSON(http.StatusOK, tokenResponse{
		Token: user.Token,
	})

}
