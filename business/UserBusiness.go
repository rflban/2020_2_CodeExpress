package business

import (
	"errors"
	"fmt"
	"strings"

	"github.com/go-park-mail-ru/2020_2_CodeExpress/consts"

	"github.com/go-park-mail-ru/2020_2_CodeExpress/models"
	"github.com/go-park-mail-ru/2020_2_CodeExpress/repositories"
)

func CreateUser(uRep repositories.UserRep, newUser *models.NewUser) (*models.User, error) {
	if strings.Compare(newUser.Password, newUser.RepeatedPassword) != 0 {
		return nil, errors.New("Passwords do not match")
	}

	user := &models.User{
		Name:     newUser.Name,
		Email:    newUser.Email,
		Password: newUser.Password,
		Avatar:   "",
	}

	existedUser, err := uRep.CheckUserExists(user)

	if err != nil {
		return nil, err
	}
	if existedUser != nil {
		if existedUser.Name == user.Name {
			return nil, errors.New(consts.UserNameExists)
		}
		if existedUser.Email == user.Email {
			return nil, errors.New(consts.EMailExists)
		}
	}

	err = uRep.CreateUser(user)
	return user, err
}

func UpdateProfile(uRep repositories.UserRep, profileForm *models.ProfileForm, user *models.User) (*models.User, error) {
	if profileForm.Email != user.Email {
		newUser := &models.User{
			Email: profileForm.Email,
		}
		existedUser, err := uRep.CheckUserExists(newUser)

		if err != nil {
			return nil, err
		}

		if existedUser != nil {
			if existedUser.Name == newUser.Name {
				return nil, errors.New(consts.UserNameExists)
			}
			if existedUser.Email == newUser.Email {
				return nil, errors.New(consts.EMailExists)
			}
		}
	}

	if profileForm.Username != user.Name {
		newUser := &models.User{
			Name: profileForm.Username,
		}
		existedUser, err := uRep.CheckUserExists(newUser)

		if err != nil {
			return nil, err
		}

		if existedUser != nil {
			if existedUser.Name == newUser.Name {
				return nil, errors.New(consts.UserNameExists)
			}
			if existedUser.Email == newUser.Email {
				return nil, errors.New(consts.EMailExists)
			}
		}
	}

	user.Email = profileForm.Email
	user.Name = profileForm.Username

	if err := uRep.ChangeUser(user); err != nil {
		fmt.Println("HERE")
		return nil, err
	}

	return user, nil
}
