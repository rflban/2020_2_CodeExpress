package usecase

import (
	"database/sql"

	. "github.com/go-park-mail-ru/2020_2_CodeExpress/internal/consts"
	"github.com/go-park-mail-ru/2020_2_CodeExpress/internal/models"
	. "github.com/go-park-mail-ru/2020_2_CodeExpress/internal/tools/error_response"
	"github.com/go-park-mail-ru/2020_2_CodeExpress/internal/user"
)

type UserUsecase struct {
	userRep user.UserRep
}

func NewUserUsecase(userRep user.UserRep) *UserUsecase {
	return &UserUsecase{
		userRep: userRep,
	}
}

func (uUc *UserUsecase) Create(name string, email string, password string) (*models.User, *ErrorResponse) {
	users, err := uUc.userRep.SelectByNameOrEmail(name, email)
	if err != nil {
		return nil, NewErrorResponse(ErrInternal, err)
	}
	if len(users) > 0 {
		if users[0].Name == name {
			return nil, NewErrorResponse(ErrNameAlreadyExist, nil)
		}
		return nil, NewErrorResponse(ErrEmailAlreadyExist, nil)
	}

	user, err := uUc.userRep.Insert(name, email, password)
	if err != nil {
		return nil, NewErrorResponse(ErrInternal, err)
	}
	return user, nil
}

func (uUc *UserUsecase) GetUserByLogin(login string, password string) (*models.User, *ErrorResponse) {
	user, err := uUc.userRep.SelectByLogin(login)
	if err == sql.ErrNoRows {
		return nil, NewErrorResponse(ErrIncorrectLoginOrPassword, nil)
	}
	if err != nil {
		return nil, NewErrorResponse(ErrInternal, err)
	}
	if user.Password != password {
		return nil, NewErrorResponse(ErrIncorrectLoginOrPassword, nil)
	}
	return user, nil
}

func (uUc *UserUsecase) GetById(id uint64) (*models.User, *ErrorResponse) {
	user, err := uUc.userRep.SelectById(id)
	if err == sql.ErrNoRows {
		return nil, NewErrorResponse(ErrNotAuthorized, err)
	}
	if err != nil {
		return nil, NewErrorResponse(ErrInternal, err)
	}
	return user, nil
}

func (uUc *UserUsecase) UpdateProfile(id uint64, name string, email string) (*models.User, *ErrorResponse) {
	user, errResp := uUc.GetById(id)
	if errResp != nil {
		return nil, errResp
	}

	users, err := uUc.userRep.SelectByNameOrEmail(name, email)
	if err != nil {
		return nil, NewErrorResponse(ErrInternal, err)
	}
	for _, existingUser := range users {
		if existingUser.ID == user.ID {
			continue
		}
		if existingUser.Name == name {
			return nil, NewErrorResponse(ErrNameAlreadyExist, nil)
		}
		return nil, NewErrorResponse(ErrEmailAlreadyExist, nil)
	}

	user.Name = name
	user.Email = email
	if err := uUc.userRep.Update(user); err != nil {
		return nil, NewErrorResponse(ErrInternal, err)
	}
	return user, nil
}

func (uUc *UserUsecase) UpdatePassword(id uint64, oldPassword string, newPassword string) *ErrorResponse {
	user, err := uUc.userRep.SelectById(id)
	if err != nil {
		return NewErrorResponse(ErrInternal, err)
	}

	if user.Password != oldPassword {
		return NewErrorResponse(ErrWrongOldPassword, nil)
	}
	if oldPassword == newPassword {
		return NewErrorResponse(ErrNewPasswordIsOld, nil)
	}

	user.Password = newPassword
	if err := uUc.userRep.Update(user); err != nil {
		return NewErrorResponse(ErrInternal, err)
	}
	return nil
}

func (uUc *UserUsecase) UpdateAvatar(id uint64, avatar string) (*models.User, *ErrorResponse) {
	user, errResp := uUc.GetById(id)
	if errResp != nil {
		return nil, errResp
	}

	user.Avatar = avatar
	if err := uUc.userRep.Update(user); err != nil {
		return nil, NewErrorResponse(ErrInternal, err)
	}
	return user, nil
}
