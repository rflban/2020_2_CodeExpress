package delivery

import (
	"net/http"
	"strconv"

	"github.com/go-park-mail-ru/2020_2_CodeExpress/internal/track"

	. "github.com/go-park-mail-ru/2020_2_CodeExpress/internal/consts"
	"github.com/go-park-mail-ru/2020_2_CodeExpress/internal/mwares"
	"github.com/go-park-mail-ru/2020_2_CodeExpress/internal/notification"
	. "github.com/go-park-mail-ru/2020_2_CodeExpress/internal/tools/error_response"
	. "github.com/go-park-mail-ru/2020_2_CodeExpress/internal/tools/responser"
	"github.com/go-park-mail-ru/2020_2_CodeExpress/internal/tools/validator"
	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
)

type NotificationHandler struct {
	notificationUsecase notification.NotificationUsecase
	trackUsecase        track.TrackUsecase
	upgrader            *websocket.Upgrader
}

func NewNotificationDelivery(notificationUsecase notification.NotificationUsecase, trackUsecase track.TrackUsecase) *NotificationHandler {
	return &NotificationHandler{
		notificationUsecase: notificationUsecase,
		trackUsecase:        trackUsecase,
		upgrader:            &websocket.Upgrader{},
	}
}

func (nh *NotificationHandler) Configure(e *echo.Echo, mm *mwares.MiddlewareManager) {
	e.POST("api/v1/notification", nh.HandleInitWSConnection(), mm.CheckAuth)
	e.GET("api/v1/user/:id/track", nh.HandleGetUserTrack(), mm.CheckAuth)
	e.POST("api/v1/user/track/:id", nh.HandleChangeUserTrack(), mm.CheckAuth)
	e.DELETE("api/v1/user/track", nh.HandleUserNoTrack(), mm.CheckAuth)
	e.POST("api/v1/user/:id/emoji", nh.HandleNofityUserEmoji(), mm.CheckAuth)
}

func (nh *NotificationHandler) HandleInitWSConnection() echo.HandlerFunc {
	return func(ctx echo.Context) error {
		ws, err := nh.upgrader.Upgrade(ctx.Response(), ctx.Request(), nil)
		if err != nil {
			return RespondWithError(NewErrorResponse(ErrBadRequest, err), ctx)
		}

		userID := ctx.Get(ConstAuthedUserParam).(uint64)
		nh.notificationUsecase.InitWSConnection(userID, ws)
		return ctx.JSON(http.StatusOK, OKResponse)
	}
}

func (nh *NotificationHandler) HandleGetUserTrack() echo.HandlerFunc {
	return func(ctx echo.Context) error {
		receiverUserID, err := strconv.Atoi(ctx.Param("id"))
		userID := ctx.Get(ConstAuthedUserParam).(uint64)

		if err != nil {
			return RespondWithError(NewErrorResponse(ErrBadRequest, err), ctx)
		}
		trackID := nh.notificationUsecase.GetActiveUserTrack(uint64(receiverUserID))

		if trackID == 0 {
			return RespondWithError(NewErrorResponse(ErrTrackNotExist, nil), ctx)
		}

		track, errResp := nh.trackUsecase.GetByID(trackID, userID)
		if errResp != nil {
			return RespondWithError(errResp, ctx)
		}

		return ctx.JSON(http.StatusOK, track)
	}
}

func (nh *NotificationHandler) HandleUserNoTrack() echo.HandlerFunc {
	return func(ctx echo.Context) error {
		userID := ctx.Get(ConstAuthedUserParam).(uint64)
		nh.notificationUsecase.SetUserNoTrack(userID)
		return ctx.JSON(http.StatusOK, OKResponse)
	}
}

func (nh *NotificationHandler) HandleChangeUserTrack() echo.HandlerFunc {
	return func(ctx echo.Context) error {
		trackID, err := strconv.Atoi(ctx.Param("id"))
		if err != nil {
			return RespondWithError(NewErrorResponse(ErrBadRequest, err), ctx)
		}
		userID := ctx.Get(ConstAuthedUserParam).(uint64)
		nh.notificationUsecase.TrackChanged(userID, uint64(trackID))
		return ctx.JSON(http.StatusOK, OKResponse)
	}
}

func (nh *NotificationHandler) HandleNofityUserEmoji() echo.HandlerFunc {
	type Request struct {
		Emoji rune `json:"emoji" validate:"required"`
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

		userID, err := strconv.Atoi(ctx.Param("id"))
		if err != nil {
			return RespondWithError(NewErrorResponse(ErrBadRequest, err), ctx)
		}

		nh.notificationUsecase.NotifyUserEmoji(uint64(userID), req.Emoji)
		return ctx.JSON(http.StatusOK, OKResponse)
	}
}
