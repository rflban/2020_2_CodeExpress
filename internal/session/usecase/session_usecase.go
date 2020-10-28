package usecase

import (
	"database/sql"
	"time"

	. "github.com/go-park-mail-ru/2020_2_CodeExpress/internal/consts"
	"github.com/go-park-mail-ru/2020_2_CodeExpress/internal/models"
	"github.com/go-park-mail-ru/2020_2_CodeExpress/internal/session"
	. "github.com/go-park-mail-ru/2020_2_CodeExpress/internal/tools/error_response"
)

type SessionUsecase struct {
	sessionRep session.SessionRep
}

func NewSessionUsecase(sessionRep session.SessionRep) session.SessionUsecase {
	return &SessionUsecase{
		sessionRep: sessionRep,
	}
}

func (sUc *SessionUsecase) GetByID(id string) (*models.Session, *ErrorResponse) {
	session, err := sUc.sessionRep.SelectById(id)

	if err == sql.ErrNoRows {
		return nil, NewErrorResponse(ErrNotAuthorized, err)
	}

	if err != nil {
		return nil, NewErrorResponse(ErrInternal, err)
	}

	return session, nil
}

func (sUc *SessionUsecase) CreateSession(session *models.Session) *ErrorResponse {
	if err := sUc.sessionRep.Insert(session); err != nil {
		return NewErrorResponse(ErrInternal, err)
	}
	return nil
}

func (sUc *SessionUsecase) DeleteSession(session *models.Session) *ErrorResponse {
	if err := sUc.sessionRep.Delete(session); err != nil {
		return NewErrorResponse(ErrInternal, err)
	}

	session.Expire = time.Now().AddDate(0, 0, -ConstDaysSession)

	return nil
}
