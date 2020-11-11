package repository

import (
	"testing"

	"github.com/go-park-mail-ru/2020_2_CodeExpress/internal/models"
	"github.com/stretchr/testify/assert"

	"github.com/DATA-DOG/go-sqlmock"
)

func TestCreate(t *testing.T) {
	db, mock, err := sqlmock.New()

	if err != nil {
		t.Fatalf("cant create mock: %s", err)
	}
	defer db.Close()

	repo := NewArtistRep(db)

	name := "Johny Depp"
	description := "SMTH"
	testArtist := &models.Artist{
		Name:        name,
		Description: description,
	}

	rows := sqlmock.NewRows([]string{"id"}).AddRow(1)

	mock.
		ExpectQuery(`insert into artists`).
		WithArgs(name, description).
		WillReturnRows(rows)

	if err := repo.Insert(testArtist); err != nil {
		t.Errorf("unexpected err: %s", err)
		return
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestUpdateName(t *testing.T) {
	db, mock, err := sqlmock.New()

	if err != nil {
		t.Fatalf("cant create mock: %s", err)
	}
	defer db.Close()

	repo := NewArtistRep(db)

	name := "Johny Depp"
	testArtist := &models.Artist{
		ID:   1,
		Name: name,
	}

	rows := sqlmock.NewRows([]string{"id"}).AddRow(1)

	mock.
		ExpectQuery(`update artists`).
		WithArgs(testArtist.Name, testArtist.Poster, testArtist.Avatar, testArtist.Description, testArtist.ID).
		WillReturnRows(rows)

	if err := repo.Update(testArtist); err != nil {
		t.Errorf("unexpected err: %s", err)
		return
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestUpdatePoster(t *testing.T) {
	db, mock, err := sqlmock.New()

	if err != nil {
		t.Fatalf("cant create mock: %s", err)
	}
	defer db.Close()

	repo := NewArtistRep(db)

	poster := "some poster"
	testArtist := &models.Artist{
		ID:     1,
		Poster: poster,
	}

	rows := sqlmock.NewRows([]string{"id"}).AddRow(1)

	mock.
		ExpectQuery(`update artists`).
		WithArgs(testArtist.Name, testArtist.Poster, testArtist.Avatar, testArtist.Description, testArtist.ID).
		WillReturnRows(rows)

	if err := repo.Update(testArtist); err != nil {
		t.Errorf("unexpected err: %s", err)
		return
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestDelete(t *testing.T) {
	db, mock, err := sqlmock.New()

	if err != nil {
		t.Fatalf("cant create mock: %s", err)
	}
	defer db.Close()

	repo := NewArtistRep(db)

	id := uint64(1)

	rows := sqlmock.NewRows([]string{"id"}).AddRow(1)

	mock.
		ExpectQuery(`delete from artists`).
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

func TestSelectByID(t *testing.T) {
	db, mock, err := sqlmock.New()

	if err != nil {
		t.Fatalf("cant create mock: %s", err)
	}
	defer db.Close()

	rows := sqlmock.NewRows([]string{"id", "name", "poster", "avatar", "description"}).AddRow(1, "Johny Depp", "", "", "")
	query := "select"

	repo := NewArtistRep(db)

	mock.ExpectQuery(query).WithArgs(uint64(1)).WillReturnRows(rows)

	artist, err := repo.SelectByID(uint64(1))

	assert.NoError(t, err)
	assert.NotNil(t, artist)
}

func TestSelectByName(t *testing.T) {
	db, mock, err := sqlmock.New()

	if err != nil {
		t.Fatalf("cant create mock: %s", err)
	}
	defer db.Close()

	rows := sqlmock.NewRows([]string{"id", "name", "poster", "avatar", "description"}).AddRow(1, "Johny Depp", "", "", "")
	query := "select"

	repo := NewArtistRep(db)

	mock.ExpectQuery(query).WithArgs("Johny Depp").WillReturnRows(rows)

	artist, err := repo.SelectByName("Johny Depp")

	assert.NoError(t, err)
	assert.NotNil(t, artist)
}

func TestSelectByParam(t *testing.T) {
	db, mock, err := sqlmock.New()

	if err != nil {
		t.Fatalf("cant create mock: %s", err)
	}
	defer db.Close()

	rows := sqlmock.NewRows([]string{"id", "name", "poster", "avatar", "description"})

	rows.AddRow(1, "Johny Depp", "", "", "")
	rows.AddRow(2, "Johny Depp2", "", "", "")

	query := "select"

	repo := NewArtistRep(db)

	mock.ExpectQuery(query).WithArgs(2, 0).WillReturnRows(rows)

	artists, err := repo.SelectByParam(2, 0)

	assert.NoError(t, err)
	assert.NotNil(t, artists)
	assert.Len(t, artists, 2)
}
