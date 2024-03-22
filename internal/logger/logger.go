package logger

import (
	"github.com/sirupsen/logrus"
	"go.elastic.co/ecslogrus"
	"go.opentelemetry.io/otel/trace"
)

type CustomFormatter struct {
	ecslogrus.Formatter
}

func (f *CustomFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	ctx := entry.Context
	span := trace.SpanFromContext(ctx)

	entry.Data["trace.id"] = span.SpanContext().TraceID().String()
	entry.Data["span.id"] = span.SpanContext().SpanID().String()

	return f.Formatter.Format(entry)
}

func GetLogger() *logrus.Logger {
	logger := logrus.New()
	logger.SetFormatter(&CustomFormatter{
		Formatter: ecslogrus.Formatter{},
	})
	logger.SetLevel(logrus.InfoLevel)

	return logger
}
