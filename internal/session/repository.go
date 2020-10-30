package session

import (
	"github.com/go-park-mail-ru/2020_2_CodeExpress/internal/models"
)

type SessionRep interface {
	SelectById(id string) (*models.Session, error)
	Insert(*models.Session) error
	Delete(*models.Session) error
}
