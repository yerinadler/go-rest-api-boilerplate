package observability

import (
	"context"
	"time"

	exceptions "github.com/example/go-rest-api-revision/pkg/exceptions"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.24.0"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func newExporter(ctx context.Context, otlpEndpoint string) (*otlptrace.Exporter, error) {

	conn, err := grpc.DialContext(ctx, otlpEndpoint, grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithBlock())
	exceptions.ReportError(err, "unable to reach GRPC OTLP endpoint")

	return otlptracegrpc.New(ctx, otlptracegrpc.WithGRPCConn(conn))
}

func newResource(ctx context.Context, appName string) (*resource.Resource, error) {
	return resource.New(
		ctx,
		resource.WithAttributes(
			semconv.ServiceNameKey.String(appName),
			attribute.String("application", appName),
		),
	)
}

func newTraceProvider(resource *resource.Resource, spanProcessor sdktrace.SpanProcessor) *sdktrace.TracerProvider {
	tracerProvider := sdktrace.NewTracerProvider(
		sdktrace.WithSampler(sdktrace.AlwaysSample()),
		sdktrace.WithResource(resource),
		sdktrace.WithSpanProcessor(spanProcessor),
	)

	return tracerProvider
}

func InitTracer(otlpEndpoint string, appName string) func() {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)

	resource, err := newResource(ctx, appName)
	exceptions.ReportError(err, "failed to create the OTLP resource")

	exporter, err := newExporter(ctx, otlpEndpoint)
	exceptions.ReportError(err, "failed to created the OTLP exporter")

	batchSpanProcessor := sdktrace.NewBatchSpanProcessor(exporter)
	tracerProvider := newTraceProvider(resource, batchSpanProcessor)
	otel.SetTracerProvider(tracerProvider)
	otel.SetTextMapPropagator(propagation.TraceContext{})

	return func() {
		exceptions.ReportError(tracerProvider.Shutdown(ctx), "failed to gracefully shutdown the tracer provider")
		cancel()
	}
}
