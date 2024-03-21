package main

import (
	"log"

	"github.com/example/go-rest-api-revision/config"
	"github.com/example/go-rest-api-revision/internal/api"
	gormDb "github.com/example/go-rest-api-revision/internal/db/gorm"
	"github.com/example/go-rest-api-revision/internal/db/gorm/models"
	"github.com/example/go-rest-api-revision/internal/logger"
	"github.com/example/go-rest-api-revision/internal/services"
	"github.com/example/go-rest-api-revision/pkg/middlewares"
	"github.com/example/go-rest-api-revision/pkg/observability"
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

	shutdown := observability.InitTracer(cfg.Otlp.Endpoint, cfg.ApplicationName)
	defer shutdown()

	tracer := otel.Tracer("main")

	db, err := gormDb.GetGormClient(cfg.GormDsn)
	if err != nil {
		log.Fatal(err)
	}

	db.AutoMigrate(&models.Product{})

	e := echo.New()

	e.Use(middlewares.OtelMiddleware(cfg.ApplicationName))
	customErrorHandler := middlewares.NewCustomerErrorHandler(tracer, logger)
	e.HTTPErrorHandler = customErrorHandler.Handle
	e.Use(middleware.Recover())

	systemService := services.NewSystemService(tracer)
	productService := services.NewProductService(tracer, db)

	api.InitSystemHandler(e, systemService)
	api.InitProductHandler(e, productService)

	e.Logger.Fatal(e.Start(":" + cfg.Port))
}
