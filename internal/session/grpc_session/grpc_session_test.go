package grpc_session_test

import (
	"context"
	"testing"
	"time"

	"github.com/go-park-mail-ru/2020_2_CodeExpress/internal/session/proto_session"

	"github.com/go-park-mail-ru/2020_2_CodeExpress/internal/session/grpc_session"

	"github.com/go-park-mail-ru/2020_2_CodeExpress/internal/models"
	"github.com/go-park-mail-ru/2020_2_CodeExpress/internal/session/mock_session"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestSessionUsecase_GetByID_Passed(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock_session.NewMockSessionRep(ctrl)
	mockGRPC := grpc_session.NewSessionGRPCUsecase(mockRepo)

	id := "some id"

	expectedSession := &models.Session{
		ID:     id,
		UserID: 1,
		Expire: time.Now(),
	}

	mockRepo.
		EXPECT().
		SelectById(id).
		Return(expectedSession, nil)

	_, err := mockGRPC.GetByID(context.Background(), &proto_session.SessionID{
		ID: id,
	})
	assert.Nil(t, err)
}

func TestSessionUsecase_CreateSession_Passed(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock_session.NewMockSessionRep(ctrl)

	mockGRPC := grpc_session.NewSessionGRPCUsecase(mockRepo)

	id := "some id"

	session := &models.Session{
		ID:     id,
		UserID: 1,
	}

	mockRepo.
		EXPECT().
		Insert(session).
		Return(nil)

	_, err := mockGRPC.CreateSession(context.Background(), grpc_session.SessionToGRPCSession(session))
	assert.Nil(t, err)
}

func TestSessionUsecase_DeleteSession_Passed(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock_session.NewMockSessionRep(ctrl)

	mockGRPC := grpc_session.NewSessionGRPCUsecase(mockRepo)

	id := "some id"

	session := &models.Session{
		ID:     id,
		UserID: 1,
	}

	mockRepo.
		EXPECT().
		Delete(session).
		Return(nil)

	_, err := mockGRPC.DeleteSession(context.Background(), grpc_session.SessionToGRPCSession(session))
	assert.Nil(t, err)
}
