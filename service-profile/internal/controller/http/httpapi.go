package httpapi

import (
	"context"
	"net/http"
	"strconv"
	"sync/atomic"
	_ "v1/api"
	"v1/internal/entity"

	"github.com/gin-gonic/gin"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	otelcodes "go.opentelemetry.io/otel/codes"
	oteltrace "go.opentelemetry.io/otel/trace"
)

var (
	healthy int32
	ready   int32
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

var tracer = otel.Tracer("profile-server")

func New(s Service) (*Server, *int32, *int32) {

	server := Server{
		s: s,
		R: gin.Default(),
	}

	server.R.Use(otelgin.Middleware("my-server"))

	server.R.POST("/profile/new", server.create)
	server.R.POST("/profile/login", server.login)
	server.R.DELETE("/profile/:id", server.delete)
	server.R.GET("/profile/:id", server.get)
	server.R.GET("/healthz/", server.healthz)
	server.R.GET("/readyz/", server.readyz)
	server.R.POST("/readyz/enable/", server.enableReady)
	server.R.POST("/readyz/disable/", server.disableReady)

	server.R.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	return &server, &healthy, &ready
}

// Healthz godoc
// @Summary Liveness check
// @Description used by Kubernetes liveness probe
// @Tags Kubernetes
// @Accept json
// @Produce json
// @Router /healthz [get]
// @Success 200 {string} string "OK"
func (s *Server) healthz(c *gin.Context) {
	if atomic.LoadInt32(&healthy) == 1 {
		c.JSON(http.StatusOK, map[string]string{"status": "OK"})
		return
	}
	c.Status(http.StatusServiceUnavailable)
}

// Readyz godoc
// @Summary Readiness check
// @Description used by Kubernetes readiness probe
// @Tags Kubernetes
// @Accept json
// @Produce json
// @Router /readyz [get]
// @Success 200 {string} string "OK"
func (s *Server) readyz(c *gin.Context) {
	if atomic.LoadInt32(&ready) == 1 {
		c.JSON(http.StatusOK, map[string]string{"status": "OK"})
		return
	}
	c.Status(http.StatusServiceUnavailable)
}

// EnableReady godoc
// @Summary Enable ready state
// @Description signals the Kubernetes LB that this instance is ready to receive traffic
// @Tags Kubernetes
// @Accept json
// @Produce json
// @Router /readyz/enable [post]
// @Success 202 {string} string "OK"
func (s *Server) enableReady(c *gin.Context) {
	atomic.StoreInt32(&ready, 1)
	c.Status(http.StatusAccepted)
}

// DisableReady godoc
// @Summary Disable ready state
// @Description signals the Kubernetes LB to stop sending requests to this instance
// @Tags Kubernetes
// @Accept json
// @Produce json
// @Router /readyz/disable [post]
// @Success 202 {string} string "OK"
func (s *Server) disableReady(c *gin.Context) {
	atomic.StoreInt32(&ready, 0)
	c.Status(http.StatusAccepted)
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

	commonAttrs := []attribute.KeyValue{
		attribute.String("Login", newUserRequest.Login),
		attribute.String("Pass", newUserRequest.Pass),
		attribute.String("Lastname", newUserRequest.Lastname),
		attribute.String("Secondname", newUserRequest.Secondname),
		attribute.String("Firstname", newUserRequest.Firstname),
		attribute.String("Email", newUserRequest.Email),
	}

	ctx, span := tracer.Start(c, "http_create", oteltrace.WithAttributes(commonAttrs...))
	defer span.End()

	err := newUserRequest.Validate()
	if err != nil {
		c.JSON(http.StatusInternalServerError, response{
			Message: err.Error(),
		})
		span.SetStatus(otelcodes.Error, err.Error())
		return
	}

	span.AddEvent("validation ok")

	id, err := s.s.Create(ctx, newUserRequest.Login, newUserRequest.Pass, newUserRequest.Firstname,
		newUserRequest.Secondname, newUserRequest.Lastname, newUserRequest.Email)
	if err != nil {
		_, ok := err.(validation.Error)
		if ok {
			c.JSON(http.StatusInternalServerError, response{
				Message: err.Error(),
			})
			span.SetStatus(otelcodes.Error, err.Error())
			return
		} else {
			c.JSON(http.StatusInternalServerError, response{
				Message: "failed to create profile",
			})
			span.SetStatus(otelcodes.Error, err.Error())
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

	ctx, span := tracer.Start(c, "http_delete", oteltrace.WithAttributes(attribute.String("id", idparam)))
	defer span.End()

	id, err := strconv.Atoi(idparam)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response{
			Message: "failed to delete profile",
		})
		span.SetStatus(otelcodes.Error, err.Error())
		return
	}

	result, err := s.s.Delete(ctx, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response{
			Message: "failed to delete profile",
		})
		span.SetStatus(otelcodes.Error, err.Error())
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

	ctx, span := tracer.Start(c.Request.Context(), "http_get", oteltrace.WithAttributes(attribute.String("id", idparam)))
	defer span.End()

	id, err := strconv.Atoi(idparam)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response{
			Message: "failed to get profile",
		})
		span.SetStatus(otelcodes.Error, err.Error())
		return
	}

	user, err := s.s.Get(ctx, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response{
			Message: "failed to get profile",
		})
		span.SetStatus(otelcodes.Error, err.Error())
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

	commonAttrs := []attribute.KeyValue{
		attribute.String("Login", UserRequest.Login),
		attribute.String("Pass", UserRequest.Pass),
	}

	ctx, span := tracer.Start(c.Request.Context(), "http_login", oteltrace.WithAttributes(commonAttrs...))
	defer span.End()

	err := UserRequest.Validate()
	if err != nil {
		c.JSON(http.StatusInternalServerError, response{
			Message: err.Error(),
		})
		span.SetStatus(otelcodes.Error, err.Error())
		return
	}

	span.AddEvent("validate ok")

	user, err := s.s.Login(ctx, UserRequest.Login, UserRequest.Pass)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response{
			Message: "failed to login",
		})
		span.SetStatus(otelcodes.Error, err.Error())
		return
	}

	c.JSON(http.StatusOK, tokenResponse{
		Token: user.Token,
	})

}
