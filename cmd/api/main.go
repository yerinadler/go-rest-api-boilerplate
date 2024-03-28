package main

import (
	"log"

	"github.com/example/go-rest-api-revision/config"
	"github.com/example/go-rest-api-revision/internal/api"
	gormDb "github.com/example/go-rest-api-revision/internal/db/gorm"
	"github.com/example/go-rest-api-revision/internal/db/gorm/models"
	"github.com/example/go-rest-api-revision/internal/logger"
	"github.com/example/go-rest-api-revision/internal/services"
	"github.com/example/go-rest-api-revision/pkg/messaging/kafka"
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

	shutdown := observability.InitialiseOpentelemetry(cfg.Otlp.Endpoint, cfg.Application.Name)
	defer shutdown()

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
	systemService := services.NewSystemService(tracer)
	productService := services.NewProductService(tracer, db, kafkaProducer)

	api.InitSystemHandler(e, systemService)
	api.InitProductHandler(e, productService)

	e.Logger.Fatal(e.Start(":" + cfg.Port))
}
