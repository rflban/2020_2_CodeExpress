package repository_test

import (
	"github.com/go-park-mail-ru/2020_2_CodeExpress/internal/models"
	"github.com/go-park-mail-ru/2020_2_CodeExpress/internal/session/repository"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
)

func TestInsert(t *testing.T) {
	db, mock, err := sqlmock.New()

	if err != nil {
		t.Fatalf("cant create mock: %s", err)
	}
	defer db.Close()

	repo := repository.NewSessionRep(db)

	id := "some id"

	session := &models.Session{
		ID: id,
		UserID: 1,
		Expire: time.Now(),
	}

	rows := sqlmock.NewRows([]string{"id"}).AddRow(id)

	mock.
		ExpectQuery(`insert into session`).
		WithArgs(session.ID, session.UserID, session.Expire).
		WillReturnRows(rows)

	if err := repo.Insert(session); err != nil {
		t.Errorf("unexpected err: %s", err)
		return
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestSelectById(t *testing.T) {
	db, mock, err := sqlmock.New()

	if err != nil {
		t.Fatalf("cant create mock: %s", err)
	}
	defer db.Close()

	repo := repository.NewSessionRep(db)

	id := "code_express_session_id"
	expectedSession := &models.Session{
		ID: id,
		UserID: 1,
		Expire: time.Now(),
	}

	rows := sqlmock.NewRows([]string{"id", "userID", "expire"}).AddRow(expectedSession.ID, expectedSession.UserID, expectedSession.Expire)

	mock.
		ExpectQuery(`select`).
		WithArgs(expectedSession.ID).
		WillReturnRows(rows)

	_, err = repo.SelectById(id)
	if err != nil {
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

	repo := repository.NewSessionRep(db)

	id := "some id"

	session := &models.Session{
		ID: id,
		UserID: 1,
		Expire: time.Now(),
	}

	rows := sqlmock.NewRows([]string{"expire"}).AddRow(session.Expire)

	mock.
		ExpectQuery(`delete from session`).
		WithArgs(session.ID).
		WillReturnRows(rows)

	if err := repo.Delete(session); err != nil {
		t.Errorf("unexpected err: %s", err)
		return
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}
