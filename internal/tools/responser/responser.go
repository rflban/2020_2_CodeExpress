package responser

import (
	. "github.com/go-park-mail-ru/2020_2_CodeExpress/internal/tools/error_response"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
)

func RespondWithError(err *ErrorResponse, ctx echo.Context) error {
	if err.Error != nil {
		logrus.Info(err.Error)
	}
	return ctx.JSON(err.StatusCode, err.UserError)
}
