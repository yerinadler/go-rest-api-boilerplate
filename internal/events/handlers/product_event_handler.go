package handlers

import (
	"context"

	"github.com/example/go-rest-api-revision/internal/events"
	"github.com/sirupsen/logrus"
	"go.opentelemetry.io/otel/trace"
)

type ProductEventHandler struct {
	logger *logrus.Logger
	tracer trace.Tracer
}

func NewProductEventHandler(logger *logrus.Logger, tracer trace.Tracer) *ProductEventHandler {
	return &ProductEventHandler{logger, tracer}
}

func (h *ProductEventHandler) Handle(ctx context.Context, event *events.ProductCreated) error {
	ctx, span := h.tracer.Start(ctx, "start processing ProductCreatedEvent")
	defer span.End()

	h.logger.WithContext(ctx).Infof("processed product with with the ID: %d", event.Id)

	return nil
}
