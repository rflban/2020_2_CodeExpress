package test

import (
	"reflect"
	"testing"

	"github.com/go-park-mail-ru/2020_2_CodeExpress/models"

	"github.com/go-park-mail-ru/2020_2_CodeExpress/repositories"
)

func TestCheckUserCreatesSuccess(t *testing.T) {
	createUser := func(s repositories.SignUpRep) (*models.User, error) {
		newUser := &models.NewUser{
			Name:     "Daniil",
			Email:    "daniil@mail.ru",
			Password: "123456dD",
		}
		return s.CreateUser(newUser)
	}
	expectedUser := &models.User{
		ID:       1,
		Name:     "Daniil",
		Email:    "daniil@mail.ru",
		Password: "123456dD",
	}

	sImpl := repositories.NewSignUpRepImpl()

	resultUser, err := createUser(sImpl)
	if err != nil {
		t.Fatalf("CheckUserCreates failed, error:  %s", err)
	}

	if !reflect.DeepEqual(expectedUser, resultUser) {
		t.Fatalf("CheckUserCreates failed, expected: %v, result: %v", expectedUser, resultUser)
	}
}

func TestCheckUserCreatesFailed(t *testing.T) {
	sImpl := repositories.NewSignUpRepImpl()

	createUser := func(s repositories.SignUpRep) (*models.User, error) {
		newUser := &models.NewUser{
			Name:     "Daniil1",
			Email:    "daniil1@mail.ru",
			Password: "123456dD1",
		}
		return s.CreateUser(newUser)
	}
	expectedUser := &models.User{
		ID:       1,
		Name:     "Daniil",
		Email:    "daniil@mail.ru",
		Password: "123456dD",
	}

	resultUser, _ := createUser(sImpl)

	if reflect.DeepEqual(expectedUser, resultUser) {
		t.Fatalf("CheckUserCreates not failed, expected: %v, result: %v", expectedUser, resultUser)
	}
}

func TestCheckUserExistsSuccess(t *testing.T) {
	sImpl := repositories.NewSignUpRepImpl()

	createUser := func(s repositories.SignUpRep) error {
		newUser := &models.NewUser{
			Name:     "Daniil",
			Email:    "daniil@mail.ru",
			Password: "123456dD",
		}
		return s.CheckUserExists(newUser)
	}

	err := createUser(sImpl)
	if err != nil {
		t.Fatalf("TestCheckUserExists failed, error: %s", err)
	}
}

func TestCheckUserExistsEmailFailed(t *testing.T) {
	sImpl := repositories.NewSignUpRepImpl()

	createUser := func(s repositories.SignUpRep) error {
		newUser1 := &models.NewUser{
			Name:     "Daniil21",
			Email:    "daniil@mail.ru",
			Password: "123456dD",
		}

		newUser2 := &models.NewUser{
			Name:     "Daniil",
			Email:    "daniil@mail.ru",
			Password: "123456dD",
		}
		s.CreateUser(newUser1)
		return s.CheckUserExists(newUser2)
	}

	err := createUser(sImpl)
	if err == nil {
		t.Fatalf("TestCheckUserExists not failed")
	}
}

func TestCheckUserExistsUsernameFailed(t *testing.T) {
	sImpl := repositories.NewSignUpRepImpl()

	createUser := func(s repositories.SignUpRep) error {
		newUser1 := &models.NewUser{
			Name:     "Daniil",
			Email:    "daniil2@mail.ru",
			Password: "123456dD",
		}

		newUser2 := &models.NewUser{
			Name:     "Daniil",
			Email:    "daniil@mail.ru",
			Password: "123456dD",
		}
		s.CreateUser(newUser1)
		return s.CheckUserExists(newUser2)
	}

	err := createUser(sImpl)
	if err == nil {
		t.Fatalf("TestCheckUserExists not failed")
	}
}
