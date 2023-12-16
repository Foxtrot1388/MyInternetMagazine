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
	Get(ctx context.Context, id int) (*entity.Product, error)
	Create(ctx context.Context, name string, description string) (int, error)
	Delete(ctx context.Context, id int) (bool, error)
	List(ctx context.Context) (*[]entity.ElementOfList, error)
}

type responseElement struct {
	Id   int
	Name string
}

type listResponse []responseElement

type response struct {
	Message string
}

type getResponse struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

type deleteResponse struct {
	OK bool
}

type createResponse struct {
	Id int
}

type newCatalog struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

func (a newCatalog) Validate() error {
	return validation.ValidateStruct(&a,
		validation.Field(&a.Name, validation.Required),
		validation.Field(&a.Description, validation.Required),
	)
}

func New(s Service) *Server {

	server := Server{
		s: s,
		R: gin.Default(),
	}

	server.R.GET("/ping", server.ping)
	server.R.POST("/catalog/new", server.create)
	server.R.DELETE("/catalog/:id", server.delete)
	server.R.GET("/catalog/:id", server.get)
	server.R.GET("/catalog/", server.list)

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

// @Summary List
// @Tags catalog
// @Description create
// @ID list
// @Produce json
// @Success 200 {object} listResponse
// @Failure 500 {object} response
// @Router /catalog/ [get]
func (s *Server) list(c *gin.Context) {

	result, err := s.s.List(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response{
			Message: "failed to get list of product",
		})
		return
	}

	productslist := make(listResponse, len(*result))
	for k, v := range *result {
		productslist[k] = responseElement{
			Id:   v.Id,
			Name: v.Name,
		}
	}

	c.JSON(http.StatusOK, productslist)

}

// @Summary Get
// @Tags catalog
// @Description get
// @ID get
// @Produce json
// @Param catalog_id path int true "Product ID"
// @Success 200 {object} getResponse
// @Failure 500 {object} response
// @Router /catalog/{catalog_id} [get]
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
			Message: "failed to get product",
		})
		return
	}

	product, err := s.s.Get(c, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response{
			Message: "failed to get a product",
		})
		return
	}

	c.JSON(http.StatusOK, getResponse{
		Name:        product.Name,
		Description: product.Description,
	})

}

// @Summary Create
// @Tags catalog
// @Description create
// @ID create
// @Accept json
// @Produce json
// @Param input body newCatalog true "new"
// @Success 200 {object} createResponse
// @Failure 500 {object} response
// @Router /catalog/new [post]
func (s *Server) create(c *gin.Context) {

	var newCatalogRequest newCatalog
	if err := c.BindJSON(&newCatalogRequest); err != nil {
		c.JSON(http.StatusInternalServerError, response{
			Message: "bad request",
		})
		return
	}

	err := newCatalogRequest.Validate()
	if err != nil {
		c.JSON(http.StatusInternalServerError, response{
			Message: err.Error(),
		})
		return
	}

	id, err := s.s.Create(c, newCatalogRequest.Name, newCatalogRequest.Description)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response{
			Message: "failed to create a product",
		})
		return
	} else {
		c.JSON(http.StatusOK, createResponse{
			Id: id,
		})
	}

}

// @Summary Delete
// @Tags catalog
// @Description delete
// @ID delete
// @Produce json
// @Param catalog_id path int true "Product ID"
// @Success 200 {object} deleteResponse
// @Failure 500 {object} response
// @Router /catalog/{catalog_id} [delete]
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
			Message: "failed to delete catalog",
		})
		return
	}

	result, err := s.s.Delete(c, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response{
			Message: "failed to delete a catalog",
		})
		return
	} else {
		c.JSON(http.StatusOK, deleteResponse{
			OK: result,
		})
	}

}
