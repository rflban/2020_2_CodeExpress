package user

import (
	"github.com/go-park-mail-ru/2020_2_CodeExpress/internal/models"
	. "github.com/go-park-mail-ru/2020_2_CodeExpress/internal/tools/error_response"
)

type UserUsecase interface {
	CreateUser(user *models.User) *ErrorResponse
	GetByEmail(email string) (*models.User, *ErrorResponse)
	GetByName(name string) (*models.User, *ErrorResponse)
	GetByID(id uint64) (*models.User, *ErrorResponse)
	LoginUser(login string, password string) (*models.User, *ErrorResponse)
}
