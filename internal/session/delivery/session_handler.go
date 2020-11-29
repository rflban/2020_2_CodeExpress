package delivery

import (
	"net/http"

	"github.com/go-park-mail-ru/2020_2_CodeExpress/internal/tools/csrf"

	. "github.com/go-park-mail-ru/2020_2_CodeExpress/internal/consts"
	"github.com/go-park-mail-ru/2020_2_CodeExpress/internal/models"
	"github.com/go-park-mail-ru/2020_2_CodeExpress/internal/session"
	"github.com/go-park-mail-ru/2020_2_CodeExpress/internal/tools/builder"
	. "github.com/go-park-mail-ru/2020_2_CodeExpress/internal/tools/error_response"
	. "github.com/go-park-mail-ru/2020_2_CodeExpress/internal/tools/responser"
	"github.com/go-park-mail-ru/2020_2_CodeExpress/internal/tools/validator"
	"github.com/go-park-mail-ru/2020_2_CodeExpress/internal/user"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
)

type SessionHandler struct {
	sessionUsecase session.SessionUsecase
	userUsecase    user.UserUsecase
}

func NewSessionHandler(sessionUsecase session.SessionUsecase, userUsecase user.UserUsecase) *SessionHandler {
	return &SessionHandler{
		sessionUsecase: sessionUsecase,
		userUsecase:    userUsecase,
	}
}

func (sh *SessionHandler) Configure(e *echo.Echo) {
	e.POST("/api/v1/session", sh.HandlerLogin())
	e.DELETE("/api/v1/session", sh.HandlerLogout())
}

func (sh *SessionHandler) HandlerLogin() echo.HandlerFunc {
	type Request struct {
		Login    string `json:"login" validate:"required"`
		Password string `json:"password" validate:"required"`
	}

	return func(ctx echo.Context) error {
		req := &Request{}

		if err := validator.NewRequestValidator(ctx).Validate(req); err != nil {
			if err.Error != nil { //TODO: Обрабатывать ошибки валидатора
				logrus.Info(err.Error)
				validator.GetValidationError(err)
			}
			return ctx.JSON(err.StatusCode, err.UserError)
		}

		user, err := sh.userUsecase.GetUserByLogin(req.Login, req.Password)

		if err != nil {
			return RespondWithError(err, ctx)
		}

		session := models.NewSession(user.ID)

		if err := sh.sessionUsecase.CreateSession(session); err != nil {
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

func (sh *SessionHandler) HandlerLogout() echo.HandlerFunc {
	return func(ctx echo.Context) error {
		cookie, err := ctx.Cookie(ConstSessionName)

		if err != nil {
			errResp := NewErrorResponse(ErrNotAuthorized, err)
			return RespondWithError(errResp, ctx)
		}

		session, errResp := sh.sessionUsecase.GetByID(cookie.Value)

		if errResp != nil {
			return RespondWithError(errResp, ctx)
		}

		errResp = sh.sessionUsecase.DeleteSession(session)

		if errResp != nil {
			return RespondWithError(errResp, ctx)
		}

		cookie = builder.BuildCookie(session)
		ctx.SetCookie(cookie)

		return ctx.JSON(http.StatusOK, OKResponse)
	}
}
