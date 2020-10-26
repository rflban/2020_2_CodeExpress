package test

/*import (
	"reflect"
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

func TestChangeProfileSuccess(t *testing.T) {
	uRep := repositories.NewUserRepImpl()
	profileForm := &models.ProfileForm{
		Username: "Daniil",
		Email:    "daniil@mail.ru",
	}
	err := uRep.CreateUser(&models.User{
		Name:  "Danaal",
		Email: "danaal@mail.ru",
	})

	expected := &models.User{
		Name:  "Daniil",
		Email: "daniil@mail.ru",
	}

	if err != nil {
		t.Fatalf("TestChangeProfile failed on error %s", err)
		return
	}

	user, err := uRep.GetUserByID(0)

	if err != nil {
		t.Fatalf("TestChangeProfile failed on error %s", err)
		return
	}

	user, err = business.UpdateProfile(uRep, profileForm, user)

	if err != nil {
		t.Fatalf("TestChangeProfile failed on error %s", err)
		return
	}

	if !reflect.DeepEqual(user, expected) {
		t.Fatalf("TestChangeProfile failed, result not expecteds")
		return
	}
}*/
