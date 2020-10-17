package handlers

import (
	"errors"
	"github.com/go-park-mail-ru/2020_2_CodeExpress/business"
	"github.com/go-park-mail-ru/2020_2_CodeExpress/consts"
	"github.com/go-park-mail-ru/2020_2_CodeExpress/models"
	"github.com/go-park-mail-ru/2020_2_CodeExpress/repositories"
	"github.com/labstack/echo/v4"
	"io"
	"log"
	"net/http"
	"os"
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

func (uh *UserHandler) decodeNewUser(c echo.Context) (*models.NewUser, error) {
	newUser := new(models.NewUser)
	if err := c.Bind(newUser); err != nil {
		log.Printf("Error parsing JSON %s", err)
		return nil, errors.New(consts.InternalError)
	}

	if newUser.Email == "" {
		return nil, errors.New(consts.NoEmail)
	}

	if newUser.Name == "" {
		return nil, errors.New(consts.NoUsername)
	}

	if newUser.Password == "" {
		return nil, errors.New(consts.NoPassword)
	}

	if newUser.RepeatedPassword == "" {
		return nil, errors.New(consts.NoRepeatedPassword)
	}

	if len(newUser.Password) < 8 || len(newUser.RepeatedPassword) < 8 {
		return nil, errors.New(consts.PasswordTooShort)
	}

	return newUser, nil
}

func (uh *UserHandler) decodeLogIn(c echo.Context) (*models.LogInForm, error) {
	logInForm := new(models.LogInForm)
	if err := c.Bind(logInForm); err != nil {
		log.Printf("Error parsing JSON %s", err)
		return nil, errors.New(consts.InternalError)
	}

	if logInForm.Login == "" {
		return nil, errors.New(consts.NoUsername)
	}

	if logInForm.Password == "" {
		return nil, errors.New(consts.NoPassword)
	}

	return logInForm, nil
}

func (uh *UserHandler) HandleCreateUser(c echo.Context) error {
	newUser, err := uh.decodeNewUser(c)
	if err != nil {
		return c.JSON(http.StatusBadRequest, &Error{
			Message: err.Error(),
		})
	}

	user, err := business.CreateUser(uh.UserRep, newUser)
	if err != nil {
		return c.JSON(http.StatusForbidden, &Error{
			Message: err.Error(),
		})
	}

	userSession := repositories.NewSession(user)
	err = uh.SessionRep.AddSession(userSession)
	if err != nil {
		log.Printf("Error while creating session %s", err)
		return c.JSON(http.StatusForbidden, &Error{
			Message: consts.InternalError,
		})
	}

	userCookie := &http.Cookie{
		Name:     userSession.Name,
		Value:    userSession.ID,
		Expires:  userSession.Expire,
		HttpOnly: true,
		Path:     "/",
	}
	c.SetCookie(userCookie)

	return c.JSON(http.StatusOK, &Error{
		Message: consts.NoError,
	})
}

func (uh *UserHandler) HandleLogInUser(c echo.Context) error {
	logInForm, err := uh.decodeLogIn(c)
	if err != nil {
		return c.JSON(http.StatusBadRequest, &Error{
			Message: err.Error(),
		})
	}

	user, err := uh.UserRep.LoginUser(logInForm)
	if err != nil {
		return c.JSON(http.StatusNotFound, &Error{
			Message: consts.NotAuthorized,
		})
	}

	userSession := repositories.NewSession(user)
	err = uh.SessionRep.AddSession(userSession)
	if err != nil {
		log.Printf("Error while creating session %s", err)
		return c.JSON(http.StatusInternalServerError, &Error{
			Message: consts.InternalError,
		})
	}

	userCookie := http.Cookie{
		Name:     userSession.Name,
		Value:    userSession.ID,
		Expires:  userSession.Expire,
		HttpOnly: true,
		Path:     "/",
	}
	c.SetCookie(&userCookie)

	return c.JSON(http.StatusOK, user)
}

func (uh *UserHandler) HandleLogOutUser(c echo.Context) error {
	cookie, err := c.Cookie("code_express_session_id") //TODO: доставать пользователя в middleware. не только здесь...
	if err == http.ErrNoCookie {
		return c.JSON(http.StatusNotFound, &Error{
			Message: consts.NotAuthorized,
		})
	}
	if err != nil {
		return c.JSON(http.StatusInternalServerError, &Error{
			Message: consts.InternalError,
		})
	}

	session, err := uh.SessionRep.GetSessionByID(cookie.Value)
	if err != nil {
		return c.JSON(http.StatusNotFound, &Error{
			Message: consts.NotAuthorized,
		})
	}

	err = uh.SessionRep.OutdateSession(session)
	if err != nil {
		log.Printf("Error outdating session %s", err)
		return c.JSON(http.StatusInternalServerError, &Error{
			Message: consts.InternalError,
		})
	}

	userCookie := http.Cookie{
		Name:     session.Name,
		Value:    session.ID,
		Expires:  session.Expire,
		HttpOnly: true,
		Path:     "/",
	}
	c.SetCookie(&userCookie)

	return c.JSON(http.StatusOK, &Error{
		Message: consts.NoError,
	})
}

func (uh *UserHandler) HandleUpdateProfile(c echo.Context) error {
	cookie, err := c.Cookie("code_express_session_id") //TODO: доставать пользователя в middleware
	if err == http.ErrNoCookie {
		return c.JSON(http.StatusNotFound, &Error{
			Message: consts.NotAuthorized,
		})
	}
	if err != nil {
		return c.JSON(http.StatusInternalServerError, &Error{
			Message: consts.InternalError,
		})
	}

	session, err := uh.SessionRep.GetSessionByID(cookie.Value)
	if err != nil {
		return c.JSON(http.StatusNotFound, &Error{
			Message: consts.NotAuthorized,
		})
	}

	user, err := uh.UserRep.GetUserByID(session.UserID)
	if err != nil {
		return c.JSON(http.StatusNotFound, &Error{
			Message: consts.NotAuthorized,
		})
	}

	profileForm := new(models.ProfileForm)
	if err = c.Bind(profileForm); err != nil {
		log.Printf("Error parsing JSON %s", err)
		return c.JSON(http.StatusInternalServerError, &Error{
			Message: consts.InternalError,
		})
	}

	if profileForm.Email == "" {
		return c.JSON(http.StatusBadRequest, &Error{
			Message: consts.NoEmail,
		})
	}

	if profileForm.Username == "" {
		return c.JSON(http.StatusBadRequest, &Error{
			Message: consts.NoUsername,
		})
	}

	user, err = business.UpdateProfile(uh.UserRep, profileForm, user)
	if err != nil {
		return c.JSON(http.StatusBadRequest, &Error{
			Message: err.Error(),
		})
	}

	return c.JSON(http.StatusOK, user)
}

func (uh *UserHandler) HandleUpdatePassword(c echo.Context) error {
	cookie, err := c.Cookie("code_express_session_id")
	if err != nil {
		if err == http.ErrNoCookie {
			return c.JSON(http.StatusNotFound, &Error{
				Message: consts.NotAuthorized,
			})
		}
		return c.JSON(http.StatusInternalServerError, &Error{
			Message: consts.InternalError,
		})
	}

	session, err := uh.SessionRep.GetSessionByID(cookie.Value)
	if err != nil {
		return c.JSON(http.StatusNotFound, &Error{
			Message: consts.NotAuthorized,
		})
	}

	user, err := uh.UserRep.GetUserByID(session.UserID)
	if err != nil {
		return c.JSON(http.StatusNotFound, &Error{
			Message: consts.NotAuthorized,
		})
	}

	passwordForm := new(models.PasswordForm)
	if err = c.Bind(passwordForm); err != nil {
		log.Printf("Error parsing JSON %signUpHandler", err)
		return c.JSON(http.StatusInternalServerError, &Error{
			Message: consts.InternalError,
		})
	}

	if passwordForm.Password == "" { //TODO: Подобные проверки будут частью валидации внутренней...
		return c.JSON(http.StatusBadRequest, &Error{
			Message: consts.NoPassword,
		})
	}

	if passwordForm.RepeatedPassword == "" {
		return c.JSON(http.StatusBadRequest, &Error{
			Message: consts.NoRepeatedPassword,
		})
	}

	if len(passwordForm.Password) < 8 {
		return c.JSON(http.StatusBadRequest, &Error{
			Message: consts.PasswordTooShort,
		})
	}

	if passwordForm.Password != passwordForm.RepeatedPassword {
		return c.JSON(http.StatusBadRequest, &Error{
			Message: consts.PasswordsMismatch,
		})
	}

	if passwordForm.Password == user.Password {
		return c.JSON(http.StatusBadRequest, &Error{
			Message: consts.PasswordIsOld,
		})
	}

	user.Password = passwordForm.Password //TODO: реализовать метод business.UpdatePassword

	return c.JSON(http.StatusOK, &Error{
		Message: consts.NoError,
	})
}

func (uh *UserHandler) HandleUpdateAvatar(c echo.Context) error {
	cookie, err := c.Cookie("code_express_session_id") //TODO: доставать пользователя в middleware
	if err != nil {
		if err == http.ErrNoCookie {
			return c.JSON(http.StatusNotFound, &Error{
				Message: consts.NotAuthorized,
			})
		}
		return c.JSON(http.StatusInternalServerError, &Error{
			Message: consts.InternalError,
		})
	}

	session, err := uh.SessionRep.GetSessionByID(cookie.Value)
	if err != nil {
		return c.JSON(http.StatusNotFound, &Error{
			Message: consts.NotAuthorized,
		})
	}

	user, err := uh.UserRep.GetUserByID(session.UserID)
	if err != nil {
		return c.JSON(http.StatusNotFound, &Error{
			Message: consts.NotAuthorized,
		})
	}

	formFile, err := c.FormFile("avatar")
	if err != nil {
		return c.JSON(http.StatusBadRequest, &Error{
			Message: consts.NoAvatar,
		})
	}

	if formFile.Size > int64(10<<20) { //TODO: magic number
		return c.JSON(http.StatusBadRequest, &Error{
			Message: consts.FileSizeToLarge,
		})
	}

	log.Println(formFile.Header)

	source, err := formFile.Open()
	if err != nil {
		return c.JSON(http.StatusBadRequest, &Error{
			Message: consts.NoAvatar,
		})
	}
	defer func() {
		_ = source.Close()
	}()

	fileName := business.GetFileName(user, "png") //TODO: extension
	pathToNewFile := "./avatars/" + fileName
	destination, err := os.OpenFile(pathToNewFile, os.O_WRONLY|os.O_CREATE, os.FileMode(int(0777)))
	if err != nil {
		log.Println("error in creating image file: ", err)
		return c.JSON(http.StatusInternalServerError, &Error{
			Message: consts.InternalError,
		})
	}
	defer func() {
		_ = destination.Close()
	}()

	if _, err := io.Copy(destination, source); err != nil {
		log.Println("error in copying image file: ", err)
		return c.JSON(http.StatusInternalServerError, &Error{
			Message: consts.InternalError,
		})
	}

	user.Avatar = pathToNewFile //TODO: реализовать метод business.UpdateAvatar
	if err = uh.UserRep.ChangeUser(user); err != nil {
		return c.JSON(http.StatusBadRequest, &Error{
			Message: consts.FileError,
		})
	}

	return c.JSON(http.StatusOK, user)
}

func (uh *UserHandler) HandleCurrentUser(c echo.Context) error {
	cookie, err := c.Cookie("code_express_session_id") //TODO: доставать пользователя в middleware
	if err != nil {
		if err == http.ErrNoCookie {
			return c.JSON(http.StatusNotFound, &Error{
				Message: consts.NotAuthorized,
			})
		}
		return c.JSON(http.StatusInternalServerError, &Error{
			Message: consts.InternalError,
		})
	}

	session, err := uh.SessionRep.GetSessionByID(cookie.Value)
	if err != nil {
		return c.JSON(http.StatusNotFound, &Error{
			Message: consts.NotAuthorized,
		})
	}

	user, err := uh.UserRep.GetUserByID(session.UserID)
	if err != nil {
		return c.JSON(http.StatusNotFound, &Error{
			Message: consts.NotAuthorized,
		})
	}

	return c.JSON(http.StatusOK, user)
}
