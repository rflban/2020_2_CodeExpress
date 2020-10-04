package test

import (
	"testing"

	"github.com/go-park-mail-ru/2020_2_CodeExpress/models"

	"github.com/go-park-mail-ru/2020_2_CodeExpress/repositories"
)

func TestCheckSessionExistsSuccess(t *testing.T) {
	sImpl := repositories.NewSessionRepImpl()

	user := &models.User{
		ID: 0,
	}
	session := repositories.NewSession(user)

	err := sImpl.AddSession(session)
	if err != nil {
		t.Fatalf("TestCheckSessionExists failed, error: %s", err)
	}

	newSession, err := sImpl.GetSessionByValue(session.ID)

	if err != nil {
		t.Fatalf("TestCheckSessionExists failed, error: %s", err)
	}

	if session.UserID != newSession.UserID {
		t.Fatalf("TestCheckSessionExists failed, session doesn't exist")
	}
}

func TestCheckSessionOutdatesSuccess(t *testing.T) {
	sImpl := repositories.NewSessionRepImpl()

	user := &models.User{
		ID: 0,
	}
	session := repositories.NewSession(user)

	err := sImpl.AddSession(session)
	if err != nil {
		t.Fatalf("TestCheckSessionExists failed, error: %s", err)
	}

	err = sImpl.OutdateSession(session)
	if err != nil {
		t.Fatalf("TestCheckSessionExists failed, error: %s", err)
	}

	if !sImpl.CheckSessionOutdated(session) {
		t.Fatalf("TestCheckSessionOutdates failed, session isn't outdated")
	}
}

func TestCheckSessionOutdatesFail(t *testing.T) {
	sImpl := repositories.NewSessionRepImpl()

	user := &models.User{
		ID: 0,
	}
	session := repositories.NewSession(user)

	err := sImpl.AddSession(session)
	if err != nil {
		t.Fatalf("TestCheckSessionExists failed, error: %s", err)
	}

	if sImpl.CheckSessionOutdated(session) {
		t.Fatalf("TestCheckSessionOutdates not failed, session is outdated")
	}
}

func TestProlongSessionSuccess(t *testing.T) {
	sImpl := repositories.NewSessionRepImpl()

	user := &models.User{
		ID: 0,
	}
	session := repositories.NewSession(user)

	err := sImpl.OutdateSession(session)
	if err != nil {
		t.Fatalf("TestProlongSession failed, error: %s", err)
	}

	err = sImpl.ProlongSession(session)
	if err != nil {
		t.Fatalf("TestProlongSession failed, error: %s", err)
	}

	if sImpl.CheckSessionOutdated(session) {
		t.Fatalf("TestProlongSession failed, session isn't prolonged")
	}
}
