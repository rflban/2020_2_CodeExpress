package error_response

import (
	. "github.com/go-park-mail-ru/2020_2_CodeExpress/internal/consts"
)

var OKRespose = ErrorJson{
	Message: "OK",
}

type ErrorJson struct {
	Message string `json:"message"`
}

type ErrorResponse struct {
	Error      error
	UserError  *ErrorJson
	StatusCode int
}

func NewErrorResponse(errConst int, err error) *ErrorResponse {
	return &ErrorResponse{
		Error: err,
		UserError: &ErrorJson{
			Message: Errors[errConst].Error(),
		},
		StatusCode: StatusCodes[errConst],
	}
}

func NewCustomErrorResponse(status int, err error, customText string) *ErrorResponse {
	return &ErrorResponse{
		Error: err,
		UserError: &ErrorJson{
			Message: customText,
		},
		StatusCode: status,
	}
}
