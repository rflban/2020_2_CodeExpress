package session

import (
	"github.com/go-park-mail-ru/2020_2_CodeExpress/internal/models"
	. "github.com/go-park-mail-ru/2020_2_CodeExpress/internal/tools/error_response"
)

type SessionUsecase interface {
	CreateSession(session *models.Session) *ErrorResponse
	GetByID(id string) (*models.Session, *ErrorResponse)
	DeleteSession(session *models.Session) *ErrorResponse
}
