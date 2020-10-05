package test

import (
	"testing"

	"github.com/go-park-mail-ru/2020_2_CodeExpress/business"

	"github.com/go-park-mail-ru/2020_2_CodeExpress/models"

	"github.com/go-park-mail-ru/2020_2_CodeExpress/repositories"
)

func TestCheckPasswordsSuccess(t *testing.T) {
	uRep := repositories.NewUserRepImpl()
	newUser := &models.NewUser{
		Name:             "Daniil",
		Email:            "daniil@mail.ru",
		Password:         "123456qwe",
		RepeatedPassword: "123456qwe",
	}

	_, err := business.CreateUser(uRep, newUser)

	if err != nil {
		t.Fatalf("TestCheckPasswords failed on error %s", err)
	}
}

func TestCheckPasswordsFailed(t *testing.T) {
	uRep := repositories.NewUserRepImpl()
	newUser := &models.NewUser{
		Name:             "Daniil",
		Email:            "daniil@mail.ru",
		Password:         "123456qwe",
		RepeatedPassword: "123456qweR",
	}

	_, err := business.CreateUser(uRep, newUser)

	if err == nil {
		t.Fatalf("TestCheckPasswords not failed on error")
	}
}
