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

func (s *UserHandler) HandleUpdateProfile(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	cookie, err := r.Cookie("code_express_session_id")
	if err == http.ErrNoCookie {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(&Error{
			Message: "not authorized",
		})
		return
	}

	session, err := s.SessionRep.GetSessionByValue(cookie.Value)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(&Error{
			Message: "not authorized",
		})
		return
	}

	user, err := s.UserRep.GetUserByID(session.UserID)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(&Error{
			Message: "not authorized",
		})
		return
	}

	profileForm := new(models.ProfileForm)
	err = json.NewDecoder(r.Body).Decode(profileForm)
	if err != nil {
		log.Printf("Error parsing SignUp JSON %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if profileForm.Email != user.Email {
		newUser := new(models.User)
		newUser.Email = profileForm.Email
		err = s.UserRep.CheckUserExists(newUser)

	}
	if err != nil && profileForm.Username != user.Name {
		newUser := new(models.User)
		newUser.Name = profileForm.Username
		err = s.UserRep.CheckUserExists(newUser)
	}
	if err != nil {
		w.WriteHeader(http.StatusForbidden)
		json.NewEncoder(w).Encode(&Error{
			Message: err.Error(),
		})
		return
	}

	user.Email = profileForm.Email
	user.Name = profileForm.Username

	json.NewEncoder(w).Encode(&User{
		Username: user.Name,
		Email:    user.Email,
		ID:       user.ID,
	})
}

func (s *UserHandler) HandleUpdatePassword(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	cookie, err := r.Cookie("code_express_session_id")
	if err == http.ErrNoCookie {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(&Error{
			Message: "not authorized",
		})
		return
	}

	session, err := s.SessionRep.GetSessionByValue(cookie.Value)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(&Error{
			Message: "not authorized",
		})
		return
	}

	user, err := s.UserRep.GetUserByID(session.UserID)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(&Error{
			Message: "not authorized",
		})
		return
	}

	passwordForm := new(models.PasswordForm)
	err = json.NewDecoder(r.Body).Decode(passwordForm)
	if err != nil {
		log.Printf("Error parsing SignUp JSON %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if passwordForm.Password == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(&Error{
			Message: NoPassword,
		})
		return
	}

	if passwordForm.RepeatedPassword == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(&Error{
			Message: NoRepeatedPassword,
		})
		return
	}

	if len(passwordForm.Password) < 8 {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(&Error{
			Message: PasswordTooShort,
		})
		return
	}

	if passwordForm.Password != passwordForm.RepeatedPassword {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(&Error{
			Message: PasswordsMismatch,
		})
		return
	}

	if passwordForm.Password == user.Password {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(&Error{
			Message: PasswordIsOld,
		})
		return
	}

	user.Password = passwordForm.Password

	json.NewEncoder(w).Encode(&Error{
		Message: NoError,
	})
}
