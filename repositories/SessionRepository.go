package repositories

import (
	"errors"
	"sync"
	"time"

	"github.com/go-park-mail-ru/2020_2_CodeExpress/models"
	uuid "github.com/satori/go.uuid"
)

type SessionRep interface {
	GetSessionByValue(sessionID string) (*models.Session, error)
	CheckSessionOutdated(session *models.Session) bool
	ProlongSession(session *models.Session) error
	OutdateSession(session *models.Session) error
	AddSession(session *models.Session) error
}

func NewSession(u *models.User) *models.Session {
	return &models.Session{
		Name:   "code_express_session_id",
		ID:     uuid.NewV4().String(),
		UserID: u.ID,
		Expire: time.Now().AddDate(0, 0, 1),
	}
}

type SessionRepImpl struct {
	Sessions []*models.Session
	MU       *sync.RWMutex
}

func NewSessionRepImpl() SessionRep {
	return &SessionRepImpl{
		Sessions: make([]*models.Session, 0),
		MU:       &sync.RWMutex{},
	}
}

func (s *SessionRepImpl) GetSessionByValue(sessionID string) (*models.Session, error) {
	s.MU.RLock()
	defer s.MU.RUnlock()

	for _, elemSession := range s.Sessions {
		if elemSession.ID == sessionID {
			return elemSession, nil
		}
	}
	return nil, errors.New("No session with sessionID")
}

func (s *SessionRepImpl) CheckSessionOutdated(session *models.Session) bool {
	return time.Until(session.Expire) < 0
}

func (s *SessionRepImpl) ProlongSession(session *models.Session) error {
	s.MU.Lock()
	defer s.MU.Unlock()

	for idx, elemSession := range s.Sessions {
		if elemSession.ID == session.ID {
			s.Sessions[idx].Expire = time.Now().AddDate(0, 0, 5)
		}
	}
	return nil
}

func (s *SessionRepImpl) OutdateSession(session *models.Session) error {
	s.MU.Lock()
	defer s.MU.Unlock()

	for idx, elemSession := range s.Sessions {
		if elemSession.ID == session.ID {
			s.Sessions[idx].Expire = time.Now().AddDate(0, 0, -1)
		}
	}
	return nil
}

func (s *SessionRepImpl) AddSession(session *models.Session) error {
	s.MU.RLock()
	defer s.MU.RUnlock()

	s.Sessions = append(s.Sessions, session)
	return nil
}
