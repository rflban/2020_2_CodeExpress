package delivery_test

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/go-park-mail-ru/2020_2_CodeExpress/internal/artist/delivery"

	"github.com/go-playground/assert/v2"

	"github.com/labstack/echo/v4"

	"github.com/go-park-mail-ru/2020_2_CodeExpress/internal/artist/mock_artist"
	"github.com/go-park-mail-ru/2020_2_CodeExpress/internal/models"
	. "github.com/go-park-mail-ru/2020_2_CodeExpress/internal/tools/error_response"
	"github.com/golang/mock/gomock"
)

func TestHandlerCreateArtist(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUsecase := mock_artist.NewMockArtistUsecase(ctrl)
	artist := &models.Artist{
		Name: "Imagine Dragons",
	}

	expectedArtist := &models.Artist{
		ID:   uint64(1),
		Name: "Imagine Dragons",
	}

	mockUsecase.
		EXPECT().
		CreateArtist(gomock.Eq(artist)).
		DoAndReturn(func(artist *models.Artist) error {
			artist.ID = uint64(1)
			return nil
		})

	artistHandler := delivery.NewArtistHandler(mockUsecase)
	e := echo.New()
	artistHandler.Configure(e)

	jsonArtist, err := json.Marshal(artist)
	assert.Equal(t, err, nil)

	jsonExpectedArtist, err := json.Marshal(expectedArtist)
	assert.Equal(t, err, nil)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/artists", strings.NewReader(string(jsonArtist)))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	resWriter := httptest.NewRecorder()
	ctx := e.NewContext(req, resWriter)

	handler := artistHandler.HandlerCreateArtist()

	err = handler(ctx)
	assert.Equal(t, err, nil)
	assert.Equal(t, http.StatusOK, resWriter.Code)

	resBody, err := ioutil.ReadAll(resWriter.Body)
	assert.Equal(t, err, nil)
	clearBody := resBody[:len(resBody)-1]
	assert.Equal(t, clearBody, jsonExpectedArtist)
}

func TestHandlerUpdateArtist(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUsecase := mock_artist.NewMockArtistUsecase(ctrl)
	artist := &models.Artist{
		ID:   uint64(1),
		Name: "Imagine Dragons",
	}

	expectedArtist := &models.Artist{
		ID:   uint64(1),
		Name: "Imagine Dragons",
	}

	mockUsecase.
		EXPECT().
		UpdateArtistName(gomock.Eq(artist)).
		Return(nil)

	artistHandler := delivery.NewArtistHandler(mockUsecase)
	e := echo.New()
	artistHandler.Configure(e)

	jsonArtist, err := json.Marshal(artist)
	assert.Equal(t, err, nil)

	jsonExpectedArtist, err := json.Marshal(expectedArtist)
	assert.Equal(t, err, nil)

	req := httptest.NewRequest(http.MethodPut, "/api/v1/artists/1", strings.NewReader(string(jsonArtist)))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	resWriter := httptest.NewRecorder()
	ctx := e.NewContext(req, resWriter)
	ctx.SetParamNames("id")
	ctx.SetParamValues("1")

	handler := artistHandler.HandlerUpdateArtist()

	err = handler(ctx)
	assert.Equal(t, err, nil)
	assert.Equal(t, http.StatusOK, resWriter.Code)

	resBody, err := ioutil.ReadAll(resWriter.Body)
	assert.Equal(t, err, nil)
	clearBody := resBody[:len(resBody)-1]
	assert.Equal(t, clearBody, jsonExpectedArtist)
}

func TestHandlerDeleteArtist(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUsecase := mock_artist.NewMockArtistUsecase(ctrl)
	artist := &models.Artist{
		ID:   uint64(1),
		Name: "Imagine Dragons",
	}

	mockUsecase.
		EXPECT().
		DeleteArtist(gomock.Eq(artist.ID)).
		Return(nil)

	artistHandler := delivery.NewArtistHandler(mockUsecase)
	e := echo.New()
	artistHandler.Configure(e)

	jsonArtist, err := json.Marshal(artist)
	assert.Equal(t, err, nil)

	req := httptest.NewRequest(http.MethodDelete, "/api/v1/artists/1", strings.NewReader(string(jsonArtist)))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	resWriter := httptest.NewRecorder()
	ctx := e.NewContext(req, resWriter)
	ctx.SetParamNames("id")
	ctx.SetParamValues("1")

	handler := artistHandler.HandlerDeleteArtist()

	err = handler(ctx)
	assert.Equal(t, err, nil)
	assert.Equal(t, http.StatusOK, resWriter.Code)

	resBody, err := ioutil.ReadAll(resWriter.Body)
	assert.Equal(t, err, nil)
	clearBody := resBody[:len(resBody)-1]
	byteResponse, err := json.Marshal(OKResponse)
	assert.Equal(t, err, nil)
	assert.Equal(t, clearBody, byteResponse)
}

func TestHandlerArtistByID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUsecase := mock_artist.NewMockArtistUsecase(ctrl)

	id := uint64(1)

	expectedArtist := &models.Artist{
		ID:   id,
		Name: "Imagine Dragons",
	}

	mockUsecase.
		EXPECT().
		GetByID(gomock.Eq(id)).
		Return(expectedArtist, nil)

	artistHandler := delivery.NewArtistHandler(mockUsecase)
	e := echo.New()
	artistHandler.Configure(e)

	jsonExpectedArtist, err := json.Marshal(expectedArtist)
	assert.Equal(t, err, nil)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/artists/1", strings.NewReader(""))
	resWriter := httptest.NewRecorder()
	ctx := e.NewContext(req, resWriter)
	ctx.SetParamNames("id")
	ctx.SetParamValues("1")

	handler := artistHandler.HandlerArtistByID()

	err = handler(ctx)
	assert.Equal(t, err, nil)
	assert.Equal(t, http.StatusOK, resWriter.Code)

	resBody, err := ioutil.ReadAll(resWriter.Body)
	assert.Equal(t, err, nil)
	clearBody := resBody[:len(resBody)-1]
	assert.Equal(t, clearBody, jsonExpectedArtist)
}

func TestHandlerArtistByParam(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUsecase := mock_artist.NewMockArtistUsecase(ctrl)

	count := uint64(2)
	from := uint64(0)

	expectedArtist1 := &models.Artist{
		ID:   uint64(1),
		Name: "Imagine Dragons",
	}

	expectedArtist2 := &models.Artist{
		ID:   uint64(2),
		Name: "Imagine Dragons2",
	}

	expectedArtists := []*models.Artist{expectedArtist1, expectedArtist2}

	mockUsecase.
		EXPECT().
		GetByParams(gomock.Eq(count), gomock.Eq(from)).
		Return(expectedArtists, nil)

	artistHandler := delivery.NewArtistHandler(mockUsecase)
	e := echo.New()
	artistHandler.Configure(e)

	jsonExpectedArtists, err := json.Marshal(expectedArtists)
	assert.Equal(t, err, nil)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/artists?count=2&from=0", strings.NewReader(""))
	resWriter := httptest.NewRecorder()
	ctx := e.NewContext(req, resWriter)

	handler := artistHandler.HandlerArtistsByParams()

	err = handler(ctx)
	assert.Equal(t, err, nil)
	assert.Equal(t, http.StatusOK, resWriter.Code)

	resBody, err := ioutil.ReadAll(resWriter.Body)
	assert.Equal(t, err, nil)
	clearBody := resBody[:len(resBody)-1]
	assert.Equal(t, clearBody, jsonExpectedArtists)
}
