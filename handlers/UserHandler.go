package handlers

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"github.com/go-park-mail-ru/2020_2_CodeExpress/business"

	"github.com/go-park-mail-ru/2020_2_CodeExpress/models"
	"github.com/go-park-mail-ru/2020_2_CodeExpress/repositories"
)

type UserHandler struct {
	UserRep    repositories.UserRep
	SessionRep repositories.SessionRep
}

func NewUserHandler(UserRep repositories.UserRep, SessionRep repositories.SessionRep) *UserHandler {
	return &UserHandler{
		UserRep:    UserRep,
		SessionRep: SessionRep,
	}
}

func (s *UserHandler) decodeNewUser(w http.ResponseWriter, r *http.Request) (*models.NewUser, error) {
	newUser := new(models.NewUser)

	err := json.NewDecoder(r.Body).Decode(newUser)
	if err != nil {
		log.Printf("Error parsing SignUp JSON %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return nil, errors.New(InternalError)
	}
	if newUser.Email == "" {
		w.WriteHeader(http.StatusBadRequest)
		return nil, errors.New(NoEmail)
	}

	if newUser.Name == "" {
		w.WriteHeader(http.StatusBadRequest)
		return nil, errors.New(NoUsername)
	}

	if newUser.Password == "" {
		w.WriteHeader(http.StatusBadRequest)
		return nil, errors.New(NoPassword)
	}

	if newUser.RepeatedPassword == "" {
		w.WriteHeader(http.StatusBadRequest)
		return nil, errors.New(NoRepeatedPassword)
	}

	if len(newUser.Password) < 8 || len(newUser.RepeatedPassword) < 8 {
		w.WriteHeader(http.StatusBadRequest)
		return nil, errors.New(PasswordTooShort)
	}
	return newUser, nil
}

func (s *UserHandler) decodeLogIn(w http.ResponseWriter, r *http.Request) (*models.LogInForm, error) {
	logInForm := new(models.LogInForm)
	err := json.NewDecoder(r.Body).Decode(logInForm)
	if err != nil {
		log.Printf("Error parsing SignUp JSON %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return nil, errors.New(InternalError)
	}

	if logInForm.Login == "" {
		w.WriteHeader(http.StatusBadRequest)
		return nil, errors.New("no login field")
	}

	if logInForm.Password == "" {
		w.WriteHeader(http.StatusBadRequest)
		return nil, errors.New("no password field")
	}
	return logInForm, nil
}

func (s *UserHandler) HandleLogInUser(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	logInForm, err := s.decodeLogIn(w, r)

	if err != nil {
		json.NewEncoder(w).Encode(&Error{
			Message: err.Error(),
		})
		return
	}

	user, err := s.UserRep.LoginUser(logInForm)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(&Error{
			Message: err.Error(),
		})
		return
	}

	userSession := repositories.NewSession(user)
	err = s.SessionRep.AddSession(userSession)
	if err != nil {
		log.Printf("Error while creating session %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(&Error{
			Message: InternalError,
		})
		return
	}

	userCookie := http.Cookie{
		Name:     userSession.Name,
		Value:    userSession.ID,
		Expires:  userSession.Expire,
		HttpOnly: true,
	}
	http.SetCookie(w, &userCookie)

	err = json.NewEncoder(w).Encode(user)
	if err != nil {
		log.Printf("Error marshalling LogIn JSON %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(&Error{
			Message: InternalError,
		})
		return
	}
}

func (s *UserHandler) HandleLogOutUser(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	cookie, err := r.Cookie("code_express_session_id")
	if err == http.ErrNoCookie {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(&Error{
			Message: NotAuthorized,
		})
		return
	}

	session, err := s.SessionRep.GetSessionByValue(cookie.Value)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(&Error{
			Message: NotAuthorized,
		})
		return
	}

	err = s.SessionRep.OutdateSession(session)
	if err != nil {
		log.Printf("Error outdating session %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(&Error{
			Message: InternalError,
		})
		return
	}

	userCookie := http.Cookie{
		Name:     session.Name,
		Value:    session.ID,
		Expires:  session.Expire,
		HttpOnly: true,
	}
	http.SetCookie(w, &userCookie)

	json.NewEncoder(w).Encode(&Error{
		Message: NoError,
	})
}

func (s *UserHandler) HandleCreateUser(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	newUser, err := s.decodeNewUser(w, r)
	if err != nil {
		json.NewEncoder(w).Encode(&Error{
			Message: err.Error(),
		})
		return
	}

	user, err := business.CreateUser(s.UserRep, newUser)
	if err != nil {
		w.WriteHeader(http.StatusForbidden)
		json.NewEncoder(w).Encode(&Error{
			Message: err.Error(),
		})
		return
	}

	userSession := repositories.NewSession(user)
	err = s.SessionRep.AddSession(userSession)
	if err != nil {
		log.Printf("Error while creating session %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(&Error{
			Message: InternalError,
		})
		return
	}

	userCookie := &http.Cookie{
		Name:     userSession.Name,
		Value:    userSession.ID,
		Expires:  userSession.Expire,
		HttpOnly: true,
	}
	http.SetCookie(w, userCookie)

	json.NewEncoder(w).Encode(&Error{
		Message: NoError,
	})
}
