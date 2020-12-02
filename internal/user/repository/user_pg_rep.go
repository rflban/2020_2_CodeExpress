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

func (ur *UserRep) Insert(name string, email string, password string) (*models.User, error) {
	user := &models.User{}
	if err := ur.dbConn.QueryRow(
		"INSERT INTO users (id, name, email, password) VALUES (default, $1, $2, $3) RETURNING id, name, email, password, avatar;",
		name,
		email,
		password,
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

func (ur *UserRep) Update(user *models.User) error {
	if err := ur.dbConn.QueryRow(
		"UPDATE users SET name = $1, email = $2, password = $3, avatar = $4 WHERE id = $5 RETURNING id, name, email, password, avatar;",
		user.Name,
		user.Email,
		user.Password,
		user.Avatar,
		user.ID,
	).Scan(
		&user.ID,
		&user.Name,
		&user.Email,
		&user.Password,
		&user.Avatar,
	); err != nil {
		return err
	}
	return nil
}

func (ur *UserRep) SelectById(id uint64) (*models.User, error) {
	user := &models.User{}
	if err := ur.dbConn.QueryRow(
		"SELECT id, name, email, password, avatar FROM users WHERE id = $1;", //TODO: точно ли есть смысл извлекать id, если он известен? везде так
		id,
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

func (ur *UserRep) SelectByLogin(login string) (*models.User, error) {
	user := &models.User{}
	if err := ur.dbConn.QueryRow(
		"SELECT id, name, email, password, avatar FROM users WHERE name = $1 OR email = $1;",
		login,
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

func (ur *UserRep) SelectByNameOrEmail(name string, email string) ([]*models.User, error) {
	rows, err := ur.dbConn.Query(
		"SELECT id, name, email, password, avatar FROM users WHERE name = $1 OR email = $2;",
		name,
		email,
	)
	//defer func() {
	//	_ = rows.Close()
	//}()
	if err != nil {
		return nil, err
	}
	var users []*models.User
	for rows.Next() {
		user := &models.User{}
		if err := rows.Scan(
			&user.ID,
			&user.Name,
			&user.Email,
			&user.Password,
			&user.Avatar,
		); err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	return users, nil
}

func (ur *UserRep) SelectIfAdmin(userID uint64) (bool, error) {
	query := "SELECT is_admin FROM users WHERE id = $1"

	isAdmin := false
	err := ur.dbConn.QueryRow(query, userID).Scan(&isAdmin)

	if err != nil {
		return isAdmin, err
	}

	return isAdmin, nil
}
