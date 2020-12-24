package delivery_test

import (
	"encoding/json"
	"github.com/go-park-mail-ru/2020_2_CodeExpress/internal/session/mock_session"
	"github.com/go-park-mail-ru/2020_2_CodeExpress/internal/track/mock_track"
	"github.com/go-park-mail-ru/2020_2_CodeExpress/internal/user/mock_user"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/go-park-mail-ru/2020_2_CodeExpress/internal/models"
	"github.com/go-park-mail-ru/2020_2_CodeExpress/internal/playlist/delivery"
	"github.com/go-park-mail-ru/2020_2_CodeExpress/internal/playlist/mock_playlist"
	"github.com/go-playground/assert/v2"
	"github.com/golang/mock/gomock"
	"github.com/labstack/echo/v4"

	. "github.com/go-park-mail-ru/2020_2_CodeExpress/internal/consts"
)

func TestPlaylistDelivery_HandlerCreatePlaylist(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	playlistUsecase := mock_playlist.NewMockPlaylistUsecase(ctrl)

	type Request struct {
		Title string `json:"title" validate:"required"`
	}

	title := "Some title"
	userID := uint64(3)

	request := &Request{
		Title: title,
	}

	expectedPlaylist := &models.Playlist{
		ID:     uint64(1),
		UserID: userID,
		Title:  title,
	}

	playlist := &models.Playlist{
		Title:  title,
		UserID: userID,
	}

	playlistUsecase.
		EXPECT().
		CreatePlaylist(playlist).
		DoAndReturn(func(playlist *models.Playlist) error {
			playlist.ID = uint64(1)
			return nil
		})

	playlistHandler := delivery.NewPlaylistHandler(playlistUsecase, nil, nil, nil)
	e := echo.New()
	playlistHandler.Configure(e, nil)

	jsonRequest, err := json.Marshal(request)
	assert.Equal(t, err, nil)

	jsonExpectedPlaylist, err := json.Marshal(expectedPlaylist)
	assert.Equal(t, err, nil)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/albums", strings.NewReader(string(jsonRequest)))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	resWriter := httptest.NewRecorder()
	ctx := e.NewContext(req, resWriter)
	ctx.Set(ConstAuthedUserParam, userID)

	handler := playlistHandler.HandlerCreatePlaylist()

	err = handler(ctx)
	assert.Equal(t, err, nil)
	assert.Equal(t, http.StatusOK, resWriter.Code)

	resBody, err := ioutil.ReadAll(resWriter.Body)
	assert.Equal(t, err, nil)
	assert.Equal(t, resBody, jsonExpectedPlaylist)
}

func TestPlaylistDelivery_HandlerUpdatePlaylist(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	playlistUsecase := mock_playlist.NewMockPlaylistUsecase(ctrl)

	type Request struct {
		Title string `json:"title" validate:"required"`
	}

	title := "Some title"
	userID := uint64(3)

	request := &Request{
		Title: title,
	}

	expectedPlaylist := &models.Playlist{
		ID:     uint64(1),
		UserID: userID,
		Title:  title,
	}

	playlist := &models.Playlist{
		ID:     uint64(1),
		Title:  title,
		UserID: userID,
	}

	playlistUsecase.
		EXPECT().
		GetByID(playlist.ID).
		Return(playlist, nil)

	playlistUsecase.
		EXPECT().
		UpdatePlaylist(playlist).
		Return(nil)

	playlistHandler := delivery.NewPlaylistHandler(playlistUsecase, nil, nil, nil)
	e := echo.New()
	playlistHandler.Configure(e, nil)

	jsonRequest, err := json.Marshal(request)
	assert.Equal(t, err, nil)

	jsonExpectedPlaylist, err := json.Marshal(expectedPlaylist)
	assert.Equal(t, err, nil)

	req := httptest.NewRequest(http.MethodPut, "/api/v1/albums/1", strings.NewReader(string(jsonRequest)))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	resWriter := httptest.NewRecorder()
	ctx := e.NewContext(req, resWriter)
	ctx.SetParamNames("id")
	ctx.SetParamValues("1")
	ctx.Set(ConstAuthedUserParam, userID)

	handler := playlistHandler.HandlerUpdatePlaylist()

	err = handler(ctx)
	assert.Equal(t, err, nil)
	assert.Equal(t, http.StatusOK, resWriter.Code)

	resBody, err := ioutil.ReadAll(resWriter.Body)
	assert.Equal(t, err, nil)
	assert.Equal(t, resBody, jsonExpectedPlaylist)
}

func TestPlaylistDelivery_HandlerDeletePlaylist(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	playlistUsecase := mock_playlist.NewMockPlaylistUsecase(ctrl)

	userID := uint64(3)
	playlistID := uint64(42)

	playlist := &models.Playlist{
		ID:     uint64(42),
		UserID: userID,
	}

	playlistUsecase.
		EXPECT().
		GetByID(playlist.ID).
		Return(playlist, nil)

	playlistUsecase.
		EXPECT().
		DeletePlaylist(playlistID).
		Return(nil)

	playlistHandler := delivery.NewPlaylistHandler(playlistUsecase, nil, nil, nil)
	e := echo.New()
	playlistHandler.Configure(e, nil)

	req := httptest.NewRequest(http.MethodDelete, "/api/v1/albums/42", strings.NewReader(string("")))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	resWriter := httptest.NewRecorder()
	ctx := e.NewContext(req, resWriter)
	ctx.SetParamNames("id")
	ctx.SetParamValues("42")
	ctx.Set(ConstAuthedUserParam, userID)

	handler := playlistHandler.HandlerDeletePlaylist()

	err := handler(ctx)
	assert.Equal(t, err, nil)
	assert.Equal(t, http.StatusOK, resWriter.Code)
}

func TestPlaylistDelivery_HandlerUserPlaylists(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	playlistUsecase := mock_playlist.NewMockPlaylistUsecase(ctrl)

	userID := uint64(3)

	playlists := make([]*models.Playlist, 0)

	playlistUsecase.
		EXPECT().
		GetByUserID(userID).
		Return(playlists, nil)

	playlistHandler := delivery.NewPlaylistHandler(playlistUsecase, nil, nil, nil)
	e := echo.New()
	playlistHandler.Configure(e, nil)

	jsonExpectedPlaylists, err := json.Marshal(playlists)
	assert.Equal(t, err, nil)

	req := httptest.NewRequest(http.MethodDelete, "/api/v1/albums", strings.NewReader(string("")))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	resWriter := httptest.NewRecorder()
	ctx := e.NewContext(req, resWriter)
	ctx.Set(ConstAuthedUserParam, userID)

	handler := playlistHandler.HandlerUserPlaylists()

	err = handler(ctx)
	assert.Equal(t, err, nil)
	assert.Equal(t, http.StatusOK, resWriter.Code)

	resBody, err := ioutil.ReadAll(resWriter.Body)
	assert.Equal(t, err, nil)
	assert.Equal(t, resBody, jsonExpectedPlaylists)
}

func TestPlaylistDelivery_HandlerAddTrack(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	playlistUsecase := mock_playlist.NewMockPlaylistUsecase(ctrl)

	type Request struct {
		TrackID uint64 `json:"track_id"`
	}

	userID := uint64(3)
	playlistID := uint64(42)
	trackID := uint64(64)

	playlist := &models.Playlist{
		ID:     playlistID,
		UserID: userID,
	}

	request := &Request{
		TrackID: trackID,
	}

	jsonRequest, err := json.Marshal(request)
	assert.Equal(t, err, nil)

	playlistUsecase.
		EXPECT().
		GetByID(playlist.ID).
		Return(playlist, nil)

	playlistUsecase.
		EXPECT().
		AddTrack(trackID, playlistID).
		Return(nil)

	playlistHandler := delivery.NewPlaylistHandler(playlistUsecase, nil, nil, nil)
	e := echo.New()
	playlistHandler.Configure(e, nil)

	req := httptest.NewRequest(http.MethodDelete, "/api/v1/albums/42", strings.NewReader(string(jsonRequest)))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	resWriter := httptest.NewRecorder()
	ctx := e.NewContext(req, resWriter)
	ctx.SetParamNames("id")
	ctx.SetParamValues("42")
	ctx.Set(ConstAuthedUserParam, userID)

	handler := playlistHandler.HandlerAddTrackToPlaylist()

	err = handler(ctx)
	assert.Equal(t, err, nil)
	assert.Equal(t, http.StatusOK, resWriter.Code)
}

func TestPlaylistDelivery_HandlerDeleteTrackFromPlaylist(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	playlistUsecase := mock_playlist.NewMockPlaylistUsecase(ctrl)

	userID := uint64(3)
	playlistID := uint64(42)
	trackID := uint64(64)

	playlist := &models.Playlist{
		ID:     playlistID,
		UserID: userID,
	}

	playlistUsecase.
		EXPECT().
		GetByID(playlist.ID).
		Return(playlist, nil)

	playlistUsecase.
		EXPECT().
		DeleteTrack(trackID, playlistID).
		Return(nil)

	playlistHandler := delivery.NewPlaylistHandler(playlistUsecase, nil, nil, nil)
	e := echo.New()
	playlistHandler.Configure(e, nil)

	req := httptest.NewRequest(http.MethodDelete, "/api/v1/playlists/42/tracks/64", strings.NewReader(string("")))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	resWriter := httptest.NewRecorder()
	ctx := e.NewContext(req, resWriter)
	ctx.SetParamNames("id", "track_id")
	ctx.SetParamValues("42", "64")
	ctx.Set(ConstAuthedUserParam, userID)

	handler := playlistHandler.HandlerDeleteTrackFromPlaylist()

	err := handler(ctx)
	assert.Equal(t, err, nil)
	assert.Equal(t, http.StatusOK, resWriter.Code)
}

func TestPlaylistDelivery_HandlerConcretePlaylist(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	playlistUsecase := mock_playlist.NewMockPlaylistUsecase(ctrl)
	trackUsecase := mock_track.NewMockTrackUsecase(ctrl)
	userUsecase := mock_user.NewMockUserUsecase(ctrl)
	sessionUsecase := mock_session.NewMockSessionUsecase(ctrl)

	userID := uint64(3)
	playlistID := uint64(42)
	username := "John"

	user := &models.User{
		ID:   userID,
		Name: username,
	}

	tracks := make([]*models.Track, 0)
	tracks = append(tracks, &models.Track{
		Title: "Some title",
	})

	playlist := &models.Playlist{
		ID:       playlistID,
		UserID:   userID,
		Tracks:   tracks,
		IsPublic: true,
	}

	userUsecase.
		EXPECT().
		GetById(userID).
		Return(user, nil)

	playlistUsecase.
		EXPECT().
		GetByID(playlistID).
		Return(playlist, nil)

	trackUsecase.
		EXPECT().
		GetByPlaylistID(playlistID, gomock.Eq(uint64(0))).
		Return(tracks, nil)

	playlistHandler := delivery.NewPlaylistHandler(playlistUsecase, trackUsecase, userUsecase, sessionUsecase)
	e := echo.New()
	playlistHandler.Configure(e, nil)

	req := httptest.NewRequest(http.MethodDelete, "/api/v1/albums/42", strings.NewReader(string("")))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	resWriter := httptest.NewRecorder()
	ctx := e.NewContext(req, resWriter)
	ctx.SetParamNames("id")
	ctx.SetParamValues("42")
	ctx.Set(ConstAuthedUserParam, userID)

	handler := playlistHandler.HandlerConcretePlaylist()

	err := handler(ctx)
	assert.Equal(t, err, nil)
	assert.Equal(t, http.StatusOK, resWriter.Code)

	type Profile struct {
		Username string `json:"username"`
		Avatar   string `json:"avatar"`
	}
	type Response struct {
		Profile  Profile         `json:"profile"`
		Playlist models.Playlist `json:"playlist"`
	}

	profile := Profile{
		user.Name,
		"",
	}

	expected := Response{
		Profile:  profile,
		Playlist: *playlist,
	}

	jsonExpectedPlaylists, err := json.Marshal(expected)
	assert.Equal(t, err, nil)

	resBody, err := ioutil.ReadAll(resWriter.Body)
	assert.Equal(t, err, nil)
	assert.Equal(t, resBody, jsonExpectedPlaylists)
}

func TestPlaylistDelivery_HandlerUserPublicPlaylists(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	playlistUsecase := mock_playlist.NewMockPlaylistUsecase(ctrl)

	userID := uint64(3)

	playlists := make([]*models.Playlist, 0)

	playlistUsecase.
		EXPECT().
		GetPublicByUserID(userID).
		Return(playlists, nil)

	playlistHandler := delivery.NewPlaylistHandler(playlistUsecase, nil, nil, nil)
	e := echo.New()
	playlistHandler.Configure(e, nil)

	jsonExpectedPlaylists, err := json.Marshal(playlists)
	assert.Equal(t, err, nil)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/user/3/playlists", strings.NewReader(string("")))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	resWriter := httptest.NewRecorder()
	ctx := e.NewContext(req, resWriter)
	ctx.SetParamNames("id")
	ctx.SetParamValues("3")
	ctx.Set(ConstAuthedUserParam, userID)

	handler := playlistHandler.HandlerUserPublicPlaylists()

	err = handler(ctx)
	assert.Equal(t, err, nil)
	assert.Equal(t, http.StatusOK, resWriter.Code)

	resBody, err := ioutil.ReadAll(resWriter.Body)
	assert.Equal(t, err, nil)
	assert.Equal(t, resBody, jsonExpectedPlaylists)
}

func TestPlaylistDelivery_HandlerUpdatePrivacyPlaylist(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	playlistUsecase := mock_playlist.NewMockPlaylistUsecase(ctrl)

	type Request struct {
		IsPublic bool `json:"is_public"`
	}

	title := "Some title"
	userID := uint64(3)

	request := &Request{
		IsPublic: true,
	}

	playlist := &models.Playlist{
		ID:     uint64(1),
		Title:  title,
		UserID: userID,
	}

	playlistUsecase.
		EXPECT().
		GetByID(playlist.ID).
		Return(playlist, nil)

	playlistUsecase.
		EXPECT().
		UpdatePlaylist(playlist).
		Return(nil)

	playlistHandler := delivery.NewPlaylistHandler(playlistUsecase, nil, nil, nil)
	e := echo.New()
	playlistHandler.Configure(e, nil)

	jsonRequest, err := json.Marshal(request)
	assert.Equal(t, err, nil)

	assert.Equal(t, err, nil)

	req := httptest.NewRequest(http.MethodPut, "/api/v1/albums/1", strings.NewReader(string(jsonRequest)))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	resWriter := httptest.NewRecorder()
	ctx := e.NewContext(req, resWriter)
	ctx.SetParamNames("id")
	ctx.SetParamValues("1")
	ctx.Set(ConstAuthedUserParam, userID)

	handler := playlistHandler.HandlerUpdatePlaylist()

	err = handler(ctx)
	assert.Equal(t, err, nil)
	assert.Equal(t, http.StatusOK, resWriter.Code)

	assert.Equal(t, err, nil)
}
