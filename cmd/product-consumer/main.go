package main

import (
	"context"
	"encoding/json"
	"log"

	"github.com/IBM/sarama"
	"github.com/dnwe/otelsarama"
	"github.com/example/go-rest-api-boilerplate/internal/product-consumer/events"
	"github.com/example/go-rest-api-boilerplate/internal/product-consumer/events/handlers"
	"github.com/example/go-rest-api-boilerplate/internal/product-service/config"
	"github.com/example/go-rest-api-boilerplate/internal/product-service/logger"
	"github.com/example/go-rest-api-boilerplate/pkg/messaging/kafka"
	"github.com/example/go-rest-api-boilerplate/pkg/observability"
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

	eventHandler := handlers.NewProductEventHandler(logger, tracer)

	if err := kafka.StartConsumption(
		cfg.Kafka.Brokers,
		[]string{"test"},
		"go-rest-api-product-consumer",
		"go-rest-api",
		func(message *sarama.ConsumerMessage) error {
			log.Printf("Message claimed: value = %s, timestamp = %v, topic = %s", string(message.Value), message.Timestamp, message.Topic)

			ctx := otel.GetTextMapPropagator().Extract(context.Background(), otelsarama.NewConsumerMessageCarrier(message))

			var productCreatedEvent events.ProductCreated
			if err := json.Unmarshal(message.Value, &productCreatedEvent); err != nil {
				log.Fatal(err)
			}

			if err := eventHandler.Handle(ctx, &productCreatedEvent); err != nil {
				log.Fatal(err)
			}

			return nil
		},
	); err != nil {
		log.Fatal(err)
	}
}
