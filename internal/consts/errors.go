package consts

import (
	"errors"
	"net/http"
)

const (
	InternalError      = "internal server error"
	FormError          = "incorrect form"
	PasswordTooShort   = "password too short, at least 8 letters"
	NoError            = "ok"
	NoEmail            = "no email field"
	NoUsername         = "no username field"
	NoPassword         = "no password field"
	NoRepeatedPassword = "no repeated password field"
	PasswordsMismatch  = "passwords do not match"
	PasswordIsOld      = "password was not changed"
	NotAuthorized      = "not authorized"
	NoAvatar           = "avatar is expected"
	FileError          = "error reading file"
	FileSizeToLarge    = "file size is to large"
)

const (
	ErrInternal = iota
	ErrBadRequest
	ErrEmailAlreadyExist
	ErrNameAlreadyExist
	ErrIncorrectLoginOrPassword
	ErrNotAuthorized
	ErrNoEmail
	ErrNoUsername
	ErrNoAvatar
	ErrWrongOldPassword
	ErrNewPasswordIsOld
)

var Errors = map[int]error{
	ErrInternal:                 errors.New("Internal server error"),
	ErrBadRequest:               errors.New("Bad request received"),
	ErrEmailAlreadyExist:        errors.New("Email already exists"),
	ErrNameAlreadyExist:         errors.New("Name already exists"),
	ErrIncorrectLoginOrPassword: errors.New("Incorrect login or password"),
	ErrNotAuthorized:            errors.New("Not authorized"),
	ErrNoEmail:                  errors.New("No email field"),
	ErrNoUsername:               errors.New("No username field"),
	ErrNoAvatar:                 errors.New("No avatar field"),
	ErrWrongOldPassword:         errors.New("Wrong old password"),
	ErrNewPasswordIsOld:         errors.New("New password matches old"),
}

var StatusCodes = map[int]int{
	ErrInternal:                 http.StatusInternalServerError,
	ErrBadRequest:               http.StatusBadRequest,
	ErrEmailAlreadyExist:        http.StatusForbidden,
	ErrNameAlreadyExist:         http.StatusForbidden,
	ErrIncorrectLoginOrPassword: http.StatusNotFound,
	ErrNotAuthorized:            http.StatusNotFound,
	ErrNoEmail:                  http.StatusBadRequest,
	ErrNoUsername:               http.StatusBadRequest,
	ErrNoAvatar:                 http.StatusBadRequest,
	ErrWrongOldPassword:         http.StatusBadRequest,
	ErrNewPasswordIsOld:         http.StatusBadRequest,
}
