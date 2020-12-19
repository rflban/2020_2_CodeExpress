package user

import "github.com/go-park-mail-ru/2020_2_CodeExpress/internal/models"

type UserRep interface {
	Insert(name string, email string, password string) (*models.User, error)
	Update(user *models.User) error
	SelectById(id uint64) (*models.User, error)
	SelectByLogin(login string) (*models.User, error)
	SelectByNameOrEmail(name string, email string) ([]*models.User, error)
	SelectByName(name string, authUserId uint64) (*models.User, error)
	SelectIfAdmin(userID uint64) (bool, error)
	InsertSubscription(userSubscriberId uint64, userName string) error
	RemoveSubscription(userSubscriberId uint64, userName string) error
	SelectSubscriptions(id, authUserId uint64) (*models.Subscriptions, error)
}
