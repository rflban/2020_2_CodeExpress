package consts

import (
	"errors"
	"net/http"
)

const (
	ErrInternal = iota
	ErrBadRequest
	ErrEmailAlreadyExist
	ErrNameAlreadyExist
	ErrIncorrectLoginOrPassword
	ErrNotAuthorized
	ErrNoEmail
	ErrNoUsername
	ErrNoAvatar
	ErrWrongOldPassword
	ErrNewPasswordIsOld
	ErrArtistNotExist
	ErrTrackNotExist
	ErrAlbumNotExist
	ErrUserNotExist
	ErrNoFavoritesTracks
	ErrPlaylistNotExist
	ErrEmptySearchQuery
	ErrNotAdmin
	ErrPermissionDenied
)

var Errors = map[int]error{
	ErrInternal:                 errors.New("Внутренняя ошибка сервера"),          //Internal server error
	ErrBadRequest:               errors.New("Некорректный запрос"),                //Bad request received
	ErrEmailAlreadyExist:        errors.New("Email уже существует"),               //Email already exists
	ErrNameAlreadyExist:         errors.New("Имя пользователя уже существует"),    //Name already exists
	ErrIncorrectLoginOrPassword: errors.New("Неверный логин или пароль"),          //Incorrect login or password
	ErrNotAuthorized:            errors.New("Не авторизован"),                     //Not authorized
	ErrNoEmail:                  errors.New("Не заполнено поле Email"),            //No email field
	ErrNoUsername:               errors.New("Не заполнено поле Имя пользователя"), //No username field
	ErrNoAvatar:                 errors.New("Не заполнено поле Аватар"),           //No avatar field
	ErrWrongOldPassword:         errors.New("Неверный старый пароль"),             //Wrong old password
	ErrNewPasswordIsOld:         errors.New("Новый пароль совпадает со старым"),   //New password matches old
	ErrArtistNotExist:           errors.New("Артист не найден"),                   //Artist not found
	ErrTrackNotExist:            errors.New("Трек не найден"),                     //Track not found
	ErrAlbumNotExist:            errors.New("Альбом не найден"),                   //Album not found
	ErrUserNotExist:             errors.New("Пользователь не найден"),
	ErrNoFavoritesTracks:        errors.New("У пользователя нет избранных треков"), //User has no favorite tracks
	ErrPlaylistNotExist:         errors.New("Плейлист не найден"),
	ErrEmptySearchQuery:         errors.New("Пустой запрос на поиск"),
	ErrNotAdmin:                 errors.New("Недостаточно прав"),
	ErrPermissionDenied:         errors.New("В доступе отказано"),
}

var StatusCodes = map[int]int{
	ErrInternal:                 http.StatusInternalServerError,
	ErrBadRequest:               http.StatusBadRequest,
	ErrEmailAlreadyExist:        http.StatusForbidden,
	ErrNameAlreadyExist:         http.StatusForbidden,
	ErrIncorrectLoginOrPassword: http.StatusNotFound,
	ErrNotAuthorized:            http.StatusNotFound,
	ErrNoEmail:                  http.StatusBadRequest,
	ErrNoUsername:               http.StatusBadRequest,
	ErrNoAvatar:                 http.StatusBadRequest,
	ErrWrongOldPassword:         http.StatusBadRequest,
	ErrNewPasswordIsOld:         http.StatusBadRequest,
	ErrArtistNotExist:           http.StatusNotFound,
	ErrTrackNotExist:            http.StatusNotFound,
	ErrAlbumNotExist:            http.StatusNotFound,
	ErrUserNotExist:             http.StatusNotFound,
	ErrNoFavoritesTracks:        http.StatusNotFound,
	ErrPlaylistNotExist:         http.StatusNotFound,
	ErrEmptySearchQuery:         http.StatusBadRequest,
	ErrNotAdmin:                 http.StatusUnauthorized,
	ErrPermissionDenied:         http.StatusForbidden,
}
