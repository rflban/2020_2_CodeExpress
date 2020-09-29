package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/go-park-mail-ru/2020_2_CodeExpress/models"
	"github.com/go-park-mail-ru/2020_2_CodeExpress/repositories"
)

type SignUpHandler struct {
	SignUpRep  repositories.SignUpRep
	SessionRep repositories.SessionRep
}

func NewSignUpHandler(SignUpRep repositories.SignUpRep, SessionRep repositories.SessionRep) *SignUpHandler {
	return &SignUpHandler{
		SignUpRep:  SignUpRep,
		SessionRep: SessionRep,
	}
}

func (s *SignUpHandler) HandleCreateUser(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	newUser := new(models.NewUser)
	err := json.NewDecoder(r.Body).Decode(newUser)
	if err != nil {
		log.Printf("Error parsing SignUp JSON %s", err)
		w.Write([]byte(`{"error": "parsing_json"}`))
		return
	}

	user, err := s.createUser(newUser)
	if err != nil {
		w.Write([]byte(fmt.Sprintf(`{"error": "%s"}`, err.Error())))
		return
	}

	err = json.NewEncoder(w).Encode(user)
	if err != nil {
		log.Printf("Error marshalling SignUp JSON %s", err)
		w.Write([]byte(`{"error": "marshalling_json"}`))
		return
	}

	userSession := repositories.NewSession()
	err = s.SessionRep.AddSession(userSession)
	if err != nil {
		log.Printf("Error creating session %s", err)
		w.Write([]byte(`{"error": "creating_session"}`))
		return
	}

	userCookie := &http.Cookie{
		Name:     userSession.Name,
		Value:    userSession.ID,
		Expires:  userSession.Expire,
		HttpOnly: true,
	}
	http.SetCookie(w, userCookie)
}

func (s *SignUpHandler) createUser(newUser *models.NewUser) (*models.User, error) {
	err := s.SignUpRep.CheckUserExists(newUser)
	if err != nil {
		return nil, err
	}
	user, err := s.SignUpRep.CreateUser(newUser)
	return user, err
}
