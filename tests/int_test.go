package main_test

import (
	"encoding/json"
	"github.com/go-park-mail-ru/2020_2_CodeExpress/internal/models"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"net/http"
	"strings"
	"testing"
)
//
var (
	expectedUser = &models.User{
		ID:       1,
		Name:     "DaaSsjhdf",
		Email:    "LAKSJJDdks@maail.ru",
		Password: "12345678910",
		Avatar:   "",
	}
)


func TestRegister(t *testing.T) {
	client := http.Client{}

	type RequestRegister struct {
		Name             string `json:"username" validate:"required"`
		Email            string `json:"email" validate:"required,email"`
		Password         string `json:"password" validate:"required,gte=8"`
		RepeatedPassword string `json:"repeated_password" validate:"required,eqfield=Password"`
	}

	reqReg := &RequestRegister{
		Name: expectedUser.Name,
		Email: expectedUser.Email,
		Password: expectedUser.Password,
		RepeatedPassword: expectedUser.Password,
	}

	jsonReq, err := json.Marshal(reqReg)
	assert.Nil(t, err)

	body := strings.NewReader(string(jsonReq))

	resp, err := client.Post("http://0.0.0.0:8085/api/v1/user", "application/json", body)
	assert.Nil(t, err)

	assert.Equal(t, resp.StatusCode, 200)

	resBody, err := ioutil.ReadAll(resp.Body)

	jsonExpectedUser, err := json.Marshal(expectedUser)
	print("This: ")
	print(string(resBody))
	assert.Equal(t, err, nil)
	assert.Equal(t, string(jsonExpectedUser), string(resBody))
}

func TestGetProfile(t *testing.T) {
	type Request struct {
		Login            string `json:"login"`
		Password         string `json:"password" validate:"required,gte=8"`
	}
	client := http.Client{}


	req := &Request{
		Login: "LAKSJJDdks@maail.ru",
		Password: "12345678910",
	}

	jsonReq, err := json.Marshal(req)
	assert.Nil(t, err)

	body := strings.NewReader(string(jsonReq))

	resp, err := client.Post("http://0.0.0.0:8085/api/v1/session", "application/json", body)
	assert.Nil(t, err)

	assert.Equal(t, resp.StatusCode, 200)

	cookie := resp.Cookies()
	csrfToken := resp.Header.Get("X-CSRF-TOKEN")

	//Логаут
	reqLogout, err := http.NewRequest("DELETE", "http://0.0.0.0:8085/api/v1/session", nil)
	reqLogout.Header.Add("Content-Type", "application/json")
	reqLogout.AddCookie(cookie[0])
	reqLogout.Header.Set("X-CSRF-TOKEN", csrfToken)

	resp, err = client.Do(reqLogout)
	assert.Nil(t, err)

	assert.Equal(t, resp.StatusCode, 200)
}
