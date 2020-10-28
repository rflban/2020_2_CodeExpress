package repository

import (
	"database/sql"

	"github.com/go-park-mail-ru/2020_2_CodeExpress/internal/models"
	"github.com/go-park-mail-ru/2020_2_CodeExpress/internal/user"
)

type UserRep struct {
	dbConn *sql.DB
}

func NewUserRep(conn *sql.DB) user.UserRep {
	return &UserRep{
		dbConn: conn,
	}
}

func (ur *UserRep) Insert(u *models.User) error {
	query := "insert into users values(default, $1, $2, $3) returning id"

	err := ur.dbConn.QueryRow(query, u.Name, u.Email, u.Password).Scan(&u.ID)

	if err != nil {
		return err
	}

	return nil
}

func (ur *UserRep) Update(user *models.User) error {
	return nil
}

func (ur *UserRep) SelectByEmail(email string) (*models.User, error) {
	query := "select id, name, email, avatar from users where email = $1"
	user := &models.User{}

	err := ur.dbConn.QueryRow(query, email).
		Scan(&user.ID, &user.Name, &user.Email, &user.Avatar)

	if err != nil {
		return nil, err
	}

	return user, nil
}

func (ur *UserRep) SelectByName(name string) (*models.User, error) {
	query := "select id, name, email, avatar from users where name = $1"
	user := &models.User{}

	err := ur.dbConn.QueryRow(query, name).
		Scan(&user.ID, &user.Name, &user.Email, &user.Avatar)

	if err != nil {
		return nil, err
	}

	return user, nil
}

func (ur *UserRep) SelectByID(userID uint64) (*models.User, error) {
	query := "select id, name, email, avatar from users where id = $1"
	user := &models.User{}

	err := ur.dbConn.QueryRow(query, userID).
		Scan(&user.ID, &user.Name, &user.Email, &user.Avatar)

	if err != nil {
		return nil, err
	}

	return user, nil
}

func (ur *UserRep) SelectByLoginAndPassword(login string, password string) (*models.User, error) {
	query := "select id, name, email, avatar from users where (name = $1 or email = $1) and password = $2"
	user := &models.User{}

	err := ur.dbConn.QueryRow(query, login, password).
		Scan(&user.ID, &user.Name, &user.Email, &user.Avatar)

	if err != nil {
		return nil, err
	}

	return user, nil
}
