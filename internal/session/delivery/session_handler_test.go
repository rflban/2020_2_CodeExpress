package delivery_test

import (
	"encoding/json"
	. "github.com/go-park-mail-ru/2020_2_CodeExpress/internal/consts"
	"github.com/go-park-mail-ru/2020_2_CodeExpress/internal/models"
	"github.com/go-park-mail-ru/2020_2_CodeExpress/internal/session/delivery"
	"github.com/go-park-mail-ru/2020_2_CodeExpress/internal/session/mock_session"
	"github.com/go-park-mail-ru/2020_2_CodeExpress/internal/tools/builder"
	. "github.com/go-park-mail-ru/2020_2_CodeExpress/internal/tools/error_response"
	"github.com/go-park-mail-ru/2020_2_CodeExpress/internal/user/mock_user"
	"github.com/go-playground/assert/v2"
	"github.com/golang/mock/gomock"
	"github.com/labstack/echo/v4"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestSessionDelivery_HandlerLogin(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	sessionMockUsecase := mock_session.NewMockSessionUsecase(ctrl)
	userMockUsecase := mock_user.NewMockUserUsecase(ctrl)

	type Request struct {
		Login    string `json:"login" validate:"required"`
		Password string `json:"password" validate:"required"`
	}

	login := "somelogin"
	password := "somepassword"

	request := &Request{
		Login: login,
		Password:   password,
	}

	id := uint64(1)
	name := "somename"
	email := "someemail@mail.ru"
	avatar := ""

	expectedUser := &models.User{
		ID:      id,
		Name:   name,
		Email: email,
		Password:   password,
		Avatar: avatar,
	}

	userMockUsecase.
		EXPECT().
		GetUserByLogin(gomock.Eq(login), gomock.Eq(password)).
		Return(expectedUser, nil)

	sessionMockUsecase.
		EXPECT().
		CreateSession(gomock.Any()).
		Return(nil)

	sessionHandler := delivery.NewSessionHandler(sessionMockUsecase, userMockUsecase)
	e := echo.New()
	sessionHandler.Configure(e)

	jsonRequest, err := json.Marshal(request)
	assert.Equal(t, err, nil)

	jsonExpectedUser, err := json.Marshal(expectedUser)
	assert.Equal(t, err, nil)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/session", strings.NewReader(string(jsonRequest)))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	resWriter := httptest.NewRecorder()
	ctx := e.NewContext(req, resWriter)

	handler := sessionHandler.HandlerLogin()
	err = handler(ctx)
	assert.Equal(t, err, nil)
	assert.Equal(t, http.StatusOK, resWriter.Code)

	resBody, err := ioutil.ReadAll(resWriter.Body)
	assert.Equal(t, err, nil)
	clearBody := resBody[:len(resBody)-1]
	assert.Equal(t, clearBody, jsonExpectedUser)
}

func TestSessionDelivery_HandlerLogout(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	userMockUsecase := mock_user.NewMockUserUsecase(ctrl)
	sessionMockUsecase := mock_session.NewMockSessionUsecase(ctrl)

	type Request struct {
		Name             string `json:"username" validate:"required"`
		Email            string `json:"email" validate:"required,email"`
		Password         string `json:"password" validate:"required,gte=8"`
		RepeatedPassword string `json:"repeated_password" validate:"required,eqfield=Password"`
	}

	id := uint64(1)
	name := "somename"
	email := "someemail@mail.ru"
	password := "somepassword"

	cookieValue := "Some cookie value"

	request := &Request{
		Name:   name,
		Email: email,
		Password:   password,
		RepeatedPassword:   password,
	}

	session := &models.Session{
		ID:     cookieValue,
		UserID: id,
		Name:   ConstSessionName,
	}

	sessionMockUsecase.
		EXPECT().
		GetByID(session.ID).
		Return(session, nil)

	sessionMockUsecase.
		EXPECT().
		DeleteSession(session).
		Return(nil)

	sessionHandler := delivery.NewSessionHandler(sessionMockUsecase, userMockUsecase)
	e := echo.New()
	sessionHandler.Configure(e)

	jsonRequest, err := json.Marshal(request)
	assert.Equal(t, err, nil)

	jsonResponse, err := json.Marshal(OKResponse)
	assert.Equal(t, err, nil)

	req := httptest.NewRequest(http.MethodDelete, "/api/v1/session", strings.NewReader(string(jsonRequest)))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	req.AddCookie(builder.BuildCookie(session))
	resWriter := httptest.NewRecorder()
	ctx := e.NewContext(req, resWriter)

	handler := sessionHandler.HandlerLogout()
	err = handler(ctx)
	assert.Equal(t, err, nil)
	assert.Equal(t, http.StatusOK, resWriter.Code)

	resBody, err := ioutil.ReadAll(resWriter.Body)
	assert.Equal(t, err, nil)
	clearBody := resBody[:len(resBody)-1]
	assert.Equal(t, clearBody, jsonResponse)
}
