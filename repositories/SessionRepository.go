package repositories

import (
	"database/sql"
	"errors"
	"time"

	"github.com/go-park-mail-ru/2020_2_CodeExpress/consts"

	"github.com/go-park-mail-ru/2020_2_CodeExpress/models"
	uuid "github.com/satori/go.uuid"
)

type SessionRep interface {
	GetSessionByID(sessionID string) (*models.Session, error)
	GetSessionByUserID(userID uint64) (*models.Session, error)
	CheckSessionOutdated(session *models.Session) bool
	ProlongSession(session *models.Session) error
	OutdateSession(session *models.Session) error
	AddSession(session *models.Session) error
}

func NewSession(u *models.User) *models.Session {
	return &models.Session{
		Name:   consts.ConstSessionName,
		ID:     uuid.NewV4().String(),
		UserID: u.ID,
		Expire: time.Now().AddDate(0, 0, consts.ConstDaysSession),
	}
}

func NewSessionRepPGImpl(conn *sql.DB) SessionRep {
	return &SessionPGRep{
		dbConn: conn,
	}
}

type SessionPGRep struct {
	dbConn *sql.DB
}

func (sr *SessionPGRep) GetSessionByID(sessionID string) (*models.Session, error) {
	query := "select id, userID, expire from session where id = $1"
	session := &models.Session{
		Name: consts.ConstSessionName,
	}
	err := sr.dbConn.QueryRow(query, sessionID).Scan(&session.ID, &session.UserID, &session.Expire)

	switch err {
	case sql.ErrNoRows:
		return nil, errors.New("No session with sessionID") //TODO: вынести в константы
	case nil:
		return session, nil
	default:
		return nil, err
	}
}

func (sr *SessionPGRep) CheckSessionOutdated(session *models.Session) bool {
	return time.Until(session.Expire) < 0
}

func (sr *SessionPGRep) ProlongSession(session *models.Session) error {
	query := "update session set expire = $1 where id = $2 returning expire"

	err := sr.dbConn.QueryRow(
		query,
		time.Now().AddDate(0, 0, consts.ConstDaysSession),
		session.ID).
		Scan(&session.Expire)

	if err != nil {
		return err
	}
	return nil
}

func (sr *SessionPGRep) OutdateSession(session *models.Session) error {
	query := "update session set expire = $1 where id = $2 returning expire"

	err := sr.dbConn.QueryRow(
		query,
		time.Now().AddDate(0, 0, -consts.ConstDaysSession),
		session.ID).
		Scan(&session.Expire)

	if err != nil {
		return err
	}
	return nil
}

func (sr *SessionPGRep) AddSession(session *models.Session) error {
	query := "insert into session values($1, $2, $3) returning id"

	err := sr.dbConn.QueryRow(
		query,
		session.ID,
		session.UserID,
		session.Expire).
		Scan(&session.ID)

	if err != nil {
		return err
	}
	return nil
}

func (sr *SessionPGRep) GetSessionByUserID(userID uint64) (*models.Session, error) {
	query := "select id, userID, expire from session where userID = $1"
	session := &models.Session{
		Name: consts.ConstSessionName,
	}

	err := sr.dbConn.QueryRow(query, userID).Scan(&session.ID, &session.UserID, &session.Expire)

	switch err {
	case sql.ErrNoRows:
		return nil, errors.New("No session with userID") //TODO: вынести в константы
	case nil:
		return session, nil
	default:
		return nil, err
	}
}
