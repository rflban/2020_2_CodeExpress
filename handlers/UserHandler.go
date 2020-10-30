package handlers

import (
	"io"
	"log"
	"net/http"
	"os"

	"github.com/go-park-mail-ru/2020_2_CodeExpress/business"
	"github.com/go-park-mail-ru/2020_2_CodeExpress/consts"
	"github.com/go-park-mail-ru/2020_2_CodeExpress/models"
	"github.com/go-park-mail-ru/2020_2_CodeExpress/repositories"
	"github.com/labstack/echo/v4"
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
