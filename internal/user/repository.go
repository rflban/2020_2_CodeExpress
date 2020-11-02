package user

import "github.com/go-park-mail-ru/2020_2_CodeExpress/internal/models"

type UserRep interface {
	Insert(user *models.User) error
	Update(user *models.User) error
	SelectByEmail(email string) (*models.User, error)
	SelectByName(name string) (*models.User, error)
	SelectByNameOrEmail(name string, email string, id uint64) (*models.User, error)
	SelectByLoginAndPassword(login string, password string) (*models.User, error)
	SelectByID(userID uint64) (*models.User, error)
}
