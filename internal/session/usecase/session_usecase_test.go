package usecase_test

import (
	"github.com/go-park-mail-ru/2020_2_CodeExpress/internal/models"
	"github.com/go-park-mail-ru/2020_2_CodeExpress/internal/session/mock_session"
	"github.com/go-park-mail-ru/2020_2_CodeExpress/internal/session/usecase"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestSessionUsecase_GetByID_Passed(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock_session.NewMockSessionRep(ctrl)
	mockUsecase := usecase.NewSessionUsecase(mockRepo)

	id := "some id"

	expectedSession := &models.Session{
		ID: id,
		UserID: 1,
		Expire: time.Now(),
	}

	mockRepo.
		EXPECT().
		SelectById(id).
		Return(expectedSession, nil)

	session, err := mockUsecase.GetByID(id)
	assert.Nil(t, err)
	assert.Equal(t, expectedSession, session)
}

func TestSessionUsecase_CreateSession_Passed(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock_session.NewMockSessionRep(ctrl)
	mockUsecase := usecase.NewSessionUsecase(mockRepo)

	id := "some id"

	session := &models.Session{
		ID: id,
		UserID: 1,
		Expire: time.Now(),
	}

	mockRepo.
		EXPECT().
		Insert(session).
		Return(nil)

	err := mockUsecase.CreateSession(session)
	assert.Nil(t, err)
}

func TestSessionUsecase_DeleteSession_Passed(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock_session.NewMockSessionRep(ctrl)
	mockUsecase := usecase.NewSessionUsecase(mockRepo)

	id := "some id"

	session := &models.Session{
		ID: id,
		UserID: 1,
		Expire: time.Now(),
	}

	mockRepo.
		EXPECT().
		Delete(session).
		Return(nil)

	err := mockUsecase.DeleteSession(session)
	assert.Nil(t, err)
}
