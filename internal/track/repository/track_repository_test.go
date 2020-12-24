package repository_test

import (
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/go-park-mail-ru/2020_2_CodeExpress/internal/models"
	"github.com/go-park-mail-ru/2020_2_CodeExpress/internal/track/repository"
	"github.com/go-playground/assert/v2"
)

func TestTrackRepository_Insert(t *testing.T) {
	db, mock, err := sqlmock.New()

	if err != nil {
		t.Fatalf("cant create mock: %s", err)
	}
	defer db.Close()

	repo := repository.NewTrackRep(db)

	title := "Some title"
	albumID := uint64(42)
	id := uint64(1)
	index := uint8(1)
	track := &models.Track{
		Title:   title,
		AlbumID: albumID,
	}

	expectedTrack := &models.Track{
		Title:   title,
		AlbumID: albumID,
		ID:      id,
		Index:   index,
	}

	rows := sqlmock.NewRows([]string{"id", "index"}).AddRow(id, index)

	mock.
		ExpectQuery("INSERT INTO tracks").
		WithArgs(albumID, title).
		WillReturnRows(rows)

	if err := repo.Insert(track); err != nil {
		t.Errorf("unexpected err: %s", err)
		return
	}

	assert.Equal(t, track, expectedTrack)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestTrackRepository_Update(t *testing.T) {
	db, mock, err := sqlmock.New()

	if err != nil {
		t.Fatalf("cant create mock: %s", err)
	}
	defer db.Close()

	repo := repository.NewTrackRep(db)

	title := "Some title"
	albumID := uint64(42)
	id := uint64(1)
	index := uint8(1)

	expectedTrack := &models.Track{
		Title:   title,
		AlbumID: albumID,
		ID:      id,
		Index:   index,
	}

	rows := sqlmock.NewRows([]string{"id"}).AddRow(id)

	mock.
		ExpectQuery("UPDATE tracks").
		WithArgs(albumID, title, index, id).
		WillReturnRows(rows)

	if err := repo.Update(expectedTrack); err != nil {
		t.Errorf("unexpected err: %s", err)
		return
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestTrackRepository_UpdateAudio(t *testing.T) {
	db, mock, err := sqlmock.New()

	if err != nil {
		t.Fatalf("cant create mock: %s", err)
	}
	defer db.Close()

	repo := repository.NewTrackRep(db)

	title := "Some title"
	albumID := uint64(42)
	id := uint64(1)
	index := uint8(1)
	audio := "Some audio"
	duration := 200

	expectedTrack := &models.Track{
		Title:    title,
		AlbumID:  albumID,
		ID:       id,
		Index:    index,
		Audio:    audio,
		Duration: duration,
	}

	rows := sqlmock.NewRows([]string{"id"}).AddRow(id)

	mock.
		ExpectQuery("UPDATE tracks").
		WithArgs(audio, duration, id).
		WillReturnRows(rows)

	if err := repo.UpdateAudio(expectedTrack); err != nil {
		t.Errorf("unexpected err: %s", err)
		return
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestTrackRepository_Delete(t *testing.T) {
	db, mock, err := sqlmock.New()

	if err != nil {
		t.Fatalf("cant create mock: %s", err)
	}
	defer db.Close()

	repo := repository.NewTrackRep(db)

	id := uint64(1)
	rows := sqlmock.NewRows([]string{"id"}).AddRow(id)

	mock.
		ExpectQuery("DELETE FROM tracks").
		WithArgs(id).
		WillReturnRows(rows)

	if err := repo.Delete(id); err != nil {
		t.Errorf("unexpected err: %s", err)
		return
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestTrackRepository_DeleteTrackFromUsersTracks(t *testing.T) {
	db, mock, err := sqlmock.New()

	if err != nil {
		t.Fatalf("cant create mock: %s", err)
	}
	defer db.Close()

	repo := repository.NewTrackRep(db)

	trackID := uint64(1)
	userID := uint64(24)
	rows := sqlmock.NewRows([]string{"track_id"}).AddRow(trackID)

	mock.
		ExpectQuery("DELETE FROM user_track").
		WithArgs(userID, trackID).
		WillReturnRows(rows)

	if err := repo.DeleteTrackFromUsersTracks(userID, trackID); err != nil {
		t.Errorf("unexpected err: %s", err)
		return
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestTrackRepository_InsertTrackToUser(t *testing.T) {
	db, mock, err := sqlmock.New()

	if err != nil {
		t.Fatalf("cant create mock: %s", err)
	}
	defer db.Close()

	repo := repository.NewTrackRep(db)

	trackID := uint64(1)
	userID := uint64(24)
	rows := sqlmock.NewRows([]string{"track_id"}).AddRow(trackID)

	mock.
		ExpectQuery("INSERT INTO user_track").
		WithArgs(userID, trackID).
		WillReturnRows(rows)

	if err := repo.InsertTrackToUser(userID, trackID); err != nil {
		t.Errorf("unexpected err: %s", err)
		return
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestTrackRepository_SelectByID(t *testing.T) {
	db, mock, err := sqlmock.New()

	if err != nil {
		t.Fatalf("cant create mock: %s", err)
	}
	defer db.Close()

	repo := repository.NewTrackRep(db)

	title := "Some title"
	albumID := uint64(42)
	id := uint64(1)
	index := uint8(1)
	audio := "Some audio"
	duration := 200
	poster := "Some poster"
	artistName := "Some artist name"
	artistID := uint64(50)

	expectedTrack := &models.Track{
		Title:       title,
		AlbumID:     albumID,
		ID:          id,
		Index:       index,
		Audio:       audio,
		Duration:    duration,
		ArtistID:    artistID,
		Artist:      artistName,
		AlbumPoster: poster,
	}

	rows := sqlmock.NewRows([]string{"id", "album_id", "title", "duration", "index", "audio", "poster", "name",
		"artist_id", "is_favorite", "is_like"})
	rows.AddRow(id, albumID, title, duration, index, audio, poster, artistName, artistID, nil, nil)

	mock.
		ExpectQuery("SELECT").
		WithArgs(id, uint64(0)).
		WillReturnRows(rows)

	track, err := repo.SelectByID(id, 0)

	if err != nil {
		t.Errorf("unexpected err: %s", err)
		return
	}

	assert.Equal(t, track, expectedTrack)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestTrackRepository_SelectByArtistID(t *testing.T) {
	db, mock, err := sqlmock.New()

	if err != nil {
		t.Fatalf("cant create mock: %s", err)
	}
	defer db.Close()

	repo := repository.NewTrackRep(db)

	title := "Some title"
	albumID := uint64(42)
	id := uint64(1)
	index := uint8(1)
	audio := "Some audio"
	duration := 200
	poster := "Some poster"
	artistName := "Some artist name"
	artistID := uint64(50)

	expectedTrack := &models.Track{
		Title:       title,
		AlbumID:     albumID,
		ID:          id,
		Index:       index,
		Audio:       audio,
		Duration:    duration,
		ArtistID:    artistID,
		Artist:      artistName,
		AlbumPoster: poster,
	}

	expectedTracks := []*models.Track{expectedTrack}

	rows := sqlmock.NewRows([]string{"id", "album_id", "title", "duration", "index", "audio", "poster", "artist_id",
		"name", "is_favorite", "is_like"})
	rows.AddRow(id, albumID, title, duration, index, audio, poster, artistID, artistName, nil, nil)

	mock.
		ExpectQuery("SELECT").
		WithArgs(artistID, uint64(0)).
		WillReturnRows(rows)

	tracks, err := repo.SelectByArtistId(artistID, 0)

	if err != nil {
		t.Errorf("unexpected err: %s", err)
		return
	}

	assert.Equal(t, tracks, expectedTracks)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestTrackRepository_SelectByAlbumID(t *testing.T) {
	db, mock, err := sqlmock.New()

	if err != nil {
		t.Fatalf("cant create mock: %s", err)
	}
	defer db.Close()

	repo := repository.NewTrackRep(db)

	title := "Some title"
	albumID := uint64(42)
	id := uint64(1)
	index := uint8(1)
	audio := "Some audio"
	duration := 200
	poster := "Some poster"
	artistName := "Some artist name"
	artistID := uint64(50)

	expectedTrack := &models.Track{
		Title:       title,
		AlbumID:     albumID,
		ID:          id,
		Index:       index,
		Audio:       audio,
		Duration:    duration,
		ArtistID:    artistID,
		Artist:      artistName,
		AlbumPoster: poster,
	}

	expectedTracks := []*models.Track{expectedTrack}

	rows := sqlmock.NewRows([]string{"id", "album_id", "title", "duration", "index", "audio", "poster", "name",
		"artist_id", "is_favorite", "is_like"})
	rows.AddRow(id, albumID, title, duration, index, audio, poster, artistName, artistID, nil, nil)

	mock.
		ExpectQuery("SELECT").
		WithArgs(albumID, uint64(0)).
		WillReturnRows(rows)

	tracks, err := repo.SelectByAlbumID(albumID, 0)

	if err != nil {
		t.Errorf("unexpected err: %s", err)
		return
	}

	assert.Equal(t, tracks, expectedTracks)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestTrackRepository_SelectByParams(t *testing.T) {
	db, mock, err := sqlmock.New()

	if err != nil {
		t.Fatalf("cant create mock: %s", err)
	}
	defer db.Close()

	repo := repository.NewTrackRep(db)

	title := "Some title"
	albumID := uint64(42)
	id := uint64(1)
	index := uint8(1)
	audio := "Some audio"
	duration := 200
	poster := "Some poster"
	artistName := "Some artist name"
	artistID := uint64(50)

	expectedTrack := &models.Track{
		Title:       title,
		AlbumID:     albumID,
		ID:          id,
		Index:       index,
		Audio:       audio,
		Duration:    duration,
		ArtistID:    artistID,
		Artist:      artistName,
		AlbumPoster: poster,
	}

	expectedTracks := []*models.Track{expectedTrack}

	count, from := uint64(1), uint64(0)

	rows := sqlmock.NewRows([]string{"id", "album_id", "title", "duration", "index", "audio", "poster", "artist_id",
		"name", "is_favorite", "is_like"})
	rows.AddRow(id, albumID, title, duration, index, audio, poster, artistID, artistName, nil, nil)

	mock.
		ExpectQuery("SELECT").
		WithArgs(count, from, uint64(0)).
		WillReturnRows(rows)

	tracks, err := repo.SelectByParams(count, from, 0)

	if err != nil {
		t.Errorf("unexpected err: %s", err)
		return
	}

	assert.Equal(t, tracks, expectedTracks)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestTrackRepository_SelectTopByParams(t *testing.T) {
	db, mock, err := sqlmock.New()

	if err != nil {
		t.Fatalf("cant create mock: %s", err)
	}
	defer db.Close()

	repo := repository.NewTrackRep(db)

	title := "Some title"
	albumID := uint64(42)
	id := uint64(1)
	index := uint8(1)
	audio := "Some audio"
	duration := 200
	poster := "Some poster"
	artistName := "Some artist name"
	artistID := uint64(50)

	expectedTrack := &models.Track{
		Title:       title,
		AlbumID:     albumID,
		ID:          id,
		Index:       index,
		Audio:       audio,
		Duration:    duration,
		ArtistID:    artistID,
		Artist:      artistName,
		AlbumPoster: poster,
	}

	expectedTracks := []*models.Track{expectedTrack}

	count, from := uint64(1), uint64(0)

	rows := sqlmock.NewRows([]string{"id", "album_id", "title", "duration", "index", "audio", "poster", "artist_id",
		"name", "is_favorite", "is_like"})
	rows.AddRow(id, albumID, title, duration, index, audio, poster, artistID, artistName, nil, nil)

	mock.
		ExpectQuery("SELECT").
		WithArgs(count, from, uint64(0)).
		WillReturnRows(rows)

	tracks, err := repo.SelectTopByParams(count, from, 0)

	if err != nil {
		t.Errorf("unexpected err: %s", err)
		return
	}

	assert.Equal(t, tracks, expectedTracks)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestTrackRepository_SelectFavoritesByUserID(t *testing.T) {
	db, mock, err := sqlmock.New()

	if err != nil {
		t.Fatalf("cant create mock: %s", err)
	}
	defer db.Close()

	repo := repository.NewTrackRep(db)

	title := "Some title"
	albumID := uint64(42)
	id := uint64(1)
	index := uint8(1)
	audio := "Some audio"
	duration := 200
	poster := "Some poster"
	artistName := "Some artist name"
	artistID := uint64(50)

	expectedTrack := &models.Track{
		Title:       title,
		AlbumID:     albumID,
		ID:          id,
		Index:       index,
		Audio:       audio,
		Duration:    duration,
		ArtistID:    artistID,
		Artist:      artistName,
		AlbumPoster: poster,
		IsFavorite:  true,
	}

	expectedTracks := []*models.Track{expectedTrack}

	userID := uint64(1)

	rows := sqlmock.NewRows([]string{"id", "album_id", "title", "duration", "index", "audio", "poster", "artist_id",
		"name", "is_like"})
	rows.AddRow(id, albumID, title, duration, index, audio, poster, artistID, artistName, nil)

	mock.
		ExpectQuery("SELECT").
		WithArgs(userID).
		WillReturnRows(rows)

	tracks, err := repo.SelectFavoritesByUserID(userID)

	if err != nil {
		t.Errorf("unexpected err: %s", err)
		return
	}

	assert.Equal(t, tracks, expectedTracks)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestTrackRepository_SelectByPlaylistID(t *testing.T) {
	db, mock, err := sqlmock.New()

	if err != nil {
		t.Fatalf("cant create mock: %s", err)
	}
	defer db.Close()

	repo := repository.NewTrackRep(db)

	title := "Some title"
	albumID := uint64(42)
	id := uint64(1)
	index := uint8(1)
	audio := "Some audio"
	duration := 200
	poster := "Some poster"
	artistName := "Some artist name"
	artistID := uint64(50)

	expectedTrack := &models.Track{
		Title:       title,
		AlbumID:     albumID,
		ID:          id,
		Index:       index,
		Audio:       audio,
		Duration:    duration,
		ArtistID:    artistID,
		Artist:      artistName,
		AlbumPoster: poster,
	}

	expectedTracks := []*models.Track{expectedTrack}

	playlistID := uint64(1)

	rows := sqlmock.NewRows([]string{"id", "album_id", "title", "duration", "index", "audio", "poster", "name",
		"artist_id", "is_favorite", "is_like"})
	rows.AddRow(id, albumID, title, duration, index, audio, poster, artistName, artistID, nil, nil)

	mock.
		ExpectQuery("SELECT").
		WithArgs(playlistID, uint64(0)).
		WillReturnRows(rows)

	tracks, err := repo.SelectByPlaylistID(playlistID, 0)

	if err != nil {
		t.Errorf("unexpected err: %s", err)
		return
	}

	assert.Equal(t, tracks, expectedTracks)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}
