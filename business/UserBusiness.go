package business

import (
	"errors"
	"strings"

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

	err := uRep.CheckUserExists(user)
	if err != nil {
		return nil, err
	}
	err = uRep.CreateUser(user)
	return user, err
}

func UpdateProfile(uRep repositories.UserRep, profileForm *models.ProfileForm, user *models.User) (*models.User, error) {
	var err error
	if profileForm.Email != user.Email {
		newUser := new(models.User)
		newUser.Email = profileForm.Email
		err = uRep.CheckUserExists(newUser)

	}
	if err == nil && profileForm.Username != user.Name {
		newUser := new(models.User)
		newUser.Name = profileForm.Username
		err = uRep.CheckUserExists(newUser)
	}

	if err != nil {
		return nil, err
	}

	user.Email = profileForm.Email
	user.Name = profileForm.Username
	return user, nil
}
