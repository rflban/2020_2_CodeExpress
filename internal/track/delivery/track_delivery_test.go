package delivery_test

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/go-park-mail-ru/2020_2_CodeExpress/internal/session/mock_session"
	"github.com/go-park-mail-ru/2020_2_CodeExpress/internal/user/mock_user"

	. "github.com/go-park-mail-ru/2020_2_CodeExpress/internal/consts"
	"github.com/go-park-mail-ru/2020_2_CodeExpress/internal/tools/builder"

	"github.com/go-park-mail-ru/2020_2_CodeExpress/internal/track/delivery"

	"github.com/go-park-mail-ru/2020_2_CodeExpress/internal/track/mock_track"

	"github.com/go-park-mail-ru/2020_2_CodeExpress/internal/models"
	"github.com/go-playground/assert/v2"
	"github.com/golang/mock/gomock"
	"github.com/labstack/echo/v4"
)

func TestAlbumDelivery_HandlerCreateTrack(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	trackMockUsecase := mock_track.NewMockTrackUsecase(ctrl)

	type Request struct {
		Title   string `json:"title" validate:"required"`
		AlbumID uint64 `json:"album_id" validate:"required"`
	}

	id := uint64(42)
	index := uint8(32)
	title := "Some title"
	albumID := uint64(3)

	request := &Request{
		Title:   title,
		AlbumID: albumID,
	}

	track := &models.Track{
		Title:   title,
		AlbumID: albumID,
	}

	expectedTrack := &models.Track{
		ID:      id,
		Title:   title,
		AlbumID: albumID,
		Index:   index,
	}

	trackMockUsecase.
		EXPECT().
		CreateTrack(gomock.Eq(track)).
		DoAndReturn(func(track *models.Track) error {
			track.ID = id
			track.Index = index
			return nil
		})

	albumHandler := delivery.NewTrackHandler(trackMockUsecase, nil, nil)
	e := echo.New()
	albumHandler.Configure(e)

	jsonRequest, err := json.Marshal(request)
	assert.Equal(t, err, nil)

	jsonExpectedAlbum, err := json.Marshal(expectedTrack)
	assert.Equal(t, err, nil)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/tracks", strings.NewReader(string(jsonRequest)))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	resWriter := httptest.NewRecorder()
	ctx := e.NewContext(req, resWriter)

	handler := albumHandler.HandlerCreateTrack()

	err = handler(ctx)
	assert.Equal(t, err, nil)
	assert.Equal(t, http.StatusOK, resWriter.Code)

	resBody, err := ioutil.ReadAll(resWriter.Body)
	assert.Equal(t, err, nil)
	clearBody := resBody[:len(resBody)-1]
	assert.Equal(t, clearBody, jsonExpectedAlbum)
}

func TestAlbumDelivery_HandlerUpdateTrack(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	trackMockUsecase := mock_track.NewMockTrackUsecase(ctrl)

	type Request struct {
		Title   string `json:"title" validate:"required"`
		AlbumID uint64 `json:"album_id" validate:"required"`
		Index   uint8  `json:"index" validate:"required"`
	}

	id := uint64(42)
	index := uint8(32)
	title := "Some title"
	albumID := uint64(3)

	request := &Request{
		Title:   title,
		AlbumID: albumID,
		Index:   index,
	}

	expectedTrack := &models.Track{
		ID:      id,
		Title:   title,
		AlbumID: albumID,
		Index:   index,
	}

	trackMockUsecase.
		EXPECT().
		UpdateTrack(gomock.Eq(expectedTrack)).
		Return(nil)

	albumHandler := delivery.NewTrackHandler(trackMockUsecase, nil, nil)
	e := echo.New()
	albumHandler.Configure(e)

	jsonRequest, err := json.Marshal(request)
	assert.Equal(t, err, nil)

	req := httptest.NewRequest(http.MethodPut, "/api/v1/tracks/42", strings.NewReader(string(jsonRequest)))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	resWriter := httptest.NewRecorder()
	ctx := e.NewContext(req, resWriter)
	ctx.SetParamNames("id")
	ctx.SetParamValues("42")

	handler := albumHandler.HandlerUpdateTrack()

	err = handler(ctx)
	assert.Equal(t, err, nil)
	assert.Equal(t, http.StatusOK, resWriter.Code)
}

func TestAlbumDelivery_HandlerDeleteTrack(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	trackMockUsecase := mock_track.NewMockTrackUsecase(ctrl)

	id := uint64(42)

	trackMockUsecase.
		EXPECT().
		DeleteTrack(gomock.Eq(id)).
		Return(nil)

	albumHandler := delivery.NewTrackHandler(trackMockUsecase, nil, nil)
	e := echo.New()
	albumHandler.Configure(e)

	req := httptest.NewRequest(http.MethodDelete, "/api/v1/tracks/42", strings.NewReader(string("")))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	resWriter := httptest.NewRecorder()
	ctx := e.NewContext(req, resWriter)
	ctx.SetParamNames("id")
	ctx.SetParamValues("42")

	handler := albumHandler.HandlerDeleteTrack()

	err := handler(ctx)
	assert.Equal(t, err, nil)
	assert.Equal(t, http.StatusOK, resWriter.Code)
}

func TestAlbumDelivery_HandlerTracksByParams(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	trackMockUsecase := mock_track.NewMockTrackUsecase(ctrl)

	count, from := uint64(1), uint64(0)

	expectedTracks := []*models.Track{
		&models.Track{ID: 0},
		&models.Track{ID: 1},
	}

	trackMockUsecase.
		EXPECT().
		GetByParams(gomock.Eq(count), gomock.Eq(from)).
		Return(expectedTracks, nil)

	albumHandler := delivery.NewTrackHandler(trackMockUsecase, nil, nil)
	e := echo.New()
	albumHandler.Configure(e)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/tracks?count=1&from=0", strings.NewReader(string("")))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	resWriter := httptest.NewRecorder()
	ctx := e.NewContext(req, resWriter)

	handler := albumHandler.HandlerTracksByParams()

	err := handler(ctx)
	assert.Equal(t, err, nil)
	assert.Equal(t, http.StatusOK, resWriter.Code)
}

func TestAlbumDelivery_HandlerTracksByArtistID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	trackMockUsecase := mock_track.NewMockTrackUsecase(ctrl)

	id := uint64(42)

	expectedTracks := []*models.Track{
		&models.Track{ID: 0},
		&models.Track{ID: 1},
	}

	trackMockUsecase.
		EXPECT().
		GetByArtistID(id).
		Return(expectedTracks, nil)

	albumHandler := delivery.NewTrackHandler(trackMockUsecase, nil, nil)
	e := echo.New()
	albumHandler.Configure(e)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/artists/42/tracks", strings.NewReader(string("")))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	resWriter := httptest.NewRecorder()
	ctx := e.NewContext(req, resWriter)
	ctx.SetParamNames("id")
	ctx.SetParamValues("42")

	handler := albumHandler.HandlerTracksByArtistID()

	err := handler(ctx)
	assert.Equal(t, err, nil)
	assert.Equal(t, http.StatusOK, resWriter.Code)
}

func TestAlbumDelivery_HandlerFavouritesByUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	trackMockUsecase := mock_track.NewMockTrackUsecase(ctrl)
	sessionMockUsecase := mock_session.NewMockSessionUsecase(ctrl)
	userMockUsecase := mock_user.NewMockUserUsecase(ctrl)

	userID := uint64(42)
	cookieValue := "Some cookie value"

	expectedTracks := []*models.Track{
		&models.Track{ID: 0},
		&models.Track{ID: 1},
	}

	session := &models.Session{
		ID:     cookieValue,
		UserID: userID,
		Name:   ConstSessionName,
	}

	user := &models.User{ID: userID}

	sessionMockUsecase.
		EXPECT().
		GetByID(cookieValue).
		Return(session, nil)

	userMockUsecase.
		EXPECT().
		GetById(session.UserID).
		Return(user, nil)

	trackMockUsecase.
		EXPECT().
		GetFavoritesByUserID(user.ID).
		Return(expectedTracks, nil)

	albumHandler := delivery.NewTrackHandler(trackMockUsecase, sessionMockUsecase, userMockUsecase)
	e := echo.New()
	albumHandler.Configure(e)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/favorite/tracks", strings.NewReader(string("")))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	req.AddCookie(builder.BuildCookie(session))
	resWriter := httptest.NewRecorder()
	ctx := e.NewContext(req, resWriter)

	handler := albumHandler.HandlerFavouritesByUser()
	err := handler(ctx)
	assert.Equal(t, err, nil)
	assert.Equal(t, http.StatusOK, resWriter.Code)
}

func TestAlbumDelivery_HandlerAddToUsersFavourites(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	trackMockUsecase := mock_track.NewMockTrackUsecase(ctrl)
	sessionMockUsecase := mock_session.NewMockSessionUsecase(ctrl)
	userMockUsecase := mock_user.NewMockUserUsecase(ctrl)

	userID := uint64(42)
	trackID := uint64(156)
	cookieValue := "Some cookie value"

	session := &models.Session{
		ID:     cookieValue,
		UserID: userID,
		Name:   ConstSessionName,
	}

	user := &models.User{ID: userID}

	sessionMockUsecase.
		EXPECT().
		GetByID(cookieValue).
		Return(session, nil)

	userMockUsecase.
		EXPECT().
		GetById(session.UserID).
		Return(user, nil)

	trackMockUsecase.
		EXPECT().
		AddToFavourites(user.ID, trackID).
		Return(nil)

	albumHandler := delivery.NewTrackHandler(trackMockUsecase, sessionMockUsecase, userMockUsecase)
	e := echo.New()
	albumHandler.Configure(e)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/favorite/tracks/156", strings.NewReader(string("")))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	req.AddCookie(builder.BuildCookie(session))
	resWriter := httptest.NewRecorder()
	ctx := e.NewContext(req, resWriter)
	ctx.SetParamNames("id")
	ctx.SetParamValues("156")

	handler := albumHandler.HandlerAddToUsersFavourites()
	err := handler(ctx)
	assert.Equal(t, err, nil)
	assert.Equal(t, http.StatusOK, resWriter.Code)
}

func TestAlbumDelivery_HandlerDeleteFromUsersFavourites(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	trackMockUsecase := mock_track.NewMockTrackUsecase(ctrl)
	sessionMockUsecase := mock_session.NewMockSessionUsecase(ctrl)
	userMockUsecase := mock_user.NewMockUserUsecase(ctrl)

	userID := uint64(42)
	trackID := uint64(156)
	cookieValue := "Some cookie value"

	session := &models.Session{
		ID:     cookieValue,
		UserID: userID,
		Name:   ConstSessionName,
	}

	user := &models.User{ID: userID}

	sessionMockUsecase.
		EXPECT().
		GetByID(cookieValue).
		Return(session, nil)

	userMockUsecase.
		EXPECT().
		GetById(session.UserID).
		Return(user, nil)

	trackMockUsecase.
		EXPECT().
		DeleteFromFavourites(user.ID, trackID).
		Return(nil)

	albumHandler := delivery.NewTrackHandler(trackMockUsecase, sessionMockUsecase, userMockUsecase)
	e := echo.New()
	albumHandler.Configure(e)

	req := httptest.NewRequest(http.MethodDelete, "/api/v1/favorite/tracks/156", strings.NewReader(string("")))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	req.AddCookie(builder.BuildCookie(session))
	resWriter := httptest.NewRecorder()
	ctx := e.NewContext(req, resWriter)
	ctx.SetParamNames("id")
	ctx.SetParamValues("156")

	handler := albumHandler.HandlerDeleteFromUsersFavourites()
	err := handler(ctx)
	assert.Equal(t, err, nil)
	assert.Equal(t, http.StatusOK, resWriter.Code)
}
