package delivery_test

import (
	"database/sql"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/go-park-mail-ru/2020_2_CodeExpress/internal/track/mock_track"

	"github.com/go-park-mail-ru/2020_2_CodeExpress/internal/album/delivery"
	"github.com/go-park-mail-ru/2020_2_CodeExpress/internal/album/mock_album"
	"github.com/go-park-mail-ru/2020_2_CodeExpress/internal/artist/mock_artist"
	. "github.com/go-park-mail-ru/2020_2_CodeExpress/internal/tools/error_response"

	. "github.com/go-park-mail-ru/2020_2_CodeExpress/internal/consts"

	"github.com/go-park-mail-ru/2020_2_CodeExpress/internal/models"
	"github.com/go-playground/assert/v2"
	"github.com/golang/mock/gomock"
	"github.com/labstack/echo/v4"
)

func TestAlbumDelivery_HandlerCreateAlbum(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	albumMockUsecase := mock_album.NewMockAlbumUsecase(ctrl)
	artistMockUsecase := mock_artist.NewMockArtistUsecase(ctrl)

	type Request struct {
		Title    string `json:"title" validate:"required"`
		ArtistID uint64 `json:"artist_id" validate:"required"`
	}

	title := "Some title"
	artistID := uint64(3)
	artistName := "Some name"

	request := &Request{
		Title:    title,
		ArtistID: artistID,
	}

	album := &models.Album{
		Title:      title,
		ArtistID:   artistID,
		ArtistName: artistName,
	}

	artist := &models.Artist{
		ID:   artistID,
		Name: artistName,
	}

	expectedAlbum := &models.Album{
		ID:         uint64(1),
		Title:      title,
		ArtistID:   artistID,
		ArtistName: artistName,
	}

	artistMockUsecase.
		EXPECT().
		GetByID(gomock.Eq(artistID)).
		Return(artist, nil)

	albumMockUsecase.
		EXPECT().
		CreateAlbum(gomock.Eq(album)).
		DoAndReturn(func(album *models.Album) error {
			album.ID = uint64(1)
			return nil
		})

	albumHandler := delivery.NewAlbumHandler(albumMockUsecase, artistMockUsecase, nil, nil, nil)
	e := echo.New()
	albumHandler.Configure(e, nil)

	jsonRequest, err := json.Marshal(request)
	assert.Equal(t, err, nil)

	jsonExpectedAlbum, err := json.Marshal(expectedAlbum)
	assert.Equal(t, err, nil)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/albums", strings.NewReader(string(jsonRequest)))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	resWriter := httptest.NewRecorder()
	ctx := e.NewContext(req, resWriter)

	handler := albumHandler.HandlerCreateAlbum()

	err = handler(ctx)
	assert.Equal(t, err, nil)
	assert.Equal(t, http.StatusOK, resWriter.Code)

	resBody, err := ioutil.ReadAll(resWriter.Body)
	assert.Equal(t, err, nil)
	assert.Equal(t, resBody, jsonExpectedAlbum)
}

func TestAlbumDelivery_HandlerCreateAlbum_Failed(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	albumMockUsecase := mock_album.NewMockAlbumUsecase(ctrl)
	artistMockUsecase := mock_artist.NewMockArtistUsecase(ctrl)

	type Request struct {
		Title    string `json:"title" validate:"required"`
		ArtistID uint64 `json:"artist_id" validate:"required"`
	}

	title := "Some title"
	artistID := uint64(3)

	request := &Request{
		Title:    title,
		ArtistID: artistID,
	}

	artistMockUsecase.
		EXPECT().
		GetByID(gomock.Eq(artistID)).
		Return(nil, NewErrorResponse(ErrArtistNotExist, sql.ErrNoRows))

	albumHandler := delivery.NewAlbumHandler(albumMockUsecase, artistMockUsecase, nil, nil, nil)
	e := echo.New()
	albumHandler.Configure(e, nil)

	jsonRequest, err := json.Marshal(request)
	assert.Equal(t, err, nil)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/albums", strings.NewReader(string(jsonRequest)))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	resWriter := httptest.NewRecorder()
	ctx := e.NewContext(req, resWriter)

	handler := albumHandler.HandlerCreateAlbum()

	err = handler(ctx)
	assert.Equal(t, err, nil)
	assert.Equal(t, http.StatusNotFound, resWriter.Code)
}

func TestAlbumDelivery_HandlerUpdateAlbum(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	albumMockUsecase := mock_album.NewMockAlbumUsecase(ctrl)
	artistMockUsecase := mock_artist.NewMockArtistUsecase(ctrl)

	type Request struct {
		Title    string `json:"title" validate:"required"`
		ArtistID uint64 `json:"artist_id" validate:"required"`
	}

	title := "Some title"
	artistID := uint64(3)
	artistName := "Some name"

	request := &Request{
		Title:    title,
		ArtistID: artistID,
	}

	album := &models.Album{
		ID:         uint64(3),
		Title:      title,
		ArtistID:   artistID,
		ArtistName: artistName,
	}

	artist := &models.Artist{
		ID:   artistID,
		Name: artistName,
	}

	artistMockUsecase.
		EXPECT().
		GetByID(gomock.Eq(artistID)).
		Return(artist, nil)

	albumMockUsecase.
		EXPECT().
		UpdateAlbum(gomock.Eq(album)).
		Return(nil)

	albumHandler := delivery.NewAlbumHandler(albumMockUsecase, artistMockUsecase, nil, nil, nil)
	e := echo.New()
	albumHandler.Configure(e, nil)

	jsonRequest, err := json.Marshal(request)
	assert.Equal(t, err, nil)

	req := httptest.NewRequest(http.MethodPut, "/api/v1/albums/3", strings.NewReader(string(jsonRequest)))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	resWriter := httptest.NewRecorder()
	ctx := e.NewContext(req, resWriter)
	ctx.SetParamNames("id")
	ctx.SetParamValues("3")

	handler := albumHandler.HandlerUpdateAlbum()

	err = handler(ctx)
	assert.Equal(t, err, nil)
	assert.Equal(t, http.StatusOK, resWriter.Code)
}

func TestAlbumDelivery_HandlerDeleteAlbum(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	albumMockUsecase := mock_album.NewMockAlbumUsecase(ctrl)

	albumMockUsecase.
		EXPECT().
		DeleteAlbum(gomock.Eq(uint64(3))).
		Return(nil)

	albumHandler := delivery.NewAlbumHandler(albumMockUsecase, nil, nil, nil, nil)
	e := echo.New()
	albumHandler.Configure(e, nil)

	jsonRequest, err := json.Marshal("")
	assert.Equal(t, err, nil)

	req := httptest.NewRequest(http.MethodPut, "/api/v1/albums/3", strings.NewReader(string(jsonRequest)))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	resWriter := httptest.NewRecorder()
	ctx := e.NewContext(req, resWriter)
	ctx.SetParamNames("id")
	ctx.SetParamValues("3")

	handler := albumHandler.HandlerDeleteAlbum()

	err = handler(ctx)
	assert.Equal(t, err, nil)
	assert.Equal(t, http.StatusOK, resWriter.Code)
}

func TestAlbumDelivery_HandlerAlbumsByArtist(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	albumMockUsecase := mock_album.NewMockAlbumUsecase(ctrl)
	artistMockUsecase := mock_artist.NewMockArtistUsecase(ctrl)

	title := "Some title"
	artistID := uint64(3)
	artistName := "Some name"

	album1 := &models.Album{
		ID:         uint64(5),
		Title:      title,
		ArtistID:   artistID,
		ArtistName: artistName,
	}

	album2 := &models.Album{
		ID:         uint64(34),
		Title:      title,
		ArtistID:   artistID,
		ArtistName: artistName,
	}

	albums := []*models.Album{album1, album2}

	artist := &models.Artist{
		ID:   artistID,
		Name: artistName,
	}

	artistMockUsecase.
		EXPECT().
		GetByID(gomock.Eq(artistID)).
		Return(artist, nil)

	albumMockUsecase.
		EXPECT().
		GetByArtistID(gomock.Eq(artistID)).
		Return(albums, nil)

	albumHandler := delivery.NewAlbumHandler(albumMockUsecase, artistMockUsecase, nil, nil, nil)
	e := echo.New()
	albumHandler.Configure(e, nil)

	jsonRequest, err := json.Marshal("")
	assert.Equal(t, err, nil)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/artists/3/albums", strings.NewReader(string(jsonRequest)))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	resWriter := httptest.NewRecorder()
	ctx := e.NewContext(req, resWriter)
	ctx.SetParamNames("id")
	ctx.SetParamValues("3")

	handler := albumHandler.HandlerAlbumsByArtist()

	err = handler(ctx)
	assert.Equal(t, err, nil)
	assert.Equal(t, http.StatusOK, resWriter.Code)
}

func TestAlbumDelivery_HandlerAlbumsByParams(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	albumMockUsecase := mock_album.NewMockAlbumUsecase(ctrl)

	title := "Some title"
	artistID := uint64(3)
	artistName := "Some name"

	album1 := &models.Album{
		ID:         uint64(5),
		Title:      title,
		ArtistID:   artistID,
		ArtistName: artistName,
	}

	album2 := &models.Album{
		ID:         uint64(34),
		Title:      title,
		ArtistID:   artistID,
		ArtistName: artistName,
	}

	albums := []*models.Album{album1, album2}

	albumMockUsecase.
		EXPECT().
		GetByParams(gomock.Eq(uint64(2)), gomock.Eq(uint64(0))).
		Return(albums, nil)

	albumHandler := delivery.NewAlbumHandler(albumMockUsecase, nil, nil, nil, nil)
	e := echo.New()
	albumHandler.Configure(e, nil)

	jsonRequest, err := json.Marshal("")
	assert.Equal(t, err, nil)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/albums?count=2&from=0", strings.NewReader(string(jsonRequest)))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	resWriter := httptest.NewRecorder()
	ctx := e.NewContext(req, resWriter)

	handler := albumHandler.HandlerAlbumsByParams()

	err = handler(ctx)
	assert.Equal(t, err, nil)
	assert.Equal(t, http.StatusOK, resWriter.Code)
}

func TestAlbumDelivery_HandlerAlbumTracks(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	albumMockUsecase := mock_album.NewMockAlbumUsecase(ctrl)
	trackMockUsecase := mock_track.NewMockTrackUsecase(ctrl)

	title := "Some title"
	artistID := uint64(3)
	artistName := "Some name"
	albumID := uint64(5)

	album := &models.Album{
		ID:         albumID,
		Title:      title,
		ArtistID:   artistID,
		ArtistName: artistName,
	}

	track := &models.Track{
		AlbumID:  albumID,
		ArtistID: artistID,
		Artist:   artistName,
	}

	tracks := []*models.Track{track}

	albumMockUsecase.
		EXPECT().
		GetByID(gomock.Eq(albumID)).
		Return(album, nil)

	trackMockUsecase.
		EXPECT().
		GetByAlbumID(gomock.Eq(albumID), gomock.Eq(uint64(0))).
		Return(tracks, nil)

	albumHandler := delivery.NewAlbumHandler(albumMockUsecase, nil, trackMockUsecase, nil, nil)
	e := echo.New()
	albumHandler.Configure(e, nil)

	jsonRequest, err := json.Marshal("")
	assert.Equal(t, err, nil)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/albums/5", strings.NewReader(string(jsonRequest)))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	resWriter := httptest.NewRecorder()
	ctx := e.NewContext(req, resWriter)
	ctx.SetParamNames("id")
	ctx.SetParamValues("5")

	handler := albumHandler.HandlerAlbumTracks()

	err = handler(ctx)
	assert.Equal(t, err, nil)
	assert.Equal(t, http.StatusOK, resWriter.Code)
}
