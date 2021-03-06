package delivery_test

import (
	"encoding/json"
	"github.com/go-park-mail-ru/2020_2_CodeExpress/internal/session/mock_session"
	"github.com/go-park-mail-ru/2020_2_CodeExpress/internal/tools/builder"
	"github.com/go-park-mail-ru/2020_2_CodeExpress/internal/user/mock_user"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/go-park-mail-ru/2020_2_CodeExpress/internal/track/delivery"

	. "github.com/go-park-mail-ru/2020_2_CodeExpress/internal/consts"
	"github.com/go-park-mail-ru/2020_2_CodeExpress/internal/track/mock_track"

	"github.com/go-park-mail-ru/2020_2_CodeExpress/internal/models"
	"github.com/go-playground/assert/v2"
	"github.com/golang/mock/gomock"
	"github.com/labstack/echo/v4"
)

func TestTrackDelivery_HandlerCreateTrack(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	trackMockUsecase := mock_track.NewMockTrackUsecase(ctrl)

	type Request struct {
		Title   string `json:"title" validate:"required"`
		AlbumID uint64 `json:"album_id" validate:"required"`
	}

	track := &models.Track{
		Title:   "Some title",
		AlbumID: 1,
	}

	trackMockUsecase.
		EXPECT().
		CreateTrack(gomock.Eq(track), gomock.Eq(uint64(0))).
		Return(nil)

	trackHandler := delivery.NewTrackHandler(trackMockUsecase, nil, nil)
	e := echo.New()
	trackHandler.Configure(e, nil)

	jsonRequest, err := json.Marshal(&Request{
		Title:   track.Title,
		AlbumID: track.AlbumID,
	})
	assert.Equal(t, err, nil)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/tracks", strings.NewReader(string(jsonRequest)))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	resWriter := httptest.NewRecorder()
	ctx := e.NewContext(req, resWriter)

	handler := trackHandler.HandlerCreateTrack()

	err = handler(ctx)
	assert.Equal(t, err, nil)
	assert.Equal(t, http.StatusOK, resWriter.Code)
}

func TestTrackDelivery_HandlerUpdateTrack(t *testing.T) {
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
		UpdateTrack(gomock.Eq(expectedTrack), uint64(0)).
		Return(nil)

	albumHandler := delivery.NewTrackHandler(trackMockUsecase, nil, nil)
	e := echo.New()
	albumHandler.Configure(e, nil)

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

func TestTrackDelivery_HandlerDeleteTrack(t *testing.T) {
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
	albumHandler.Configure(e, nil)

	req := httptest.NewRequest(http.MethodDelete, "/api/v1/tracks/42", strings.NewReader(""))
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

func TestTrackDelivery_HandlerTracksByParams(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	trackMockUsecase := mock_track.NewMockTrackUsecase(ctrl)

	count, from := uint64(1), uint64(0)

	expectedTracks := []*models.Track{
		{
			ID: 0,
		},
		{
			ID: 1,
		},
	}

	trackMockUsecase.
		EXPECT().
		GetByParams(gomock.Eq(count), gomock.Eq(from), uint64(0)).
		Return(expectedTracks, nil)

	albumHandler := delivery.NewTrackHandler(trackMockUsecase, nil, nil)
	e := echo.New()
	albumHandler.Configure(e, nil)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/tracks?count=1&from=0", strings.NewReader(""))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	resWriter := httptest.NewRecorder()
	ctx := e.NewContext(req, resWriter)

	handler := albumHandler.HandlerTracksByParams()

	err := handler(ctx)
	assert.Equal(t, err, nil)
	assert.Equal(t, http.StatusOK, resWriter.Code)
}

func TestTrackDelivery_HandlerTopTracksByParams(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	trackMockUsecase := mock_track.NewMockTrackUsecase(ctrl)

	count, from := uint64(1), uint64(0)

	expectedTracks := []*models.Track{
		{
			ID: 0,
		},
		{
			ID: 1,
		},
	}

	trackMockUsecase.
		EXPECT().
		GetTopByParams(gomock.Eq(count), gomock.Eq(from), uint64(0)).
		Return(expectedTracks, nil)

	albumHandler := delivery.NewTrackHandler(trackMockUsecase, nil, nil)
	e := echo.New()
	albumHandler.Configure(e, nil)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/tracks/top?count=1&from=0", strings.NewReader(""))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	resWriter := httptest.NewRecorder()
	ctx := e.NewContext(req, resWriter)

	handler := albumHandler.HandlerTopTracksByParams()

	err := handler(ctx)
	assert.Equal(t, err, nil)
	assert.Equal(t, http.StatusOK, resWriter.Code)
}

func TestTrackDelivery_HandlerTracksByArtistID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	trackMockUsecase := mock_track.NewMockTrackUsecase(ctrl)

	id := uint64(42)

	expectedTracks := []*models.Track{
		{
			ID: 0,
		},
		{
			ID: 1,
		},
	}

	trackMockUsecase.
		EXPECT().
		GetByArtistId(id, uint64(0)).
		Return(expectedTracks, nil)

	albumHandler := delivery.NewTrackHandler(trackMockUsecase, nil, nil)
	e := echo.New()
	albumHandler.Configure(e, nil)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/artists/42/tracks", strings.NewReader(""))
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

func TestTrackDelivery_HandlerFavouritesByUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	trackMockUsecase := mock_track.NewMockTrackUsecase(ctrl)
	userMockUsecase := mock_user.NewMockUserUsecase(ctrl)
	sessionMockUsecase := mock_session.NewMockSessionUsecase(ctrl)

	id := uint64(1)
	cookieValue := "Some cookie value"

	session := &models.Session{
		ID:     cookieValue,
		UserID: id,
		Name:   ConstSessionName,
	}

	expectedTracks := []*models.Track{
		{
			ID: 0,
		},
		{
			ID: 1,
		},
	}

	trackMockUsecase.
		EXPECT().
		GetFavoritesByUserID(uint64(1)).
		Return(expectedTracks, nil)

	trackHandler := delivery.NewTrackHandler(trackMockUsecase, sessionMockUsecase, userMockUsecase)
	e := echo.New()
	trackHandler.Configure(e, nil)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/favorite/tracks", strings.NewReader(""))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	req.AddCookie(builder.BuildCookie(session))
	resWriter := httptest.NewRecorder()
	ctx := e.NewContext(req, resWriter)
	ctx.Set(ConstAuthedUserParam, uint64(1))

	handler := trackHandler.HandlerFavouritesByUser()

	err := handler(ctx)
	assert.Equal(t, err, nil)
	assert.Equal(t, http.StatusOK, resWriter.Code)
}

func TestTrackDelivery_HandlerAddToUsersFavourites(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	trackMockUsecase := mock_track.NewMockTrackUsecase(ctrl)
	userMockUsecase := mock_user.NewMockUserUsecase(ctrl)
	sessionMockUsecase := mock_session.NewMockSessionUsecase(ctrl)

	id := uint64(1)
	cookieValue := "Some cookie value"

	session := &models.Session{
		ID:     cookieValue,
		UserID: id,
		Name:   ConstSessionName,
	}

	trackMockUsecase.
		EXPECT().
		AddToFavourites(uint64(1), uint64(1)).
		Return(nil)

	trackHandler := delivery.NewTrackHandler(trackMockUsecase, sessionMockUsecase, userMockUsecase)
	e := echo.New()
	trackHandler.Configure(e, nil)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/favorite/tracks", strings.NewReader(""))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	req.AddCookie(builder.BuildCookie(session))
	resWriter := httptest.NewRecorder()
	ctx := e.NewContext(req, resWriter)
	ctx.Set(ConstAuthedUserParam, uint64(1))
	ctx.SetParamNames("id")
	ctx.SetParamValues("1")

	handler := trackHandler.HandlerAddToUsersFavourites()

	err := handler(ctx)
	assert.Equal(t, err, nil)
	assert.Equal(t, http.StatusOK, resWriter.Code)
}

func TestTrackDelivery_HandlerDeleteFromUsersFavourites(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	trackMockUsecase := mock_track.NewMockTrackUsecase(ctrl)
	userMockUsecase := mock_user.NewMockUserUsecase(ctrl)
	sessionMockUsecase := mock_session.NewMockSessionUsecase(ctrl)

	id := uint64(1)
	cookieValue := "Some cookie value"

	session := &models.Session{
		ID:     cookieValue,
		UserID: id,
		Name:   ConstSessionName,
	}

	trackMockUsecase.
		EXPECT().
		DeleteFromFavourites(uint64(1), uint64(1)).
		Return(nil)

	trackHandler := delivery.NewTrackHandler(trackMockUsecase, sessionMockUsecase, userMockUsecase)
	e := echo.New()
	trackHandler.Configure(e, nil)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/favorite/tracks", strings.NewReader(""))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	req.AddCookie(builder.BuildCookie(session))
	resWriter := httptest.NewRecorder()
	ctx := e.NewContext(req, resWriter)
	ctx.Set(ConstAuthedUserParam, uint64(1))
	ctx.SetParamNames("id")
	ctx.SetParamValues("1")

	handler := trackHandler.HandlerDeleteFromUsersFavourites()

	err := handler(ctx)
	assert.Equal(t, err, nil)
	assert.Equal(t, http.StatusOK, resWriter.Code)
}

func TestTrackDelivery_HandlerLikeTrack(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	trackMockUsecase := mock_track.NewMockTrackUsecase(ctrl)
	userMockUsecase := mock_user.NewMockUserUsecase(ctrl)
	sessionMockUsecase := mock_session.NewMockSessionUsecase(ctrl)

	id := uint64(1)
	cookieValue := "Some cookie value"

	session := &models.Session{
		ID:     cookieValue,
		UserID: id,
		Name:   ConstSessionName,
	}

	trackMockUsecase.
		EXPECT().
		LikeTrack(uint64(1), uint64(1)).
		Return(nil)

	trackHandler := delivery.NewTrackHandler(trackMockUsecase, sessionMockUsecase, userMockUsecase)
	e := echo.New()
	trackHandler.Configure(e, nil)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/tracks/1/like", strings.NewReader(""))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	req.AddCookie(builder.BuildCookie(session))
	resWriter := httptest.NewRecorder()
	ctx := e.NewContext(req, resWriter)
	ctx.Set(ConstAuthedUserParam, uint64(1))
	ctx.SetParamNames("id")
	ctx.SetParamValues("1")

	handler := trackHandler.HandlerLikeTrack()

	err := handler(ctx)
	assert.Equal(t, err, nil)
	assert.Equal(t, http.StatusOK, resWriter.Code)
}

func TestTrackDelivery_HandlerDislikeTrack(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	trackMockUsecase := mock_track.NewMockTrackUsecase(ctrl)
	userMockUsecase := mock_user.NewMockUserUsecase(ctrl)
	sessionMockUsecase := mock_session.NewMockSessionUsecase(ctrl)

	id := uint64(1)
	cookieValue := "Some cookie value"

	session := &models.Session{
		ID:     cookieValue,
		UserID: id,
		Name:   ConstSessionName,
	}

	trackMockUsecase.
		EXPECT().
		DislikeTrack(uint64(1), uint64(1)).
		Return(nil)

	trackHandler := delivery.NewTrackHandler(trackMockUsecase, sessionMockUsecase, userMockUsecase)
	e := echo.New()
	trackHandler.Configure(e, nil)

	req := httptest.NewRequest(http.MethodDelete, "/api/v1/tracks/1/like", strings.NewReader(""))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	req.AddCookie(builder.BuildCookie(session))
	resWriter := httptest.NewRecorder()
	ctx := e.NewContext(req, resWriter)
	ctx.Set(ConstAuthedUserParam, uint64(1))
	ctx.SetParamNames("id")
	ctx.SetParamValues("1")

	handler := trackHandler.HandlerDislikeTrack()

	err := handler(ctx)
	assert.Equal(t, err, nil)
	assert.Equal(t, http.StatusOK, resWriter.Code)
}
