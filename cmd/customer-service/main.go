package main

import (
	"log"

	"github.com/example/go-rest-api-boilerplate/internal/customer-service/api"
	"github.com/example/go-rest-api-boilerplate/internal/customer-service/config"
	gormDb "github.com/example/go-rest-api-boilerplate/internal/customer-service/db/gorm"
	"github.com/example/go-rest-api-boilerplate/internal/customer-service/db/gorm/models"
	"github.com/example/go-rest-api-boilerplate/internal/customer-service/logger"
	"github.com/example/go-rest-api-boilerplate/internal/customer-service/services"
	"github.com/example/go-rest-api-boilerplate/pkg/middlewares"
	"github.com/example/go-rest-api-boilerplate/pkg/observability"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"go.opentelemetry.io/otel"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal(err)
	}

	logger := logger.GetLogger()

	shutdown := observability.InitialiseOpentelemetry(cfg.Otlp.Endpoint, cfg.Application.Name)
	defer shutdown()

	tracer := otel.Tracer("main")

	db, err := gormDb.GetGormClient(cfg.Gorm.Dsn)
	if err != nil {
		log.Fatal(err)
	}

	db.AutoMigrate(&models.Customer{})

	e := echo.New()

	e.Use(middlewares.OtelMiddleware(cfg.Application.Name))
	customErrorHandler := middlewares.NewCustomerErrorHandler(tracer, logger)
	e.HTTPErrorHandler = customErrorHandler.Handle
	e.Use(middleware.Recover())

	systemService := services.NewSystemService(tracer)
	customerService := services.NewCustomerService(tracer, logger, db)

	api.InitSystemHandler(e, systemService)
	api.InitCustomerHandler(e, customerService)

	e.Logger.Fatal(e.Start(":" + cfg.Port))
}
