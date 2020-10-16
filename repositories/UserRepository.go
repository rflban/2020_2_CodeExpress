package repositories

import (
	"database/sql"
	"errors"

	"github.com/go-park-mail-ru/2020_2_CodeExpress/models"
)

type UserRep interface {
	CheckUserExists(u *models.User) (*models.User, error)
	CreateUser(u *models.User) error
	GetUserByID(userID uint64) (*models.User, error)
	LoginUser(loginForm *models.LogInForm) (*models.User, error)
	GetUserBySession(session *models.Session) (*models.User, error)
	ChangeUser(user *models.User) error
}

func NewUserRepPGImpl(conn *sql.DB) UserRep {
	return &UserPGImpl{
		dbConn: conn,
	}
}

type UserPGImpl struct {
	dbConn *sql.DB
}

func (ur *UserPGImpl) CreateUser(u *models.User) error {
	query := "insert into users values(default, $1, $2, $3) returning id"

	err := ur.dbConn.QueryRow(query, u.Name, u.Email, u.Password).Scan(&u.ID)
	if err != nil {
		return err
	}
	return nil
}

func (ur *UserPGImpl) CheckUserExists(u *models.User) (*models.User, error) {
	query := "select id, name, email, avatar from users where name = $1 or email = $2"
	newUser := &models.User{}

	err := ur.dbConn.QueryRow(query, u.Name, u.Email).
		Scan(&newUser.ID, &newUser.Name, &newUser.Email, &newUser.Avatar)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return newUser, nil
}

func (ur *UserPGImpl) GetUserByID(userID uint64) (*models.User, error) {
	query := "select id, name, email, password, avatar from users where id = $1"
	user := &models.User{}

	err := ur.dbConn.QueryRow(query, userID).Scan(
		&user.ID,
		&user.Name,
		&user.Email,
		&user.Password,
		&user.Avatar)

	switch err {
	case sql.ErrNoRows:
		return nil, errors.New("No such user") //TODO: вынести в константы
	case nil:
		return user, nil
	default:
		return nil, err
	}
}

func (ur *UserPGImpl) LoginUser(loginForm *models.LogInForm) (*models.User, error) {
	query := "select id, name, email, password, avatar from users where (name = $1 or email = $1) and password = $2"
	user := &models.User{}

	err := ur.dbConn.QueryRow(query, loginForm.Login, loginForm.Password).Scan(
		&user.ID,
		&user.Name,
		&user.Email,
		&user.Password,
		&user.Avatar)

	switch err {
	case sql.ErrNoRows:
		return nil, errors.New("Incorrect login or password") //TODO: вынести в константы
	case nil:
		return user, nil
	default:
		return nil, err
	}
}

func (ur *UserPGImpl) GetUserBySession(session *models.Session) (*models.User, error) {
	query := "select id, name, email, password, avatar from users where id = $1"
	user := &models.User{}

	err := ur.dbConn.QueryRow(query, session.UserID).Scan(
		&user.ID,
		&user.Name,
		&user.Email,
		&user.Password,
		&user.Avatar)

	switch err {
	case sql.ErrNoRows:
		return nil, errors.New("User not exists") //TODO: вынести в константы
	case nil:
		return user, nil
	default:
		return nil, err
	}
}

func (ur *UserPGImpl) ChangeUser(user *models.User) error {
	query := "update users set name = $1, email = $2, password = $3, avatar = $4 where id = $5" +
		"returning id"

	err := ur.dbConn.QueryRow(query,
		user.Name,
		user.Email,
		user.Password,
		user.Avatar,
		user.ID).Scan(
		&user.ID)

	switch err {
	case sql.ErrNoRows:
		return errors.New("User not exists")
	case nil:
		return nil
	default:
		return err
	}
}
