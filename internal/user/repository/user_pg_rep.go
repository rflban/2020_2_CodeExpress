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
	if err := ur.dbConn.QueryRow(
		"INSERT INTO users (id, name, email, password) VALUES (default, $1, $2, $3) RETURNING id;",
		user.Name,
		user.Email,
		user.Password,
	).Scan(
		&user.ID,
	); err != nil {
		return err
	}
	return nil
}

func (ur *UserRep) Update(user *models.User) error {
	if _, err := ur.dbConn.Exec(
		"UPDATE users SET name = $1, email = $2, avatar = $3 WHERE id = $4;",
		user.Name,
		user.Email,
		user.Avatar,
		user.ID,
	); err != nil {
		return err
	}
	return nil
}

func (ur *UserRep) UpdatePassword(user *models.User) error {
	if _, err := ur.dbConn.Exec(
		"UPDATE users SET password = $1 WHERE id = $2;",
		user.Password,
		user.ID,
	); err != nil {
		return err
	}
	return nil
}

func (ur *UserRep) SelectByEmail(email string) (*models.User, error) {
	user := &models.User{}
	if err := ur.dbConn.QueryRow(
		"SELECT id, name, email, avatar FROM users WHERE email = $1;",
		email,
	).Scan(
		&user.ID,
		&user.Name,
		&user.Email,
		&user.Avatar,
	); err != nil {
		return nil, err
	}
	return user, nil
}

func (ur *UserRep) SelectByName(name string) (*models.User, error) {
	user := &models.User{}
	if err := ur.dbConn.QueryRow(
		"SELECT id, name, email, avatar FROM users WHERE name = $1;",
		name,
	).Scan(
		&user.ID,
		&user.Name,
		&user.Email,
		&user.Avatar,
	); err != nil {
		return nil, err
	}
	return user, nil
}

func (ur *UserRep) SelectByNameWithPassword(name string) (*models.User, error) {
	user := &models.User{}
	if err := ur.dbConn.QueryRow(
		"SELECT id, name, email, password, avatar FROM users WHERE name = $1;",
		name,
	).Scan(
		&user.ID,
		&user.Name,
		&user.Email,
		&user.Password,
		&user.Avatar,
	); err != nil {
		return nil, err
	}
	return user, nil
}

func (ur *UserRep) SelectByNameOrEmail(name string, email string, id uint64) (*models.User, error) {
	user := &models.User{}
	if err := ur.dbConn.QueryRow(
		"SELECT id, name, email, avatar FROM users WHERE (name = $1 OR email = $2) AND id != $3;",
		name,
		email,
		id,
	).Scan(&user.ID, &user.Name, &user.Email, &user.Avatar); err != nil {
		return nil, err
	}
	return user, nil
}

func (ur *UserRep) SelectByID(userID uint64) (*models.User, error) {
	user := &models.User{}
	if err := ur.dbConn.QueryRow(
		"SELECT id, name, email, avatar FROM users WHERE id = $1;",
		userID,
	).Scan(&user.ID, &user.Name, &user.Email, &user.Avatar); err != nil {
		return nil, err
	}
	return user, nil
}
