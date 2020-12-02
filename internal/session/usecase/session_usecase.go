package usecase

import (
	"database/sql"
	"time"

	"golang.org/x/net/context"

	"github.com/go-park-mail-ru/2020_2_CodeExpress/internal/session/grpc_session"
	"github.com/go-park-mail-ru/2020_2_CodeExpress/internal/session/proto_session"

	. "github.com/go-park-mail-ru/2020_2_CodeExpress/internal/consts"
	"github.com/go-park-mail-ru/2020_2_CodeExpress/internal/models"
	"github.com/go-park-mail-ru/2020_2_CodeExpress/internal/session"
	. "github.com/go-park-mail-ru/2020_2_CodeExpress/internal/tools/error_response"
)

type SessionUsecase struct {
	sessionGRPC proto_session.SessionServiceClient
}

func NewSessionUsecase(sessionGRPC proto_session.SessionServiceClient) session.SessionUsecase {
	return &SessionUsecase{
		sessionGRPC: sessionGRPC,
	}
}

func (sUc *SessionUsecase) GetByID(id string) (*models.Session, *ErrorResponse) {
	sessionID := &proto_session.SessionID{
		ID: id,
	}

	session, err := sUc.sessionGRPC.GetByID(context.Background(), sessionID)

	if err == sql.ErrNoRows {
		return nil, NewErrorResponse(ErrNotAuthorized, err)
	}

	if err != nil {
		return nil, NewErrorResponse(ErrInternal, err)
	}

	return grpc_session.GRPCSessionToSession(session), nil
}

func (sUc *SessionUsecase) CreateSession(session *models.Session) *ErrorResponse {
	if _, err := sUc.sessionGRPC.CreateSession(context.Background(), grpc_session.SessionToGRPCSession(session)); err != nil {
		return NewErrorResponse(ErrInternal, err)
	}
	return nil
}

func (sUc *SessionUsecase) DeleteSession(session *models.Session) *ErrorResponse {
	if _, err := sUc.sessionGRPC.DeleteSession(context.Background(), grpc_session.SessionToGRPCSession(session)); err != nil {
		return NewErrorResponse(ErrInternal, err)
	}

	session.Expire = time.Now().AddDate(0, 0, -ConstDaysSession)

	return nil
}
