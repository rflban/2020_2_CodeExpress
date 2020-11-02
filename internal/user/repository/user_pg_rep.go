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

func (ur *UserRep) Insert(user *models.User) error {
	const query string = "INSERT INTO users (id, name, email, password) VALUES (default, $1, $2, $3) RETURNING id;"
	err := ur.dbConn.QueryRow(query, user.Name, user.Email, user.Password).Scan(&user.ID)
	if err != nil {
		return err
	}

	return nil
}

func (ur *UserRep) Update(user *models.User) error {
	const query string = "UPDATE users SET name = $1, email = $2, password = $3, avatar = $4 WHERE id = $5;"
	_, err := ur.dbConn.Exec(query, user.Name, user.Email, user.Password, user.Avatar, user.ID)
	if err != nil {
		return err
	}

	return nil
}

func (ur *UserRep) SelectByEmail(email string) (*models.User, error) {
	query := "SELECT id, name, email, avatar FROM users WHERE email = $1;"
	user := &models.User{}
	err := ur.dbConn.QueryRow(query, email).Scan(&user.ID, &user.Name, &user.Email, &user.Avatar)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (ur *UserRep) SelectByName(name string) (*models.User, error) {
	query := "SELECT id, name, email, avatar FROM users WHERE name = $1;"
	user := &models.User{}
	err := ur.dbConn.QueryRow(query, name).Scan(&user.ID, &user.Name, &user.Email, &user.Avatar)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (ur *UserRep) SelectByNameOrEmail(name string, email string, id uint64) (*models.User, error) {
	query := "SELECT id, name, email, avatar FROM users WHERE (name = $1 OR email = $2) AND id != $3;"
	user := &models.User{}
	err := ur.dbConn.QueryRow(query, name, email, id).Scan(&user.ID, &user.Name, &user.Email, &user.Avatar)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (ur *UserRep) SelectByID(userID uint64) (*models.User, error) {
	query := "SELECT id, name, email, avatar FROM users WHERE id = $1;"
	user := &models.User{}
	err := ur.dbConn.QueryRow(query, userID).Scan(&user.ID, &user.Name, &user.Email, &user.Avatar)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (ur *UserRep) SelectByLoginAndPassword(login string, password string) (*models.User, error) {
	query := "SELECT id, name, email, avatar FROM users WHERE (name = $1 OR email = $1) AND password = $2;"
	user := &models.User{}
	err := ur.dbConn.QueryRow(query, login, password).Scan(&user.ID, &user.Name, &user.Email, &user.Avatar)
	if err != nil {
		return nil, err
	}

	return user, nil
}
