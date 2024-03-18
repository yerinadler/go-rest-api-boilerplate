package main

import (
	"log"

	"github.com/example/go-rest-api-revision/config"
	"github.com/example/go-rest-api-revision/internal/api"
	gormDb "github.com/example/go-rest-api-revision/internal/db/gorm"
	"github.com/example/go-rest-api-revision/internal/db/gorm/models"
	"github.com/example/go-rest-api-revision/internal/services"
	"github.com/example/go-rest-api-revision/pkg/observability"
	"github.com/labstack/echo/v4"
	"go.opentelemetry.io/contrib/instrumentation/github.com/labstack/echo/otelecho"
	"go.opentelemetry.io/otel"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal(err)
	}

	shutdown := observability.InitTracer(cfg.Otlp.Endpoint)
	defer shutdown()

	tracer := otel.Tracer("main")

	db, err := gormDb.GetGormClient(cfg.GormDsn)
	if err != nil {
		log.Fatal(err)
	}

	db.AutoMigrate(&models.Product{})

	echo := echo.New()
	echo.Use(otelecho.Middleware("go-rest-api"))

	systemService := services.NewSystemService(tracer)
	productService := services.NewProductService(tracer, db)

	api.InitSystemHandler(echo, systemService)
	api.InitProductHandler(echo, productService)

	echo.Logger.Fatal(echo.Start(":" + cfg.Port))
}
