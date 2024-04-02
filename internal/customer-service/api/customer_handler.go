package api

import (
	"net/http"

	"github.com/example/go-rest-api-boilerplate/internal/customer-service/dtos"
	"github.com/example/go-rest-api-boilerplate/internal/customer-service/services"
	"github.com/example/go-rest-api-boilerplate/pkg/responses"
	"github.com/labstack/echo/v4"
)

type CustomerHandler struct {
	customerService *services.CustomerService
}

func InitCustomerHandler(e *echo.Echo, customerService *services.CustomerService) {
	handler := &CustomerHandler{
		customerService,
	}
	e.Group("/")
	e.POST("/customers", handler.CreateCustomer)
	e.GET("/customers/:id", handler.GetCustomerById)
}

func (h *CustomerHandler) CreateCustomer(c echo.Context) error {

	var dto dtos.CustomerDto

	if err := c.Bind(&dto); err != nil {
		return err
	}

	if err := h.customerService.CreateCustomer(c.Request().Context(), &dto); err != nil {
		return err
	}

	return c.JSON(http.StatusCreated, responses.Ok("Success", map[string]string{
		"message": "successfully created the new product",
	}))
}

func (h *CustomerHandler) GetCustomerById(c echo.Context) error {
	productDto, err := h.customerService.GetCustomerById(c.Request().Context(), c.Param("id"))
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, responses.Ok("successfully retrieved the product", productDto))
}
