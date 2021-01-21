package main_test

import (
	"encoding/json"
	"github.com/go-park-mail-ru/2020_2_CodeExpress/internal/models"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"testing"
)

var (
	defaultName = "name"
	expectedUser = &models.User{
		ID:       1,
		Name:     "DaaSsjhdf",
		Email:    "LAKSJJDdks@maail.ru",
		Password: "12345678910",
		Avatar:   "",
	}
)


func register(t *testing.T, client http.Client) ([]*http.Cookie, string) {
	type Request struct {
		Name             string `json:"username" validate:"required"`
		Email            string `json:"email" validate:"required,email"`
		Password         string `json:"password" validate:"required,gte=8"`
		RepeatedPassword string `json:"repeated_password" validate:"required,eqfield=Password"`
	}

	request := &Request{
		Name: expectedUser.Name,
		Email: expectedUser.Email,
		Password: expectedUser.Password,
		RepeatedPassword: expectedUser.Password,
	}

	jsonReq, err := json.Marshal(request)

	assert.Nil(t, err)

	body := strings.NewReader(string(jsonReq))

	resp, err := client.Post("http://127.0.0.1:8085/api/v1/user", "application/json", body)

	assert.Nil(t, err)
	assert.Equal(t, resp.StatusCode, 200)

	resBody, err := ioutil.ReadAll(resp.Body)

	assert.Nil(t, err)

	jsonExpectedUser, err := json.Marshal(expectedUser)

	assert.Equal(t, err, nil)
	assert.Equal(t, string(jsonExpectedUser), string(resBody))

	return resp.Cookies(), resp.Header.Get("X-CSRF-TOKEN")
}

func logout(t *testing.T, client http.Client, cookie *http.Cookie, csrfToken string) {
	request, err := http.NewRequest("DELETE", "http://127.0.0.1:8085/api/v1/session", nil)

	assert.Nil(t, err)

	request.Header.Add("Content-Type", "application/json")
	request.AddCookie(cookie)
	request.Header.Set("X-CSRF-TOKEN", csrfToken)

	resp, err := client.Do(request)
	assert.Nil(t, err)

	assert.Equal(t, resp.StatusCode, http.StatusOK)
}

func login(t *testing.T, client http.Client) ([]*http.Cookie, string) {
	type Request struct {
		Login    string `json:"login" validate:"required"`
		Password string `json:"password" validate:"required"`
	}

	request := &Request{
		Login: expectedUser.Email,
		Password: expectedUser.Password,
	}

	jsonReq, err := json.Marshal(request)
	assert.Nil(t, err)

	body := strings.NewReader(string(jsonReq))

	resp, err := client.Post("http://127.0.0.1:8085/api/v1/session", "application/json", body)
	assert.Nil(t, err)

	assert.Equal(t, resp.StatusCode, 200)

	cookie := resp.Cookies()
	csrfToken := resp.Header.Get("X-CSRF-TOKEN")

	return cookie, csrfToken
}

func getProfile(t *testing.T, client http.Client, cookie *http.Cookie) {
	request, err := http.NewRequest("GET", "http://127.0.0.1:8085/api/v1/user", nil)

	assert.Nil(t, err)

	request.AddCookie(cookie)
	resp, err := client.Do(request)
	assert.Nil(t, err)

	assert.Equal(t, resp.StatusCode, 200)

	resBody, err := ioutil.ReadAll(resp.Body)
	assert.Nil(t, err)

	jsonExpectedUser, err := json.Marshal(expectedUser)

	assert.Equal(t, err, nil)
	assert.Equal(t, string(jsonExpectedUser), string(resBody))
}

func changeProfile(t *testing.T, client http.Client, cookie *http.Cookie, csrfToken string) {
	type Request struct {
		Name  string `json:"username" validate:"required"`
		Email string `json:"email" validate:"required,email"`
	}

	expectedUser.Name += "new"
	expectedUser.Email = expectedUser.Name + "@mail.ru"

	reqChangeProfile := &Request{
		Name: expectedUser.Name,
		Email: expectedUser.Email,
	}

	jsonReq, err := json.Marshal(reqChangeProfile)

	assert.Nil(t, err)

	body := strings.NewReader(string(jsonReq))

	request, err := http.NewRequest("PUT", "http://127.0.0.1:8085/api/v1/user/profile", body)

	assert.Nil(t, err)

	request.Header.Add("Content-Type", "application/json")
	request.AddCookie(cookie)
	request.Header.Set("X-CSRF-TOKEN", csrfToken)

	resp, err := client.Do(request)

	assert.Nil(t, err)

	assert.Equal(t, resp.StatusCode, http.StatusOK)

	resBody, err := ioutil.ReadAll(resp.Body)
	assert.Nil(t, err)

	jsonExpectedUser, err := json.Marshal(expectedUser)

	assert.Equal(t, err, nil)
	assert.Equal(t, string(jsonExpectedUser), string(resBody))
}

func changePassword(t *testing.T, client http.Client, cookie *http.Cookie, csrfToken string) {
	type Request struct {
		OldPassword      string `json:"old_password" validate:"required"`
		Password         string `json:"password" validate:"required,gte=8"`
		RepeatedPassword string `json:"repeated_password" validate:"required,eqfield=Password"`
	}

	newPassowrd := "newpassword"

	reqChangePassword := &Request{
		OldPassword: expectedUser.Password,
		Password: newPassowrd,
		RepeatedPassword: newPassowrd,
	}

	jsonReq, err := json.Marshal(reqChangePassword)
	assert.Nil(t, err)

	body := strings.NewReader(string(jsonReq))

	request, err := http.NewRequest("PUT", "http://127.0.0.1:8085/api/v1/user/password", body)
	assert.Nil(t, err)

	request.Header.Add("Content-Type", "application/json")
	request.AddCookie(cookie)
	request.Header.Set("X-CSRF-TOKEN", csrfToken)

	resp, err := client.Do(request)

	assert.Nil(t, err)

	assert.Equal(t, resp.StatusCode, http.StatusOK)
}

func TestEndToEnd(t *testing.T) {
	client := http.Client{}

	for i := 0; i <= 1000; i++ {
		expectedUser.ID = uint64(i) + 1
		expectedUser.Name = defaultName + strconv.Itoa(i)
		expectedUser.Email = defaultName + strconv.Itoa(i) + "@mail.ru"

		cookie, csrfToken := register(t, client)

		logout(t, client, cookie[0], csrfToken)

		cookie, csrfToken = login(t, client)

		getProfile(t, client, cookie[0])

		changeProfile(t, client, cookie[0], csrfToken)

		changePassword(t, client, cookie[0], csrfToken)

		logout(t, client, cookie[0], csrfToken)
	}
}
