package delivery

import (
	"encoding/json"
	"github.com/go-park-mail-ru/2020_2_CodeExpress/internal/models"
	"github.com/go-park-mail-ru/2020_2_CodeExpress/internal/search/mock_search"
	"github.com/golang/mock/gomock"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strconv"
	"testing"
)

func TestSearchDelivery_HandlerSearch_Passed(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	searchMockUsecase := mock_search.NewMockSearchUsecase(ctrl)

	query := "thIs"
	var offset, limit uint64 = 0, 10

	expectedSearch := &models.Search{
		Albums: []*models.Album{
			{
				ID:         1,
				Title:      "This is album title!",
				ArtistID:   1,
				ArtistName: "Jean Elan",
				Poster:     "",
			},
		},
		Artists: []*models.Artist{
			{
				ID:          1,
				Name:        "ARTHIST",
				Poster:      "",
				Avatar:      "",
				Description: "",
			},
		},
		Tracks: []*models.Track{
			{
				ID:          2,
				Title:       "Feel This Alive",
				Duration:    0,
				AlbumPoster: "",
				AlbumID:     1,
				Index:       2,
				Audio:       "",
				Artist:      "Jean Elan",
				ArtistID:    1,
			},
			{
				ID:          1,
				Title:       "This is great track!",
				Duration:    0,
				AlbumPoster: "",
				AlbumID:     2,
				Index:       1,
				Audio:       "",
				Artist:      "Cosmo Klein",
				ArtistID:    2,
			},
		},
	}

	searchMockUsecase.
		EXPECT().
		Search(gomock.Eq(query), gomock.Eq(offset), gomock.Eq(limit)).
		Return(expectedSearch, nil)

	jsonExpectedSearch, err := json.Marshal(expectedSearch)
	assert.Nil(t, err)

	searchHandler := NewSearchHandler(searchMockUsecase)
	e := echo.New()
	searchHandler.Configure(e)

	param := make(url.Values)
	param["query"] = []string{query}
	param["offset"] = []string{strconv.FormatUint(offset, 10)}
	param["limit"] = []string{strconv.FormatUint(limit, 10)}

	request := httptest.NewRequest(http.MethodGet, "/api/v1/search?"+param.Encode(), nil)
	responseWriter := httptest.NewRecorder()
	context := e.NewContext(request, responseWriter)

	handler := searchHandler.HandlerSearch()
	err = handler(context)
	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, responseWriter.Code)

	responseBody, err := ioutil.ReadAll(responseWriter.Body)
	assert.Nil(t, err)
	assert.Equal(t, responseBody, jsonExpectedSearch)
}

func TestSearchDelivery_HandlerSearch_Failed(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	searchMockUsecase := mock_search.NewMockSearchUsecase(ctrl)

	query := "  "
	var offset, limit uint64 = 0, 10

	searchHandler := NewSearchHandler(searchMockUsecase)
	e := echo.New()
	searchHandler.Configure(e)

	param := make(url.Values)
	param["query"] = []string{query}
	param["offset"] = []string{strconv.FormatUint(offset, 10)}
	param["limit"] = []string{strconv.FormatUint(limit, 10)}

	request := httptest.NewRequest(http.MethodGet, "/api/v1/search?"+param.Encode(), nil)
	responseWriter := httptest.NewRecorder()
	context := e.NewContext(request, responseWriter)

	handler := searchHandler.HandlerSearch()
	err := handler(context)
	assert.Nil(t, err)
	assert.Equal(t, http.StatusBadRequest, responseWriter.Code)
}
