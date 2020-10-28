package validator

import (
	"fmt"

	. "github.com/go-park-mail-ru/2020_2_CodeExpress/internal/consts"
	. "github.com/go-park-mail-ru/2020_2_CodeExpress/internal/tools/error_response"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

type ValidationErrors = validator.ValidationErrors

type RequestValidator struct {
	ctx       echo.Context
	validator *validator.Validate
}

func NewRequestValidator(ctx echo.Context) *RequestValidator {
	return &RequestValidator{
		ctx:       ctx,
		validator: validator.New(),
	}
}

func (rv *RequestValidator) Validate(request interface{}) *ErrorResponse {
	if err := rv.ctx.Bind(request); err != nil {
		return NewErrorResponse(ErrBadRequest, err)
	}

	if err := rv.validator.Struct(request); err != nil {
		return NewErrorResponse(ErrBadRequest, err)
	}

	return nil
}

func GetValidationError(errResp *ErrorResponse) {
	if _, ok := errResp.Error.(validator.ValidationErrors); !ok {
		return
	}

	for _, err := range errResp.Error.(validator.ValidationErrors) {
		switch err.Tag() {
		case "required":
			errResp.UserError = &ErrorJson{
				Message: fmt.Sprintf("Field %s is required", err.Field()),
			}
			return
		case "gte":
			errResp.UserError = &ErrorJson{
				Message: fmt.Sprintf("Field %s should be at least %s symbols", err.Field(), err.Param()),
			}
			return
		case "eqfield":
			errResp.UserError = &ErrorJson{
				Message: fmt.Sprintf("Field %s and %s do not match", err.Field(), err.Param()),
			}
			return
		case "email":
			errResp.UserError = &ErrorJson{
				Message: fmt.Sprintf("String %s is not an %s", err.Value(), err.Field()),
			}
			return
		}
	}
}
