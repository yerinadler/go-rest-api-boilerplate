package api

import (
	"net/http"

	"github.com/example/go-rest-api-revision/packages/product-service/services"
	"github.com/example/go-rest-api-revision/pkg/responses"
	"github.com/labstack/echo/v4"
)

type SystemHandler struct {
	systemService *services.SystemService
}

func InitSystemHandler(e *echo.Echo, systemService *services.SystemService) {
	handler := &SystemHandler{
		systemService,
	}
	e.Group("/")
	e.GET("/healthz", handler.HealthCheck)
	e.GET("/hello", handler.GetHello)
}

func (h *SystemHandler) HealthCheck(c echo.Context) error {
	return c.JSON(http.StatusOK, responses.Ok("Success", nil))
}

func (h *SystemHandler) GetHello(c echo.Context) error {
	message := h.systemService.GetHelloMessage(c.Request().Context())

	return c.JSON(http.StatusOK, responses.Ok("Success", map[string]string{
		"hello_message": message,
	}))
}
