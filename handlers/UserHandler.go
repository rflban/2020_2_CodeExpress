package handlers

import (
	"strings"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"errors"

	"github.com/go-park-mail-ru/2020_2_CodeExpress/models"
	"github.com/go-park-mail-ru/2020_2_CodeExpress/repositories"
)

type SignUpHandler struct {
	UserRep  repositories.UserRep
	SessionRep repositories.SessionRep
}

func NewSignUpHandler(UserRep repositories.UserRep, SessionRep repositories.SessionRep) *SignUpHandler {
	return &SignUpHandler{
		UserRep:  UserRep,
		SessionRep: SessionRep,
	}
}

func (s *SignUpHandler) HandleCreateUser(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	newUser := new(models.NewUser)
	err := json.NewDecoder(r.Body).Decode(newUser)
	if err != nil {
		log.Printf("Error parsing SignUp JSON %s", err)
		http.Error(w, `{"error": "internal server error"}`, http.StatusBadRequest)
		return
	}

	user, err := s.createUser(newUser)
	if err != nil {
		http.Error(w, fmt.Sprintf(`{"error": "%s"}`, err.Error()), http.StatusForbidden)
		return
	}

	err = json.NewEncoder(w).Encode(user)
	if err != nil {
		log.Printf("Error marshalling SignUp JSON %s", err)
		http.Error(w, `{"error": "internal server error"}`, http.StatusInternalServerError)
		return
	}

	userSession := repositories.NewSession(user)
	err = s.SessionRep.AddSession(userSession)
	if err != nil {
		log.Printf("Error while creating session %s", err)
		http.Error(w, `{"error": "internal server error"}`, http.StatusInternalServerError)
		return
	}

	userCookie := &http.Cookie{
		Name:     userSession.Name,
		Value:    userSession.ID,
		Expires:  userSession.Expire,
		HttpOnly: true,
	}
	http.SetCookie(w, userCookie)
	w.WriteHeader(http.StatusOK)
}

func (s *SignUpHandler) createUser(newUser *models.NewUser) (*models.User, error) {
	if strings.Compare(newUser.Password, newUser.RepeatedPassword) != 0{
		return nil, errors.New("Passwords do not match")
	}

	user := &models.User{
		Name: newUser.Name,
		Email: newUser.Email,
		Password: newUser.Password,
	}

	err := s.UserRep.CheckUserExists(user)
	if err != nil {
		return nil, err
	}
	err = s.UserRep.CreateUser(user)
	return user, err
}
