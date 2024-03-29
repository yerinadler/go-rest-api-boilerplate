package api

import (
	"net/http"

	"github.com/example/go-rest-api-revision/packages/product-service/dtos"
	"github.com/example/go-rest-api-revision/packages/product-service/services"
	"github.com/example/go-rest-api-revision/pkg/responses"
	"github.com/labstack/echo/v4"
)

type ProductHandler struct {
	productService *services.ProductService
}

func InitProductHandler(e *echo.Echo, productService *services.ProductService) {
	handler := &ProductHandler{
		productService,
	}
	e.Group("/")
	e.POST("/products", handler.CreateProduct)
	e.GET("/products", handler.GetProducts)
	e.GET("/products/:id", handler.GetProductById)
}

func (h *ProductHandler) CreateProduct(c echo.Context) error {

	var dto dtos.ProductDto

	if err := c.Bind(&dto); err != nil {
		return err
	}

	if err := h.productService.CreateProduct(c.Request().Context(), &dto); err != nil {
		return err
	}

	return c.JSON(http.StatusCreated, responses.Ok("Success", map[string]string{
		"message": "successfully created the new product",
	}))
}

func (h *ProductHandler) GetProducts(c echo.Context) error {
	productDtos, err := h.productService.GetAllProducts(c.Request().Context())
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, responses.Ok("successfully retrieved the product list", productDtos))
}

func (h *ProductHandler) GetProductById(c echo.Context) error {
	productDto, err := h.productService.GetProductById(c.Request().Context(), c.Param("id"))
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, responses.Ok("successfully retrieved the product", productDto))
}
