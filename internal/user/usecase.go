package user

import (
	"github.com/go-park-mail-ru/2020_2_CodeExpress/internal/models"
	. "github.com/go-park-mail-ru/2020_2_CodeExpress/internal/tools/error_response"
)

type UserUsecase interface {
	Create(string, string, string) (*models.User, *ErrorResponse)
	GetUserByLogin(string, string) (*models.User, *ErrorResponse)
	GetById(uint64) (*models.User, *ErrorResponse)
	UpdateProfile(uint64, string, string) (*models.User, *ErrorResponse)
	UpdatePassword(uint64, string, string) *ErrorResponse
	UpdateAvatar(uint64, string) (*models.User, *ErrorResponse)
}
