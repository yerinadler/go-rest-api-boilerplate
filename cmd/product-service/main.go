package main

import (
	"context"
	"log"

	"github.com/example/go-rest-api-boilerplate/internal/product-service/adapters"
	"github.com/example/go-rest-api-boilerplate/internal/product-service/api"
	"github.com/example/go-rest-api-boilerplate/internal/product-service/config"
	gormDb "github.com/example/go-rest-api-boilerplate/internal/product-service/db/gorm"
	"github.com/example/go-rest-api-boilerplate/internal/product-service/db/gorm/models"
	"github.com/example/go-rest-api-boilerplate/internal/product-service/logger"
	"github.com/example/go-rest-api-boilerplate/internal/product-service/services"
	"github.com/example/go-rest-api-boilerplate/pkg/messaging/kafka"
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

	ctx := context.Background()

	shutdownFunctions := observability.InitialiseOpentelemetry(ctx, cfg.Otlp.Endpoint, cfg.Application.Name)
	defer func() {
		for _, shutdownFunction := range shutdownFunctions {
			if err := shutdownFunction(ctx); err != nil {
				log.Fatal(err)
			}
		}
	}()

	tracer := otel.Tracer("main")

	db, err := gormDb.GetGormClient(cfg.Gorm.Dsn)
	if err != nil {
		log.Fatal(err)
	}

	db.AutoMigrate(&models.Product{})

	e := echo.New()

	e.Use(middlewares.OtelMiddleware(cfg.Application.Name))
	customErrorHandler := middlewares.NewCustomerErrorHandler(tracer, logger)
	e.HTTPErrorHandler = customErrorHandler.Handle
	e.Use(middleware.Recover())

	kafkaProducer, err := kafka.NewKafkaProducer(cfg.Kafka.Brokers, logger, tracer)
	if err != nil {
		log.Fatal(err)
	}

	customerAdapter := adapters.NewCustomerAdapter(cfg.External.Services.Customer)
	systemService := services.NewSystemService(tracer)
	productService := services.NewProductService(tracer, logger, db, kafkaProducer, customerAdapter)

	api.InitSystemHandler(e, systemService)
	api.InitProductHandler(e, productService)

	e.Logger.Fatal(e.Start(":" + cfg.Port))
}
