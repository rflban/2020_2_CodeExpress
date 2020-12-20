package repository_test

import (
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/go-park-mail-ru/2020_2_CodeExpress/internal/models"
	"github.com/go-park-mail-ru/2020_2_CodeExpress/internal/playlist/repository"
	"github.com/go-playground/assert/v2"
)

func TestPlaylistRepository_Insert(t *testing.T) {
	db, mock, err := sqlmock.New()

	if err != nil {
		t.Fatalf("cant create mock: %s", err)
	}
	defer db.Close()

	repo := repository.NewPlaylistRep(db)

	title := "Some title"
	userID := uint64(0)
	playlist := &models.Playlist{
		Title:  title,
		UserID: userID,
	}

	rows := sqlmock.NewRows([]string{"id"}).AddRow(1)

	mock.
		ExpectQuery(`insert into playlists`).
		WithArgs(userID, title).
		WillReturnRows(rows)

	if err := repo.Insert(playlist); err != nil {
		t.Errorf("unexpected err: %s", err)
		return
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestPlaylistRepository_Update(t *testing.T) {
	db, mock, err := sqlmock.New()

	if err != nil {
		t.Fatalf("cant create mock: %s", err)
	}
	defer db.Close()

	repo := repository.NewPlaylistRep(db)

	id := uint64(0)
	title := "Some title"
	userID := uint64(0)
	poster := "Some poster"
	isPublic:= false

	playlist := &models.Playlist{
		ID:     id,
		Title:  title,
		UserID: userID,
		Poster: poster,
		IsPublic: isPublic,
	}

	rows := sqlmock.NewRows([]string{"id"}).AddRow(1)

	mock.
		ExpectQuery(`update playlists`).
		WithArgs(title, poster, isPublic, id).
		WillReturnRows(rows)

	if err := repo.Update(playlist); err != nil {
		t.Errorf("unexpected err: %s", err)
		return
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestPlaylistRepository_Delete(t *testing.T) {
	db, mock, err := sqlmock.New()

	if err != nil {
		t.Fatalf("cant create mock: %s", err)
	}
	defer db.Close()

	repo := repository.NewPlaylistRep(db)

	id := uint64(0)
	rows := sqlmock.NewRows([]string{"id"}).AddRow(1)

	mock.
		ExpectQuery(`delete from playlists`).
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

func TestPlaylistRepository_SelectByID(t *testing.T) {
	db, mock, err := sqlmock.New()

	if err != nil {
		t.Fatalf("cant create mock: %s", err)
	}
	defer db.Close()

	repo := repository.NewPlaylistRep(db)

	id := uint64(0)
	title := "Some title"
	userID := uint64(0)
	poster := "Some poster"
	isPublic := false

	expectedPlaylist := &models.Playlist{
		ID:     id,
		Title:  title,
		UserID: userID,
		Poster: poster,
		IsPublic: isPublic,
	}

	rows := sqlmock.NewRows([]string{"id", "user_id", "title", "poster", "is_public"})
	rows.AddRow(id, userID, title, poster, isPublic)

	mock.
		ExpectQuery(`select`).
		WithArgs(id).
		WillReturnRows(rows)

	playlist, err := repo.SelectByID(id)
	if err != nil {
		t.Errorf("unexpected err: %s", err)
		return
	}

	assert.Equal(t, playlist, expectedPlaylist)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestPlaylistRepository_SelectByUserID(t *testing.T) {
	db, mock, err := sqlmock.New()

	if err != nil {
		t.Fatalf("cant create mock: %s", err)
	}
	defer db.Close()

	repo := repository.NewPlaylistRep(db)

	id := uint64(0)
	title := "Some title"
	userID := uint64(0)
	poster := "Some poster"
	isPublic := false

	expectedPlaylist := &models.Playlist{
		ID:     id,
		Title:  title,
		UserID: userID,
		Poster: poster,
		IsPublic: isPublic,
	}

	rows := sqlmock.NewRows([]string{"id", "user_id", "title", "poster", "is_public"})
	rows.AddRow(id, userID, title, poster, isPublic)

	mock.
		ExpectQuery(`select`).
		WithArgs(userID).
		WillReturnRows(rows)

	playlist, err := repo.SelectByUserID(userID)
	if err != nil {
		t.Errorf("unexpected err: %s", err)
		return
	}

	assert.Equal(t, playlist[0], expectedPlaylist)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestPlaylistRepository_InsertTrack(t *testing.T) {
	db, mock, err := sqlmock.New()

	if err != nil {
		t.Fatalf("cant create mock: %s", err)
	}
	defer db.Close()

	repo := repository.NewPlaylistRep(db)

	trackID := uint64(1)
	playlistID := uint64(10)

	rows := sqlmock.NewRows([]string{"track_id"}).AddRow(1)

	mock.
		ExpectQuery(`insert into track_playlist`).
		WithArgs(trackID, playlistID).
		WillReturnRows(rows)

	if err := repo.InsertTrack(trackID, playlistID); err != nil {
		t.Errorf("unexpected err: %s", err)
		return
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestPlaylistRepository_DeleteTrack(t *testing.T) {
	db, mock, err := sqlmock.New()

	if err != nil {
		t.Fatalf("cant create mock: %s", err)
	}
	defer db.Close()

	repo := repository.NewPlaylistRep(db)

	trackID := uint64(1)
	playlistID := uint64(10)

	rows := sqlmock.NewRows([]string{"track_id"}).AddRow(1)

	mock.
		ExpectQuery(`delete from track_playlist`).
		WithArgs(trackID, playlistID).
		WillReturnRows(rows)

	if err := repo.DeleteTrack(trackID, playlistID); err != nil {
		t.Errorf("unexpected err: %s", err)
		return
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}
