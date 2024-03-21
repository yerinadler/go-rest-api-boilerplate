package middlewares

import (
	"errors"
	"fmt"

	exception "github.com/example/go-rest-api-revision/pkg/exceptions"
	"github.com/example/go-rest-api-revision/pkg/responses"
	"github.com/labstack/echo/v4"
)

func ErrorHandler(err error, c echo.Context) {
	var appErr *exception.ApplicationError
	if errors.As(err, &appErr) {
		if err := c.JSON(exception.GetHttpStatusForCode(appErr.Code), responses.Response{
			Code:    appErr.Code,
			Message: appErr.Message,
		}); err != nil {
			fmt.Println(err)
		}
	} else {
		fmt.Println(err)
	}
}
