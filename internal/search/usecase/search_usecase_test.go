package usecase

import (
	"github.com/go-park-mail-ru/2020_2_CodeExpress/internal/models"
	"github.com/go-park-mail-ru/2020_2_CodeExpress/internal/search/mock_search"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSearchUsecase_SearchAlbums(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock_search.NewMockSearchRep(ctrl)
	mockUsecase := NewSearchUsecase(mockRepo)

	query := "aLbUM"
	expected := []*models.Album{{
		ID:         1,
		Title:      "1 album",
		ArtistID:   1,
		ArtistName: "Jean Elan",
		Poster:     "",
	}, {
		ID:         2,
		Title:      "2 AlbuM",
		ArtistID:   2,
		ArtistName: "Cosmo Klein",
		Poster:     "",
	}}

	mockRepo.
		EXPECT().
		SelectAlbums(gomock.Eq(query), gomock.Eq(uint64(0)), gomock.Eq(uint64(2))).
		Return(expected, nil)

	result, err := mockUsecase.SearchAlbums(query, 0, 2)
	assert.Nil(t, err)
	assert.Equal(t, expected, result)

	mockRepo.
		EXPECT().
		SelectAlbums(gomock.Eq(query), gomock.Eq(uint64(0)), gomock.Eq(uint64(1))).
		Return(expected[:1], nil)

	result, err = mockUsecase.SearchAlbums(query, 0, 1)
	assert.Nil(t, err)
	assert.Equal(t, expected[:1], result)

	mockRepo.
		EXPECT().
		SelectAlbums(gomock.Eq(query), gomock.Eq(uint64(1)), gomock.Eq(uint64(2))).
		Return(expected[1:], nil)

	result, err = mockUsecase.SearchAlbums(query, 1, 2)
	assert.Nil(t, err)
	assert.Equal(t, expected[1:], result)

	query = "abracadabra"
	mockRepo.
		EXPECT().
		SelectAlbums(gomock.Eq(query), gomock.Eq(uint64(0)), gomock.Eq(uint64(2))).
		Return([]*models.Album{}, nil)

	result, err = mockUsecase.SearchAlbums(query, 0, 2)
	assert.Nil(t, err)
	assert.Equal(t, []*models.Album{}, result)
}

func TestSearchUsecase_SearchArtists(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock_search.NewMockSearchRep(ctrl)
	mockUsecase := NewSearchUsecase(mockRepo)

	query := "E"
	expected := []*models.Artist{{
		ID:          2,
		Name:        "Cosmo Klein",
		Poster:      "",
		Avatar:      "",
		Description: "",
	}, {
		ID:          1,
		Name:        "Jean Elan",
		Poster:      "",
		Avatar:      "",
		Description: "",
	}}

	mockRepo.
		EXPECT().
		SelectArtists(gomock.Eq(query), gomock.Eq(uint64(0)), gomock.Eq(uint64(2))).
		Return(expected, nil)

	result, err := mockUsecase.SearchArtists(query, 0, 2)
	assert.Nil(t, err)
	assert.Equal(t, expected, result)

	mockRepo.
		EXPECT().
		SelectArtists(gomock.Eq(query), gomock.Eq(uint64(0)), gomock.Eq(uint64(1))).
		Return(expected[:1], nil)

	result, err = mockUsecase.SearchArtists(query, 0, 1)
	assert.Nil(t, err)
	assert.Equal(t, expected[:1], result)

	mockRepo.
		EXPECT().
		SelectArtists(gomock.Eq(query), gomock.Eq(uint64(1)), gomock.Eq(uint64(2))).
		Return(expected[1:], nil)

	result, err = mockUsecase.SearchArtists(query, 1, 2)
	assert.Nil(t, err)
	assert.Equal(t, expected[1:], result)

	query = "abracadabra"
	mockRepo.
		EXPECT().
		SelectArtists(gomock.Eq(query), gomock.Eq(uint64(0)), gomock.Eq(uint64(2))).
		Return([]*models.Artist{}, nil)

	result, err = mockUsecase.SearchArtists(query, 0, 2)
	assert.Nil(t, err)
	assert.Equal(t, []*models.Artist{}, result)
}

func TestSearchUsecase_SearchTracks(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock_search.NewMockSearchRep(ctrl)
	mockUsecase := NewSearchUsecase(mockRepo)

	query := "raC"
	expected := []*models.Track{{
		ID:          1,
		Title:       "Track 1",
		Duration:    0,
		AlbumPoster: "",
		AlbumID:     1,
		Index:       1,
		Audio:       "",
		Artist:      "Jean Elan",
		ArtistID:    1,
	}, {
		ID:          1,
		Title:       "Track 2",
		Duration:    0,
		AlbumPoster: "",
		AlbumID:     2,
		Index:       1,
		Audio:       "",
		Artist:      "Cosmo Klein",
		ArtistID:    2,
	}}

	mockRepo.
		EXPECT().
		SelectTracks(gomock.Eq(query), gomock.Eq(uint64(0)), gomock.Eq(uint64(2))).
		Return(expected, nil)

	result, err := mockUsecase.SearchTracks(query, 0, 2)
	assert.Nil(t, err)
	assert.Equal(t, expected, result)

	mockRepo.
		EXPECT().
		SelectTracks(gomock.Eq(query), gomock.Eq(uint64(0)), gomock.Eq(uint64(1))).
		Return(expected[:1], nil)

	result, err = mockUsecase.SearchTracks(query, 0, 1)
	assert.Nil(t, err)
	assert.Equal(t, expected[:1], result)

	mockRepo.
		EXPECT().
		SelectTracks(gomock.Eq(query), gomock.Eq(uint64(1)), gomock.Eq(uint64(2))).
		Return(expected[1:], nil)

	result, err = mockUsecase.SearchTracks(query, 1, 2)
	assert.Nil(t, err)
	assert.Equal(t, expected[1:], result)

	query = "abracadabra"
	mockRepo.
		EXPECT().
		SelectTracks(gomock.Eq(query), gomock.Eq(uint64(0)), gomock.Eq(uint64(2))).
		Return([]*models.Track{}, nil)

	result, err = mockUsecase.SearchTracks(query, 0, 2)
	assert.Nil(t, err)
	assert.Equal(t, []*models.Track{}, result)
}
