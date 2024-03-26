package middlewares

import (
	"errors"
	"fmt"
	"net/http"

	exception "github.com/example/go-rest-api-revision/pkg/exceptions"
	"github.com/example/go-rest-api-revision/pkg/responses"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

type CustomErrorHandler struct {
	tracer trace.Tracer
	logger *logrus.Logger
}

func NewCustomerErrorHandler(tracer trace.Tracer, logger *logrus.Logger) *CustomErrorHandler {
	return &CustomErrorHandler{
		tracer,
		logger,
	}
}

func (eh *CustomErrorHandler) Handle(err error, c echo.Context) {
	var appErr *exception.ApplicationError

	ctx := c.Request().Context()
	span := trace.SpanFromContext(ctx)
	defer span.End()

	if errors.As(err, &appErr) {
		eh.logger.WithContext(ctx).WithError(err).WithField("code", appErr.Code).Error(err.Error())
		if err := c.JSON(exception.GetHttpStatusForCode(appErr.Code), responses.Response{
			Code:    appErr.Code,
			Message: appErr.Message,
		}); err != nil {
			fmt.Println(err)
		}

		span.SetStatus(codes.Error, appErr.Message)
		span.SetAttributes(attribute.Int("http.status_code", exception.GetHttpStatusForCode(appErr.Code)))
	} else {
		eh.logger.WithContext(ctx).WithError(err).Error(err.Error())
		if err := c.JSON(http.StatusInternalServerError, responses.Response{
			Code:    exception.UNEXPECTED_EXCEPTION,
			Message: err.Error(),
		}); err != nil {
			fmt.Println(err)
		}
		span.SetStatus(codes.Error, err.Error())
		span.SetAttributes(attribute.Int("http.status_code", http.StatusInternalServerError))
	}
}
