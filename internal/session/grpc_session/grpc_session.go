package grpc_session

import (
	"github.com/go-park-mail-ru/2020_2_CodeExpress/internal/models"
	"github.com/go-park-mail-ru/2020_2_CodeExpress/internal/session"
	proto_session "github.com/go-park-mail-ru/2020_2_CodeExpress/internal/session/proto_session"
	"golang.org/x/net/context"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type SessionGRPCUsecase struct {
	sessionRep session.SessionRep
}

func SessionToGRPCSession(session *models.Session) *proto_session.Session {
	return &proto_session.Session{
		ID:     session.ID,
		Name:   session.Name,
		UserID: session.UserID,
		Expire: timestamppb.New(session.Expire),
	}
}

func (sgu *SessionGRPCUsecase) CreateSession(ctx context.Context, grpcSession *proto_session.Session) (*proto_session.Session, error) {
	session := GRPCSessionToSession(grpcSession)

	err := sgu.sessionRep.Insert(session)

	if err != nil {
		return new(proto_session.Session), err
	}

	return grpcSession, nil
}

func (sgu *SessionGRPCUsecase) DeleteSession(ctx context.Context, grpcSession *proto_session.Session) (*proto_session.Session, error) {
	session := GRPCSessionToSession(grpcSession)

	err := sgu.sessionRep.Delete(session)

	if err != nil {
		return new(proto_session.Session), err
	}

	return grpcSession, nil
}

func (sgu *SessionGRPCUsecase) GetByID(ctx context.Context, grpcSession *proto_session.SessionID) (*proto_session.Session, error) {
	sessionID := grpcSession.ID

	session, err := sgu.sessionRep.SelectById(sessionID)

	if err != nil {
		return new(proto_session.Session), err
	}

	return SessionToGRPCSession(session), nil
}

func NewSessionGRPCUsecase(sessionRep session.SessionRep) *SessionGRPCUsecase {
	return &SessionGRPCUsecase{
		sessionRep: sessionRep,
	}
}

func GRPCSessionToSession(session *proto_session.Session) *models.Session {
	return &models.Session{
		ID:     session.ID,
		Name:   session.Name,
		UserID: session.UserID,
		Expire: session.Expire.AsTime(),
	}
}
