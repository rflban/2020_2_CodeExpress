package user

import (
	"github.com/go-park-mail-ru/2020_2_CodeExpress/internal/models"
	. "github.com/go-park-mail-ru/2020_2_CodeExpress/internal/tools/error_response"
)

type UserUsecase interface {
	CreateUser(user *models.User) *ErrorResponse
	GetByEmail(email string) (*models.User, *ErrorResponse)
	GetByName(name string) (*models.User, *ErrorResponse)
	GetById(id uint64) (*models.User, *ErrorResponse)
	LoginUser(name string, password string) (*models.User, *ErrorResponse)
	UpdateProfile(user *models.User) *ErrorResponse
	UpdatePassword(userId uint64, oldPassword string, password string) (*models.User, *ErrorResponse)
}
