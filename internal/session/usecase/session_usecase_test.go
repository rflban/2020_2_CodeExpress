package usecase

import (
	"context"
	"github.com/go-park-mail-ru/2020_2_CodeExpress/internal/models"
	"github.com/go-park-mail-ru/2020_2_CodeExpress/internal/session/grpc_session"
	"github.com/go-park-mail-ru/2020_2_CodeExpress/internal/session/mock_session"
	"github.com/go-park-mail-ru/2020_2_CodeExpress/internal/session/proto_session"
	"github.com/go-playground/assert/v2"
	"github.com/golang/mock/gomock"
	"testing"
	"time"
)

func TestSessionUsecase_GetByID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockClient := mock_session.NewMockSessionServiceClient(ctrl)
	mockUsecase := NewSessionUsecase(mockClient)

	id := "1"
	sessionID := &proto_session.SessionID{
		ID: id,
	}

	mockClient.
		EXPECT().
		GetByID(context.Background(), sessionID).
		Return(&proto_session.Session{ ID: id }, nil)

	_, err := mockUsecase.GetByID(id)
	assert.Equal(t, err, nil)
}

func TestSessionUsecase_CreateSession(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockClient := mock_session.NewMockSessionServiceClient(ctrl)
	mockUsecase := NewSessionUsecase(mockClient)

	 session := &models.Session{
		ID:     "1",
		UserID: 1,
		Expire: time.Now(),
	}

	mockClient.
		EXPECT().
		CreateSession(context.Background(), grpc_session.SessionToGRPCSession(session)).
		Return(nil, nil)

	err := mockUsecase.CreateSession(session)
	assert.Equal(t, err, nil)
}

func TestSessionUsecase_DeleteSession(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockClient := mock_session.NewMockSessionServiceClient(ctrl)
	mockUsecase := NewSessionUsecase(mockClient)

	session := &models.Session{
		ID:     "1",
		UserID: 1,
		Expire: time.Now(),
	}

	mockClient.
		EXPECT().
		DeleteSession(context.Background(), grpc_session.SessionToGRPCSession(session)).
		Return(nil, nil)

	err := mockUsecase.DeleteSession(session)
	assert.Equal(t, err, nil)
}
