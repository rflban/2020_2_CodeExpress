package handlers

import (
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
	"os"

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
		Path:     "/",
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
		Path:     "/",
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
		Path:     "/",
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

	user, err := s.UserRep.GetUserByID(session.UserID)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(&Error{
			Message: NotAuthorized,
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

	user, err = business.UpdateProfile(s.UserRep, profileForm, user)

	if err != nil {
		w.WriteHeader(http.StatusForbidden)
		json.NewEncoder(w).Encode(&Error{
			Message: err.Error(),
		})
	}

	json.NewEncoder(w).Encode(user)
}

func (s *UserHandler) HandleUpdatePassword(w http.ResponseWriter, r *http.Request) {
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

	user, err := s.UserRep.GetUserByID(session.UserID)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(&Error{
			Message: NotAuthorized,
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

func (s *UserHandler) HandleSetAvatar(w http.ResponseWriter, r *http.Request) {
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

	MaxMemorySize := int64(10 << 20)
	r.Body = http.MaxBytesReader(w, r.Body, MaxMemorySize+1024)
	err = r.ParseMultipartForm(MaxMemorySize)
	if err != nil {
		if err.Error() == "http: request body too large" {
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(&Error{
				Message: FileSizeToLarge,
			})
			return
		}
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(&Error{
			Message: FormError,
		})
		return
	}

	imageFile, _, err := r.FormFile("avatar")
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(&Error{
			Message: NoAvatar,
		})
		return
	}
	defer imageFile.Close()

	fileHeader := make([]byte, 512)
	if _, err := imageFile.Read(fileHeader); err != nil {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(&Error{
			Message: FileError,
		})
		return
	}

	extension, isAllowed := business.IsAllowedImageType(fileHeader)
	if !isAllowed {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(&Error{
			Message: FileError,
		})
		return
	}
	imageFile.Seek(0, 0)

	user, _ := s.UserRep.GetUserBySession(session)

	fileName := business.GetFileName(user, extension)
	pathToNewFile := "./avatars/" + fileName

	storageFile, err := os.OpenFile(pathToNewFile, os.O_WRONLY|os.O_CREATE, os.FileMode(int(0777)))
	if err != nil {
		log.Println("error in creating image file: ", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer storageFile.Close()

	if _, err := io.Copy(storageFile, imageFile); err != nil {
		log.Println("error in copying image file: ", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	user.Avatar = pathToNewFile

	err = json.NewEncoder(w).Encode(user)
	if err != nil {
		log.Println("error marshaling to change avatar")
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func (s *UserHandler) HandleCurrentUser(w http.ResponseWriter, r *http.Request) {
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

	user, err := s.UserRep.GetUserBySession(session)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(&Error{
			Message: NotAuthorized,
		})
		return
	}

	json.NewEncoder(w).Encode(user)
}
