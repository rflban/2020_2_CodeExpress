package repositories

import (
	"errors"
	"log"
	"sync"

	"github.com/go-park-mail-ru/2020_2_CodeExpress/models"
)

type SignUpRep interface {
	CheckUserExists(u *models.NewUser) error
	CreateUser(u *models.NewUser) (*models.User, error)
}

type SignUpRepImpl struct {
	Users []*models.User
	MU    *sync.RWMutex
}

func NewSignUpRepImpl() *SignUpRepImpl {
	return &SignUpRepImpl{
		Users: make([]*models.User, 0),
		MU:    &sync.RWMutex{},
	}
}

func (s *SignUpRepImpl) CheckUserExists(u *models.NewUser) error {
	s.MU.RLock()
	defer s.MU.RUnlock()
	for _, user := range s.Users {
		if user.Email == u.Email {
			return errors.New("Email already exists")
		}
		if user.Name == u.Name {
			return errors.New("Username already exists")
		}
	}
	return nil
}

func (s *SignUpRepImpl) CreateUser(u *models.NewUser) (*models.User, error) {
	s.MU.Lock()
	defer s.MU.Unlock()
	user := &models.User{
		ID:       s.getLastUserID() + 1,
		Name:     u.Name,
		Email:    u.Email,
		Password: u.Password,
	}
	s.Users = append(s.Users, user)
	log.Println("New user: ", user)
	return user, nil // возвращает nil так как реализация без БД
}

func (s *SignUpRepImpl) getLastUserID() uint64 {
	if len(s.Users) > 0 {
		return s.Users[len(s.Users)-1].ID
	}
	return 0
}
