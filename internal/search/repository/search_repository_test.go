package repository

import (
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/go-park-mail-ru/2020_2_CodeExpress/internal/models"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSelectRepository_SelectAlbums(t *testing.T) {
	db, mock, err := sqlmock.New()
	defer func() {
		_ = db.Close()
	}()
	if err != nil {
		t.Fatalf("cant create mock: %s", err)
	}

	repository := NewSearchRep(db)

	expected := []*models.Album{
		{
			ID:         1,
			Title:      "This is album title!",
			ArtistID:   1,
			ArtistName: "Jean Elan",
			Poster:     "",
		},
	}

	rows := sqlmock.NewRows([]string{"albums.id", "albums.title", "albums.artist_id", "artists.name", "albums.poster"})
	for _, album := range expected {
		rows.AddRow(album.ID, album.Title, album.ArtistID, album.ArtistName, album.Poster)
	}

	query := "tHiS"
	var offset, limit uint64 = 0, 10

	mock.
		ExpectQuery("SELECT").
		WithArgs(query, limit, offset).
		WillReturnRows(rows)

	result, err := repository.SelectAlbums(query, offset, limit)
	if err != nil {
		t.Errorf("unexpected err: %s", err)
		return
	}

	assert.Equal(t, result, expected)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestSelectRepository_SelectArtists(t *testing.T) {
	db, mock, err := sqlmock.New()
	defer func() {
		_ = db.Close()
	}()
	if err != nil {
		t.Fatalf("cant create mock: %s", err)
	}

	repository := NewSearchRep(db)

	expected := []*models.Artist{
		{
			ID:          1,
			Name:        "ARTHIST",
			Poster:      "",
			Avatar:      "",
			Description: "",
		},
	}

	rows := sqlmock.NewRows([]string{"artists.id", "artists.name", "artists.description", "artists.poster",
		"artists.avatar"})
	for _, artist := range expected {
		rows.AddRow(artist.ID, artist.Name, artist.Description, artist.Poster, artist.Avatar)
	}

	query := "tHiS"
	var offset, limit uint64 = 0, 10

	mock.
		ExpectQuery("SELECT").
		WithArgs(query, limit, offset).
		WillReturnRows(rows)

	result, err := repository.SelectArtists(query, offset, limit)
	if err != nil {
		t.Errorf("unexpected err: %s", err)
		return
	}

	assert.Equal(t, result, expected)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestSelectRepository_SelectTracks(t *testing.T) {
	db, mock, err := sqlmock.New()
	defer func() {
		_ = db.Close()
	}()
	if err != nil {
		t.Fatalf("cant create mock: %s", err)
	}

	repository := NewSearchRep(db)

	expected := []*models.Track{
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
	}

	rows := sqlmock.NewRows([]string{"tracks.id", "tracks.album_id", "albums.poster", "albums.artist_id",
		"artists.name", "tracks.title", "tracks.duration", "tracks.index", "tracks.audio"})
	for _, track := range expected {
		rows.AddRow(track.ID, track.AlbumID, track.AlbumPoster, track.ArtistID, track.Artist, track.Title,
			track.Duration, track.Index, track.Audio)
	}

	query := "tHiS"
	var offset, limit uint64 = 0, 10

	mock.
		ExpectQuery("SELECT").
		WithArgs(query, limit, offset).
		WillReturnRows(rows)

	result, err := repository.SelectTracks(query, offset, limit)
	if err != nil {
		t.Errorf("unexpected err: %s", err)
		return
	}

	assert.Equal(t, result, expected)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}
