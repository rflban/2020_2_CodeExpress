package repository_test

import (
	"testing"

	"github.com/go-park-mail-ru/2020_2_CodeExpress/internal/album/repository"
	"github.com/go-playground/assert/v2"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/go-park-mail-ru/2020_2_CodeExpress/internal/models"
)

func TestAlbumRepository_Insert(t *testing.T) {
	db, mock, err := sqlmock.New()

	if err != nil {
		t.Fatalf("cant create mock: %s", err)
	}
	defer db.Close()

	repo := repository.NewAlbumRep(db)

	title := "Some title"
	artistId := uint64(0)
	album := &models.Album{
		Title:    title,
		ArtistID: artistId,
	}

	rows := sqlmock.NewRows([]string{"id"}).AddRow(1)

	mock.
		ExpectQuery(`insert into albums`).
		WithArgs(artistId, title).
		WillReturnRows(rows)

	if err := repo.Insert(album); err != nil {
		t.Errorf("unexpected err: %s", err)
		return
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestAlbumRepository_Update(t *testing.T) {
	db, mock, err := sqlmock.New()

	if err != nil {
		t.Fatalf("cant create mock: %s", err)
	}
	defer db.Close()

	repo := repository.NewAlbumRep(db)

	id := uint64(0)
	title := "Some title"
	artistId := uint64(0)
	album := &models.Album{
		ID:       id,
		Title:    title,
		ArtistID: artistId,
	}

	rows := sqlmock.NewRows([]string{"id"}).AddRow(1)

	mock.
		ExpectQuery(`update albums`).
		WithArgs(title, artistId, id).
		WillReturnRows(rows)

	if err := repo.Update(album); err != nil {
		t.Errorf("unexpected err: %s", err)
		return
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestAlbumRepository_UpdatePoster(t *testing.T) {
	db, mock, err := sqlmock.New()

	if err != nil {
		t.Fatalf("cant create mock: %s", err)
	}
	defer db.Close()

	repo := repository.NewAlbumRep(db)

	id := uint64(0)
	poster := "Some poster"
	album := &models.Album{
		ID:     id,
		Poster: poster,
	}

	rows := sqlmock.NewRows([]string{"id"}).AddRow(1)

	mock.
		ExpectQuery(`update albums`).
		WithArgs(poster, id).
		WillReturnRows(rows)

	if err := repo.UpdatePoster(album); err != nil {
		t.Errorf("unexpected err: %s", err)
		return
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestAlbumRepository_Delete(t *testing.T) {
	db, mock, err := sqlmock.New()

	if err != nil {
		t.Fatalf("cant create mock: %s", err)
	}
	defer db.Close()

	repo := repository.NewAlbumRep(db)

	id := uint64(0)
	rows := sqlmock.NewRows([]string{"id"}).AddRow(1)

	mock.
		ExpectQuery(`delete from albums`).
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

func TestAlbumRepository_SelectByID(t *testing.T) {
	db, mock, err := sqlmock.New()

	if err != nil {
		t.Fatalf("cant create mock: %s", err)
	}
	defer db.Close()

	repo := repository.NewAlbumRep(db)

	id := uint64(0)
	title := "Some title"
	artistID := uint64(0)
	artistName := "Some artist name"
	poster := "Some poster"

	expectedAlbum := &models.Album{
		ID:         id,
		Title:      title,
		ArtistID:   artistID,
		ArtistName: artistName,
		Poster:     poster,
	}

	rows := sqlmock.NewRows([]string{"id", "artist_id", "title", "poster", "name"})
	rows.AddRow(id, artistID, title, poster, artistName)

	mock.
		ExpectQuery(`select`).
		WithArgs(id).
		WillReturnRows(rows)

	album, err := repo.SelectByID(id)
	if err != nil {
		t.Errorf("unexpected err: %s", err)
		return
	}

	assert.Equal(t, album, expectedAlbum)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestAlbumRepository_SelectByArtistID(t *testing.T) {
	db, mock, err := sqlmock.New()

	if err != nil {
		t.Fatalf("cant create mock: %s", err)
	}
	defer db.Close()

	repo := repository.NewAlbumRep(db)

	id := uint64(0)
	title := "Some title"
	artistID := uint64(0)
	artistName := "Some artist name"
	poster := "Some poster"

	expectedAlbum1 := &models.Album{
		ID:         id,
		Title:      title,
		ArtistID:   artistID,
		ArtistName: artistName,
		Poster:     poster,
	}

	expectedAlbum2 := &models.Album{
		ID:         id + 1,
		Title:      title,
		ArtistID:   artistID,
		ArtistName: artistName,
		Poster:     poster,
	}

	expectedAlbums := []*models.Album{expectedAlbum1, expectedAlbum2}

	rows := sqlmock.NewRows([]string{"id", "artist_id", "title", "poster", "name"})
	rows.AddRow(id, artistID, title, poster, artistName)
	rows.AddRow(id+1, artistID, title, poster, artistName)

	mock.
		ExpectQuery(`select`).
		WithArgs(artistID).
		WillReturnRows(rows)

	albums, err := repo.SelectByArtistID(id)
	if err != nil {
		t.Errorf("unexpected err: %s", err)
		return
	}

	assert.Equal(t, albums, expectedAlbums)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestAlbumRepository_SelectByParam(t *testing.T) {
	db, mock, err := sqlmock.New()

	if err != nil {
		t.Fatalf("cant create mock: %s", err)
	}
	defer db.Close()

	repo := repository.NewAlbumRep(db)

	id := uint64(0)
	title := "Some title"
	artistID := uint64(0)
	artistName := "Some artist name"
	poster := "Some poster"

	expectedAlbum1 := &models.Album{
		ID:         id,
		Title:      title,
		ArtistID:   artistID,
		ArtistName: artistName,
		Poster:     poster,
	}

	expectedAlbum2 := &models.Album{
		ID:         id + 1,
		Title:      title,
		ArtistID:   artistID,
		ArtistName: artistName,
		Poster:     poster,
	}

	count, from := uint64(2), uint64(0)

	expectedAlbums := []*models.Album{expectedAlbum1, expectedAlbum2}

	rows := sqlmock.NewRows([]string{"id", "artist_id", "title", "poster", "name"})
	rows.AddRow(id, artistID, title, poster, artistName)
	rows.AddRow(id+1, artistID, title, poster, artistName)

	mock.
		ExpectQuery(`select`).
		WithArgs(count, from).
		WillReturnRows(rows)

	albums, err := repo.SelectByParam(count, from)
	if err != nil {
		t.Errorf("unexpected err: %s", err)
		return
	}

	assert.Equal(t, albums, expectedAlbums)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}
