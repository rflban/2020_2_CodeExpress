package delivery

import (
	"crypto/md5"
	"encoding/hex"
	"io"
	"net/http"
	"os"

	"github.com/go-park-mail-ru/2020_2_CodeExpress/internal/mwares"
	"github.com/go-park-mail-ru/2020_2_CodeExpress/internal/tools/csrf"

	"github.com/go-park-mail-ru/2020_2_CodeExpress/internal/models"
	"github.com/go-park-mail-ru/2020_2_CodeExpress/internal/tools/builder"

	. "github.com/go-park-mail-ru/2020_2_CodeExpress/internal/consts"
	"github.com/go-park-mail-ru/2020_2_CodeExpress/internal/session"
	. "github.com/go-park-mail-ru/2020_2_CodeExpress/internal/tools/error_response"
	. "github.com/go-park-mail-ru/2020_2_CodeExpress/internal/tools/responser"
	"github.com/go-park-mail-ru/2020_2_CodeExpress/internal/tools/validator"
	"github.com/go-park-mail-ru/2020_2_CodeExpress/internal/user"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
)

type UserHandler struct {
	userUsecase    user.UserUsecase
	sessionUsecase session.SessionUsecase
}

func NewUserHandler(userUsecase user.UserUsecase, sessionUsecase session.SessionUsecase) *UserHandler {
	return &UserHandler{
		userUsecase:    userUsecase,
		sessionUsecase: sessionUsecase,
	}
}

func (uh *UserHandler) Configure(e *echo.Echo, mm *mwares.MiddlewareManager) {
	e.POST("/api/v1/user", uh.handlerRegisterUser())
	e.GET("/api/v1/user", uh.handlerCurrentUserInfo(), mm.CheckAuth)
	e.PUT("/api/v1/user/profile", uh.handlerUpdateProfile(), mm.CheckCSRF, mm.CheckAuth)
	e.PUT("/api/v1/user/password", uh.handlerUpdatePassword(), mm.CheckCSRF, mm.CheckAuth)
	e.PUT("/api/v1/user/photo", uh.handlerUpdateAvatar(), mm.CheckCSRF, mm.CheckAuth)
}

func (uh *UserHandler) handlerRegisterUser() echo.HandlerFunc {
	type Request struct {
		Name             string `json:"username" validate:"required"`
		Email            string `json:"email" validate:"required,email"`
		Password         string `json:"password" validate:"required,gte=8"`
		RepeatedPassword string `json:"repeated_password" validate:"required,eqfield=Password"`
	}

	return func(ctx echo.Context) error {
		req := &Request{}
		if err := validator.NewRequestValidator(ctx).Validate(req); err != nil {
			if err.Error != nil {
				logrus.Info(err.Error)
				validator.GetValidationError(err)
			}
			return ctx.JSON(err.StatusCode, err.UserError)
		}

		user, err := uh.userUsecase.Create(req.Name, req.Email, req.Password)
		if err != nil {
			return RespondWithError(err, ctx)
		}

		session := models.NewSession(user.ID)
		if err := uh.sessionUsecase.CreateSession(session); err != nil {
			return RespondWithError(err, ctx)
		}

		token, e := csrf.NewCSRFToken(session)

		if e != nil {
			return RespondWithError(NewErrorResponse(ErrInternal, e), ctx)
		}

		ctx.Response().Header().Set("X-CSRF-TOKEN", token)

		cookie := builder.BuildCookie(session)
		ctx.SetCookie(cookie)

		return ctx.JSON(http.StatusOK, user)
	}
}

func (uh *UserHandler) handlerCurrentUserInfo() echo.HandlerFunc {
	return func(ctx echo.Context) error {
		user_id := ctx.Get(ConstAuthedUserParam)

		user, errResp := uh.userUsecase.GetById(user_id.(uint64))
		if errResp != nil {
			return RespondWithError(errResp, ctx)
		}

		return ctx.JSON(http.StatusOK, user)
	}
}

func (uh *UserHandler) handlerUpdateProfile() echo.HandlerFunc {
	type Request struct {
		Name  string `json:"username" validate:"required"`
		Email string `json:"email" validate:"required,email"`
	}

	return func(ctx echo.Context) error {
		req := &Request{}
		if err := validator.NewRequestValidator(ctx).Validate(req); err != nil {
			if err.Error != nil {
				logrus.Info(err.Error)
				validator.GetValidationError(err)
			}
			return ctx.JSON(err.StatusCode, err.UserError)
		}

		user_id := ctx.Get(ConstAuthedUserParam)

		user, errResp := uh.userUsecase.UpdateProfile(user_id.(uint64), req.Name, req.Email)
		if errResp != nil {
			return RespondWithError(errResp, ctx)
		}

		return ctx.JSON(http.StatusOK, user)
	}
}

func (uh *UserHandler) handlerUpdatePassword() echo.HandlerFunc {
	type Request struct {
		OldPassword      string `json:"old_password" validate:"required"`
		Password         string `json:"password" validate:"required,gte=8"`
		RepeatedPassword string `json:"repeated_password" validate:"required,eqfield=Password"`
	}

	return func(ctx echo.Context) error {
		req := &Request{}
		if err := validator.NewRequestValidator(ctx).Validate(req); err != nil {
			if err.Error != nil {
				logrus.Info(err.Error)
				validator.GetValidationError(err)
			}
			return ctx.JSON(err.StatusCode, err.UserError)
		}

		user_id := ctx.Get(ConstAuthedUserParam)

		errResp := uh.userUsecase.UpdatePassword(user_id.(uint64), req.OldPassword, req.Password)
		if errResp != nil {
			return RespondWithError(errResp, ctx)
		}

		return ctx.JSON(http.StatusOK, OKResponse)
	}
}

func (uh *UserHandler) handlerUpdateAvatar() echo.HandlerFunc {
	return func(ctx echo.Context) error {
		user_id := ctx.Get(ConstAuthedUserParam)

		user, errResp := uh.userUsecase.GetById(user_id.(uint64))
		if errResp != nil {
			return RespondWithError(errResp, ctx)
		}

		formFile, err := ctx.FormFile("avatar")
		if err != nil {
			errResp := NewErrorResponse(ErrNoAvatar, err)
			return RespondWithError(errResp, ctx)
		}

		if formFile.Size > int64(10<<20) { //TODO: magic number
			errResp := NewErrorResponse(ErrNoAvatar, err) //TODO: другая ошибка будет
			return RespondWithError(errResp, ctx)
		}

		source, err := formFile.Open()
		if err != nil {
			errResp := NewErrorResponse(ErrNoAvatar, err)
			return RespondWithError(errResp, ctx)
		}
		defer func() {
			_ = source.Close()
		}()

		fileExtension := "png" //TODO: убрать в другое место; захордкоженное расширение
		randBytes := md5.Sum([]byte(fileExtension + user.Name))
		randString := hex.EncodeToString(randBytes[:])
		fileName := randString + "." + fileExtension
		pathToNewFile := "./avatars/" + fileName
		destination, err := os.OpenFile(pathToNewFile, os.O_WRONLY|os.O_CREATE, os.FileMode(int(0777)))
		if err != nil {
			return RespondWithError(NewErrorResponse(ErrInternal, err), ctx)
		}
		defer func() {
			_ = destination.Close()
		}()

		if _, err := io.Copy(destination, source); err != nil {
			return RespondWithError(NewErrorResponse(ErrInternal, err), ctx)
		}

		user, errResp = uh.userUsecase.UpdateAvatar(user_id.(uint64), pathToNewFile)
		if errResp != nil {
			return RespondWithError(errResp, ctx)
		}

		return ctx.JSON(http.StatusOK, user)
	}
}
