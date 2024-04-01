package services

import (
	"context"

	"go.opentelemetry.io/otel/trace"
)

type SystemService struct {
	tracer trace.Tracer
}

func NewSystemService(tracer trace.Tracer) *SystemService {
	return &SystemService{
		tracer,
	}
}

func (s *SystemService) GetHelloMessage(ctx context.Context) string {
	_, span := s.tracer.Start(ctx, "generating hello message")
	defer span.End()
	return "Hello Go"
}
