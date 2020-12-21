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

	query := "A"
	expected := &models.Search{
		Albums: []*models.Album{
			{
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
			},
		},
		Artists: []*models.Artist{
			{
				ID:          2,
				Name:        "Imagine Dragons",
				Poster:      "",
				Avatar:      "",
				Description: "",
			}, {
				ID:          1,
				Name:        "Jean Elan",
				Poster:      "",
				Avatar:      "",
				Description: "",
			},
		},
		Tracks: []*models.Track{
			{
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
			},
		},
	}

	mockRepo.
		EXPECT().
		SelectAlbums(gomock.Eq(query), gomock.Eq(uint64(0)), gomock.Eq(uint64(2))).
		Return(expected.Albums, nil)

	mockRepo.
		EXPECT().
		SelectArtists(gomock.Eq(query), gomock.Eq(uint64(0)), gomock.Eq(uint64(2))).
		Return(expected.Artists, nil)

	mockRepo.
		EXPECT().
		SelectTracks(gomock.Eq(query), gomock.Eq(uint64(0)), gomock.Eq(uint64(2)), gomock.Eq(uint64(0))).
		Return(expected.Tracks, nil)

	result, err := mockUsecase.Search(query, 0, 2, 0)
	assert.Nil(t, err)
	assert.Equal(t, expected, result)

	mockRepo.
		EXPECT().
		SelectAlbums(gomock.Eq(query), gomock.Eq(uint64(0)), gomock.Eq(uint64(1))).
		Return(expected.Albums[:1], nil)

	mockRepo.
		EXPECT().
		SelectArtists(gomock.Eq(query), gomock.Eq(uint64(0)), gomock.Eq(uint64(1))).
		Return(expected.Artists[:1], nil)

	mockRepo.
		EXPECT().
		SelectTracks(gomock.Eq(query), gomock.Eq(uint64(0)), gomock.Eq(uint64(1)), gomock.Eq(uint64(0))).
		Return(expected.Tracks[:1], nil)

	expected2 := &models.Search{
		Albums:  expected.Albums[:1],
		Artists: expected.Artists[:1],
		Tracks:  expected.Tracks[:1],
	}
	result, err = mockUsecase.Search(query, 0, 1, 0)
	assert.Nil(t, err)
	assert.Equal(t, expected2, result)

	mockRepo.
		EXPECT().
		SelectAlbums(gomock.Eq(query), gomock.Eq(uint64(1)), gomock.Eq(uint64(2))).
		Return(expected.Albums[1:], nil)

	mockRepo.
		EXPECT().
		SelectArtists(gomock.Eq(query), gomock.Eq(uint64(1)), gomock.Eq(uint64(2))).
		Return(expected.Artists[1:], nil)

	mockRepo.
		EXPECT().
		SelectTracks(gomock.Eq(query), gomock.Eq(uint64(1)), gomock.Eq(uint64(2)), gomock.Eq(uint64(0))).
		Return(expected.Tracks[1:], nil)

	expected3 := &models.Search{
		Albums:  expected.Albums[1:],
		Artists: expected.Artists[1:],
		Tracks:  expected.Tracks[1:],
	}
	result, err = mockUsecase.Search(query, 1, 2, 0)
	assert.Nil(t, err)
	assert.Equal(t, expected3, result)

	query = "abracadabra"

	mockRepo.
		EXPECT().
		SelectAlbums(gomock.Eq(query), gomock.Eq(uint64(0)), gomock.Eq(uint64(2))).
		Return([]*models.Album{}, nil)

	mockRepo.
		EXPECT().
		SelectArtists(gomock.Eq(query), gomock.Eq(uint64(0)), gomock.Eq(uint64(2))).
		Return([]*models.Artist{}, nil)

	mockRepo.
		EXPECT().
		SelectTracks(gomock.Eq(query), gomock.Eq(uint64(0)), gomock.Eq(uint64(2)), gomock.Eq(uint64(0))).
		Return([]*models.Track{}, nil)

	expected4 := &models.Search{
		Albums:  []*models.Album{},
		Artists: []*models.Artist{},
		Tracks:  []*models.Track{},
	}
	result, err = mockUsecase.Search(query, 0, 2, 0)
	assert.Nil(t, err)
	assert.Equal(t, expected4, result)
}

/*func TestSearchUsecase_SearchArtists(t *testing.T) {
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
}*/
