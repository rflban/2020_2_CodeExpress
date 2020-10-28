package repository

import (
	"database/sql"

	"github.com/go-park-mail-ru/2020_2_CodeExpress/internal/consts"
	"github.com/go-park-mail-ru/2020_2_CodeExpress/internal/models"
	"github.com/go-park-mail-ru/2020_2_CodeExpress/internal/session"
)

type SessionRep struct {
	dbConn *sql.DB
}

func NewSessionRep(dbConn *sql.DB) session.SessionRep {
	return &SessionRep{
		dbConn: dbConn,
	}
}

func (sr *SessionRep) SelectById(id string) (*models.Session, error) {
	query := "select id, userID, expire from session where id = $1"
	session := &models.Session{
		Name: consts.ConstSessionName,
	}
	err := sr.dbConn.QueryRow(query, id).Scan(&session.ID, &session.UserID, &session.Expire)

	if err != nil {
		return nil, err
	}

	return session, nil
}

func (sr *SessionRep) Insert(session *models.Session) error {
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

func (sr *SessionRep) Delete(session *models.Session) error {
	query := "delete from session where id = $1 returning expire"

	err := sr.dbConn.QueryRow(
		query,
		session.ID).Scan(&session.Expire)

	if err != nil {
		return err
	}

	return nil
}
