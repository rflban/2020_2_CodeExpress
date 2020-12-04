package repository_test

import (
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/go-park-mail-ru/2020_2_CodeExpress/internal/models"
	"github.com/go-park-mail-ru/2020_2_CodeExpress/internal/user/repository"
	"github.com/go-playground/assert/v2"
	"testing"
)

func TestInsert(t *testing.T) {
	db, mock, err := sqlmock.New()

	if err != nil {
		t.Fatalf("cant create mock: %s", err)
	}
	defer db.Close()

	repo := repository.NewUserRep(db)

	name := "some name"
	email := "some email"
	password := "some password"

	rows := sqlmock.NewRows([]string{"id", "name", "email", "password", "avatar"}).AddRow(1, name, email, password, "")

	mock.
		ExpectQuery(`INSERT INTO users`).
		WithArgs(name, email, password).
		WillReturnRows(rows)

	if _, err := repo.Insert(name, email, password); err != nil {
		t.Errorf("unexpected err: %s", err)
		return
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestUpdate(t *testing.T) {
	db, mock, err := sqlmock.New()

	if err != nil {
		t.Fatalf("cant create mock: %s", err)
	}
	defer db.Close()

	repo := repository.NewUserRep(db)

	id := uint64(1)
	name := "some name"
	email := "some email"
	password := "some password"
	avatar := ""

	expectedUser := &models.User{
		ID:      id,
		Name:   name,
		Email: email,
		Password:   password,
	}

	rows := sqlmock.NewRows([]string{"id", "name", "email", "password", "avatar"}).AddRow(1, name, email, password, avatar)

	mock.
		ExpectQuery(`UPDATE users`).
		WithArgs(name, email, password, avatar, id).
		WillReturnRows(rows)

	if err := repo.Update(expectedUser); err != nil {
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

	repo := repository.NewUserRep(db)

	id := uint64(1)
	name := "some name"
	email := "some email"
	password := "some password"
	avatar := ""

	expectedUser := &models.User{
		ID:      id,
		Name:   name,
		Email: email,
		Password:   password,
	}

	rows := sqlmock.NewRows([]string{"id", "name", "email", "password", "avatar"}).AddRow(1, name, email, password, avatar)

	mock.
		ExpectQuery(`SELECT`).
		WithArgs(id).
		WillReturnRows(rows)

	user, err := repo.SelectById(id)
	if err != nil {
		t.Errorf("unexpected err: %s", err)
		return
	}

	assert.Equal(t, user, expectedUser)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestSelectByLogin(t *testing.T) {
	db, mock, err := sqlmock.New()

	if err != nil {
		t.Fatalf("cant create mock: %s", err)
	}
	defer db.Close()

	repo := repository.NewUserRep(db)

	id := uint64(1)
	name := "some name"
	email := "some email"
	password := "some password"
	avatar := ""

	expectedUser := &models.User{
		ID:      id,
		Name:   name,
		Email: email,
		Password:   password,
	}

	rows := sqlmock.NewRows([]string{"id", "name", "email", "password", "avatar"}).AddRow(1, name, email, password, avatar)

	mock.
		ExpectQuery(`SELECT`).
		WithArgs(email).
		WillReturnRows(rows)

	user, err := repo.SelectByLogin(email)
	if err != nil {
		t.Errorf("unexpected err: %s", err)
		return
	}

	assert.Equal(t, user, expectedUser)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestSelectByNameOrEmail(t *testing.T) {
	db, mock, err := sqlmock.New()

	if err != nil {
		t.Fatalf("cant create mock: %s", err)
	}
	defer db.Close()

	repo := repository.NewUserRep(db)

	id := uint64(1)
	name := "some name"
	email := "some email"
	password := "some password"
	avatar := ""

	expectedUser := &models.User{
		ID:      id,
		Name:   name,
		Email: email,
		Password:   password,
		Avatar: avatar,
	}

	expectedUsers := []*models.User{expectedUser}

	rows := sqlmock.NewRows([]string{"id", "name", "email", "password", "avatar"}).AddRow(1, name, email, password, avatar)

	mock.
		ExpectQuery(`SELECT`).
		WithArgs(name, email).
		WillReturnRows(rows)

	users, err := repo.SelectByNameOrEmail(name, email)
	if err != nil {
		t.Errorf("unexpected err: %s", err)
		return
	}

	assert.Equal(t, users, expectedUsers)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestSelectIfAdmin(t *testing.T) {
	db, mock, err := sqlmock.New()

	if err != nil {
		t.Fatalf("cant create mock: %s", err)
	}
	defer db.Close()

	repo := repository.NewUserRep(db)

	id := uint64(1)

	rows := sqlmock.NewRows([]string{"id"}).AddRow(1)

	mock.
		ExpectQuery(`SELECT`).
		WithArgs(id).
		WillReturnRows(rows)

	_, err = repo.SelectIfAdmin(id)
	if err != nil {
		t.Errorf("unexpected err: %s", err)
		return
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}
