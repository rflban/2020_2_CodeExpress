package delivery_test

import (
	"encoding/json"
	. "github.com/go-park-mail-ru/2020_2_CodeExpress/internal/consts"
	"github.com/go-park-mail-ru/2020_2_CodeExpress/internal/models"
	"github.com/go-park-mail-ru/2020_2_CodeExpress/internal/session/mock_session"
	"github.com/go-park-mail-ru/2020_2_CodeExpress/internal/tools/builder"
	"github.com/go-park-mail-ru/2020_2_CodeExpress/internal/user/delivery"
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

func TestUserDelivery_HandlerRegisterUser(t *testing.T) {
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
	avatar := ""

	request := &Request{
		Name:   name,
		Email: email,
		Password:   password,
		RepeatedPassword:   password,
	}

	expectedUser := &models.User{
		ID:      id,
		Name:   name,
		Email: email,
		Password:   password,
		Avatar: avatar,
	}

	userMockUsecase.
		EXPECT().
		Create(gomock.Eq(name), gomock.Eq(email), gomock.Eq(password)).
		Return(expectedUser, nil)

	sessionMockUsecase.
		EXPECT().
		CreateSession(gomock.Any()).
		Return(nil)

	userHandler := delivery.NewUserHandler(userMockUsecase, sessionMockUsecase)
	e := echo.New()
	userHandler.Configure(e, nil)

	jsonRequest, err := json.Marshal(request)
	assert.Equal(t, err, nil)

	jsonExpectedUser, err := json.Marshal(expectedUser)
	assert.Equal(t, err, nil)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/user", strings.NewReader(string(jsonRequest)))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	resWriter := httptest.NewRecorder()
	ctx := e.NewContext(req, resWriter)

	handler := userHandler.HandlerRegisterUser()
	err = handler(ctx)
	assert.Equal(t, err, nil)
	assert.Equal(t, http.StatusOK, resWriter.Code)

	resBody, err := ioutil.ReadAll(resWriter.Body)
	assert.Equal(t, err, nil)
	clearBody := resBody[:len(resBody)-1]
	assert.Equal(t, clearBody, jsonExpectedUser)
}

func TestUserDelivery_HandlerRegisterUserFailed(t *testing.T) {
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

	name := "somename"
	email := "someemail@mail.ru"
	password := "short"

	request := &Request{
		Name:   name,
		Email: email,
		Password:   password,
		RepeatedPassword:   password,
	}

	userHandler := delivery.NewUserHandler(userMockUsecase, sessionMockUsecase)
	e := echo.New()
	userHandler.Configure(e, nil)

	jsonRequest, err := json.Marshal(request)
	assert.Equal(t, err, nil)

	assert.Equal(t, err, nil)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/user", strings.NewReader(string(jsonRequest)))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	resWriter := httptest.NewRecorder()
	ctx := e.NewContext(req, resWriter)

	handler := userHandler.HandlerRegisterUser()
	err = handler(ctx)
	assert.Equal(t, err, nil)
	assert.Equal(t, http.StatusBadRequest, resWriter.Code)
}

func TestUserDelivery_HandlerCurrentUserInfo(t *testing.T) {
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
	avatar := ""

	cookieValue := "Some cookie value"

	request := &Request{
		Name:   name,
		Email: email,
		Password:   password,
		RepeatedPassword:   password,
	}

	expectedUser := &models.User{
		ID:      id,
		Name:   name,
		Email: email,
		Password:   password,
		Avatar: avatar,
	}

	session := &models.Session{
		ID:     cookieValue,
		UserID: id,
		Name:   ConstSessionName,
	}

	userMockUsecase.
		EXPECT().
		GetById(session.UserID).
		Return(expectedUser, nil)

	sessionMockUsecase.
		EXPECT().
		GetByID(cookieValue).
		Return(session, nil)

	userHandler := delivery.NewUserHandler(userMockUsecase, sessionMockUsecase)
	e := echo.New()
	userHandler.Configure(e, nil)

	jsonRequest, err := json.Marshal(request)
	assert.Equal(t, err, nil)

	jsonExpectedUser, err := json.Marshal(expectedUser)
	assert.Equal(t, err, nil)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/user", strings.NewReader(string(jsonRequest)))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	req.AddCookie(builder.BuildCookie(session))
	resWriter := httptest.NewRecorder()
	ctx := e.NewContext(req, resWriter)

	handler := userHandler.HandlerCurrentUserInfo()
	err = handler(ctx)
	assert.Equal(t, err, nil)
	assert.Equal(t, http.StatusOK, resWriter.Code)

	resBody, err := ioutil.ReadAll(resWriter.Body)
	assert.Equal(t, err, nil)
	clearBody := resBody[:len(resBody)-1]
	assert.Equal(t, clearBody, jsonExpectedUser)
}

func TestUserDelivery_HandlerUpdateProfile(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	userMockUsecase := mock_user.NewMockUserUsecase(ctrl)
	sessionMockUsecase := mock_session.NewMockSessionUsecase(ctrl)

	type Request struct {
		Name  string `json:"username" validate:"required"`
		Email string `json:"email" validate:"required,email"`
	}

	id := uint64(1)
	name := "somename"
	email := "someemail@mail.ru"
	password := "somepassword"
	avatar := ""

	cookieValue := "Some cookie value"

	request := &Request{
		Name:   name,
		Email: email,
	}

	expectedUser := &models.User{
		ID:      id,
		Name:   name,
		Email: email,
		Password:   password,
		Avatar: avatar,
	}

	session := &models.Session{
		ID:     cookieValue,
		UserID: id,
		Name:   ConstSessionName,
	}

	userMockUsecase.
		EXPECT().
		UpdateProfile(session.UserID, name, email).
		Return(expectedUser, nil)

	sessionMockUsecase.
		EXPECT().
		GetByID(cookieValue).
		Return(session, nil)

	userHandler := delivery.NewUserHandler(userMockUsecase, sessionMockUsecase)
	e := echo.New()
	userHandler.Configure(e, nil)

	jsonRequest, err := json.Marshal(request)
	assert.Equal(t, err, nil)

	jsonExpectedUser, err := json.Marshal(expectedUser)
	assert.Equal(t, err, nil)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/user/profile", strings.NewReader(string(jsonRequest)))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	req.AddCookie(builder.BuildCookie(session))
	resWriter := httptest.NewRecorder()
	ctx := e.NewContext(req, resWriter)

	handler := userHandler.HandlerUpdateProfile()
	err = handler(ctx)
	assert.Equal(t, err, nil)
	assert.Equal(t, http.StatusOK, resWriter.Code)

	resBody, err := ioutil.ReadAll(resWriter.Body)
	assert.Equal(t, err, nil)
	clearBody := resBody[:len(resBody)-1]
	assert.Equal(t, clearBody, jsonExpectedUser)
}

func TestUserDelivery_HandlerUpdateProfileFailed(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	userMockUsecase := mock_user.NewMockUserUsecase(ctrl)
	sessionMockUsecase := mock_session.NewMockSessionUsecase(ctrl)

	type Request struct {
		Name  string `json:"username" validate:"required"`
		Email string `json:"email" validate:"required,email"`
	}

	id := uint64(1)
	name := "somename"
	email := "some wrong email"

	cookieValue := "Some cookie value"

	request := &Request{
		Name:   name,
		Email: email,
	}

	session := &models.Session{
		ID:     cookieValue,
		UserID: id,
		Name:   ConstSessionName,
	}

	userHandler := delivery.NewUserHandler(userMockUsecase, sessionMockUsecase)
	e := echo.New()
	userHandler.Configure(e, nil)

	jsonRequest, err := json.Marshal(request)
	assert.Equal(t, err, nil)


	req := httptest.NewRequest(http.MethodPost, "/api/v1/user/profile", strings.NewReader(string(jsonRequest)))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	req.AddCookie(builder.BuildCookie(session))
	resWriter := httptest.NewRecorder()
	ctx := e.NewContext(req, resWriter)

	handler := userHandler.HandlerUpdateProfile()
	err = handler(ctx)
	assert.Equal(t, err, nil)
	assert.Equal(t, http.StatusBadRequest, resWriter.Code)
}

func TestUserDelivery_HandlerUpdatePassword(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	userMockUsecase := mock_user.NewMockUserUsecase(ctrl)
	sessionMockUsecase := mock_session.NewMockSessionUsecase(ctrl)

	type Request struct {
		OldPassword      string `json:"old_password" validate:"required"`
		Password         string `json:"password" validate:"required,gte=8"`
		RepeatedPassword string `json:"repeated_password" validate:"required,eqfield=Password"`
	}

	id := uint64(1)
	password := "somepassword"
	newPassword := "somenewpassword"

	cookieValue := "Some cookie value"

	request := &Request{
		OldPassword:      password,
		Password:         newPassword,
		RepeatedPassword: newPassword,
	}

	session := &models.Session{
		ID:     cookieValue,
		UserID: id,
		Name:   ConstSessionName,
	}

	userMockUsecase.
		EXPECT().
		UpdatePassword(session.UserID, password, newPassword).
		Return(nil)

	sessionMockUsecase.
		EXPECT().
		GetByID(cookieValue).
		Return(session, nil)

	userHandler := delivery.NewUserHandler(userMockUsecase, sessionMockUsecase)
	e := echo.New()
	userHandler.Configure(e, nil)

	jsonRequest, err := json.Marshal(request)
	assert.Equal(t, err, nil)

	req := httptest.NewRequest(http.MethodPut, "/api/v1/user/password", strings.NewReader(string(jsonRequest)))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	req.AddCookie(builder.BuildCookie(session))
	resWriter := httptest.NewRecorder()
	ctx := e.NewContext(req, resWriter)

	handler := userHandler.HandlerUpdatePassword()
	err = handler(ctx)
	assert.Equal(t, err, nil)
	assert.Equal(t, http.StatusOK, resWriter.Code)
}

func TestUserDelivery_HandlerUpdatePasswordFailed(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	userMockUsecase := mock_user.NewMockUserUsecase(ctrl)
	sessionMockUsecase := mock_session.NewMockSessionUsecase(ctrl)

	type Request struct {
		OldPassword      string `json:"old_password" validate:"required"`
		Password         string `json:"password" validate:"required,gte=8"`
		RepeatedPassword string `json:"repeated_password" validate:"required,eqfield=Password"`
	}

	id := uint64(1)
	password := "somepassword"
	newPassword := "somenewpassword"

	cookieValue := "Some cookie value"

	request := &Request{
		OldPassword:      password,
		Password:         newPassword,
		RepeatedPassword: password,
	}

	session := &models.Session{
		ID:     cookieValue,
		UserID: id,
		Name:   ConstSessionName,
	}

	userHandler := delivery.NewUserHandler(userMockUsecase, sessionMockUsecase)
	e := echo.New()
	userHandler.Configure(e, nil)

	jsonRequest, err := json.Marshal(request)
	assert.Equal(t, err, nil)

	req := httptest.NewRequest(http.MethodPut, "/api/v1/user/password", strings.NewReader(string(jsonRequest)))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	req.AddCookie(builder.BuildCookie(session))
	resWriter := httptest.NewRecorder()
	ctx := e.NewContext(req, resWriter)

	handler := userHandler.HandlerUpdatePassword()
	err = handler(ctx)
	assert.Equal(t, err, nil)
	assert.Equal(t, http.StatusBadRequest, resWriter.Code)
}

func TestUserDelivery_HandlerUpdateAvatar(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	userMockUsecase := mock_user.NewMockUserUsecase(ctrl)
	sessionMockUsecase := mock_session.NewMockSessionUsecase(ctrl)

	type Request struct {
		OldPassword      string `json:"old_password" validate:"required"`
		Password         string `json:"password" validate:"required,gte=8"`
		RepeatedPassword string `json:"repeated_password" validate:"required,eqfield=Password"`
	}

	id := uint64(1)
	name := "somename"
	email := "someemail@mail.ru"
	password := "somepassword"
	newPassword := "somenewpassword"
	avatar := "pathToNewFile"

	cookieValue := "Some cookie value"

	request := &Request{
		OldPassword:      password,
		Password:         newPassword,
		RepeatedPassword: newPassword,
	}

	expectedUser := &models.User{
		ID:      id,
		Name:   name,
		Email: email,
		Password:   password,
		Avatar: avatar,
	}

	session := &models.Session{
		ID:     cookieValue,
		UserID: id,
		Name:   ConstSessionName,
	}

	userMockUsecase.
		EXPECT().
		GetById(session.UserID).
		Return(expectedUser, nil)

	sessionMockUsecase.
		EXPECT().
		GetByID(cookieValue).
		Return(session, nil)

	userHandler := delivery.NewUserHandler(userMockUsecase, sessionMockUsecase)
	e := echo.New()
	userHandler.Configure(e, nil)

	jsonRequest, err := json.Marshal(request)
	assert.Equal(t, err, nil)

	req := httptest.NewRequest(http.MethodPut, "/api/v1/user/photo", strings.NewReader(string(jsonRequest)))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	req.AddCookie(builder.BuildCookie(session))
	resWriter := httptest.NewRecorder()
	ctx := e.NewContext(req, resWriter)

	handler := userHandler.HandlerUpdateAvatar()
	err = handler(ctx)
	assert.Equal(t, err, nil)
	assert.Equal(t, http.StatusBadRequest, resWriter.Code)
}
