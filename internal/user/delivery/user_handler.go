package delivery

import (
	"net/http"

	"github.com/go-park-mail-ru/2020_2_CodeExpress/internal/mwares"
	"github.com/go-park-mail-ru/2020_2_CodeExpress/internal/tools/csrf"

	"github.com/go-park-mail-ru/2020_2_CodeExpress/internal/models"
	"github.com/go-park-mail-ru/2020_2_CodeExpress/internal/tools/builder"

	. "github.com/go-park-mail-ru/2020_2_CodeExpress/internal/consts"
	"github.com/go-park-mail-ru/2020_2_CodeExpress/internal/session"
	. "github.com/go-park-mail-ru/2020_2_CodeExpress/internal/tools/error_response"
	"github.com/go-park-mail-ru/2020_2_CodeExpress/internal/tools/photo_uploader"
	. "github.com/go-park-mail-ru/2020_2_CodeExpress/internal/tools/responser"
	"github.com/go-park-mail-ru/2020_2_CodeExpress/internal/tools/validator"
	"github.com/go-park-mail-ru/2020_2_CodeExpress/internal/user"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
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
	e.GET("/api/v1/user", uh.handlerCurrentUserInfo())
	e.PUT("/api/v1/user/profile", uh.handlerUpdateProfile(), mm.CheckCSRF)
	e.PUT("/api/v1/user/password", uh.handlerUpdatePassword(), mm.CheckCSRF)
	e.PUT("/api/v1/user/photo", uh.handlerUpdateAvatar(), middleware.BodyLimit("10M"))
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
		cookie, err := ctx.Cookie(ConstSessionName)
		if err != nil {
			errResp := NewErrorResponse(ErrNotAuthorized, err)
			return RespondWithError(errResp, ctx)
		}

		session, errResp := uh.sessionUsecase.GetByID(cookie.Value)
		if errResp != nil {
			return RespondWithError(errResp, ctx)
		}

		user, errResp := uh.userUsecase.GetById(session.UserID)
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

		cookie, err := ctx.Cookie(ConstSessionName)
		if err != nil {
			errResp := NewErrorResponse(ErrNotAuthorized, err)
			return RespondWithError(errResp, ctx)
		}

		session, errResp := uh.sessionUsecase.GetByID(cookie.Value)
		if errResp != nil {
			return RespondWithError(errResp, ctx)
		}

		user, errResp := uh.userUsecase.UpdateProfile(session.UserID, req.Name, req.Email)
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

		cookie, err := ctx.Cookie(ConstSessionName)
		if err != nil {
			errResp := NewErrorResponse(ErrNotAuthorized, err)
			return RespondWithError(errResp, ctx)
		}

		session, errResp := uh.sessionUsecase.GetByID(cookie.Value)
		if errResp != nil {
			return RespondWithError(errResp, ctx)
		}

		errResp = uh.userUsecase.UpdatePassword(session.UserID, req.OldPassword, req.Password)
		if errResp != nil {
			return RespondWithError(errResp, ctx)
		}

		return ctx.JSON(http.StatusOK, OKResponse)
	}
}

func (uh *UserHandler) handlerUpdateAvatar() echo.HandlerFunc {
	return func(ctx echo.Context) error {
		cookie, err := ctx.Cookie(ConstSessionName)
		if err != nil {
			errResp := NewErrorResponse(ErrNotAuthorized, err)
			return RespondWithError(errResp, ctx)
		}

		session, errResp := uh.sessionUsecase.GetByID(cookie.Value)
		if errResp != nil {
			return RespondWithError(errResp, ctx)
		}

		user, errResp := uh.userUsecase.GetById(session.UserID)
		if errResp != nil {
			return RespondWithError(errResp, ctx)
		}

		photoUploader := photo_uploader.PhotoUploader{}
		avatarPath, err := photoUploader.UploadPhoto(ctx, "avatar", "./avatars/")
		if err != nil {
			errResp := NewErrorResponse(ErrBadRequest, err)
			return RespondWithError(errResp, ctx)
		}

		user, errResp = uh.userUsecase.UpdateAvatar(session.UserID, avatarPath)
		if errResp != nil {
			return RespondWithError(errResp, ctx)
		}

		return ctx.JSON(http.StatusOK, user)
	}
}
