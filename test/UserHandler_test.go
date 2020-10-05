package test

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/go-park-mail-ru/2020_2_CodeExpress/handlers"
	"github.com/go-park-mail-ru/2020_2_CodeExpress/models"

	"github.com/go-park-mail-ru/2020_2_CodeExpress/repositories"
)

type TestCase struct {
	input  map[string]string
	output map[string]interface{}
}

func HandlerTest(testCase TestCase, t *testing.T, ts *httptest.Server, testName string) {
	requestBody, err := json.Marshal(testCase.input)
	if err != nil {
		t.Fatalf("%s not failed on error %s", testName, err)
	}

	resp, err := http.Post(ts.URL, "application/json", bytes.NewBuffer(requestBody))
	if err != nil {
		t.Fatalf("%s not failed on error %s", testName, err)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("%s not failed on error %s", testName, err)
	}

	jsonBody := make(map[string]interface{})

	err = json.Unmarshal(body, &jsonBody)
	if err != nil {
		t.Fatalf("%s not failed on error %s", testName, err)
	}

	if !reflect.DeepEqual(jsonBody, testCase.output) {
		t.Fatalf("%s failed, expected: %s, result: %s", testName, testCase.output, jsonBody)
	}
	resp.Body.Close()
}

func TestSignUpSuccess(t *testing.T) {
	testName := "TestSignUpSuccess"

	suRepImpl := repositories.NewUserRepImpl()
	sesRepImpl := repositories.NewSessionRepImpl()

	UserHandler := handlers.NewUserHandler(suRepImpl, sesRepImpl)

	testCases := []TestCase{
		TestCase{
			input: map[string]string{
				"username":          "Danaal",
				"email":             "dai@yandaex.ru",
				"password":          "1234pass",
				"repeated_password": "1234pass",
			},
			output: map[string]interface{}{
				"message": handlers.NoError,
			},
		},
		TestCase{
			input: map[string]string{
				"username":          "Danaal2",
				"email":             "dai2@yandaex.ru",
				"password":          "12342pass",
				"repeated_password": "12342pass",
			},
			output: map[string]interface{}{
				"message": handlers.NoError,
			},
		},
	}

	ts := httptest.NewServer(http.HandlerFunc(UserHandler.HandleCreateUser))

	for _, testCase := range testCases {
		HandlerTest(testCase, t, ts, testName)
	}

	ts.Close()
}

func TestSignUpEmailFail(t *testing.T) {
	testName := "TestSignUpEmail"

	suRepImpl := repositories.NewUserRepImpl()
	sesRepImpl := repositories.NewSessionRepImpl()

	UserHandler := handlers.NewUserHandler(suRepImpl, sesRepImpl)

	testCases := []TestCase{
		TestCase{
			input: map[string]string{
				"username":          "Danaal",
				"email":             "dai@yandaex.ru",
				"password":          "1234pass",
				"repeated_password": "1234pass",
			},
			output: map[string]interface{}{
				"message": handlers.NoError,
			},
		},
		TestCase{
			input: map[string]string{
				"username":          "Danaal2",
				"email":             "dai@yandaex.ru",
				"password":          "1234pass",
				"repeated_password": "1234pass",
			},
			output: map[string]interface{}{
				"message": "Email already exists",
			},
		},
	}

	ts := httptest.NewServer(http.HandlerFunc(UserHandler.HandleCreateUser))

	for _, testCase := range testCases {
		HandlerTest(testCase, t, ts, testName)
	}

	ts.Close()
}

func TestSignUpUsernameFail(t *testing.T) {
	testName := "TestSignUpUsername"

	suRepImpl := repositories.NewUserRepImpl()
	sesRepImpl := repositories.NewSessionRepImpl()

	UserHandler := handlers.NewUserHandler(suRepImpl, sesRepImpl)

	testCases := []TestCase{
		TestCase{
			input: map[string]string{
				"username":          "Danaal",
				"email":             "dai@yandaex.ru",
				"password":          "1234pass",
				"repeated_password": "1234pass",
			},
			output: map[string]interface{}{
				"message": handlers.NoError,
			},
		},
		TestCase{
			input: map[string]string{
				"username":          "Danaal",
				"email":             "dai2@yandaex.ru",
				"password":          "1234pass",
				"repeated_password": "1234pass",
			},
			output: map[string]interface{}{
				"message": "Username already exists",
			},
		},
	}

	ts := httptest.NewServer(http.HandlerFunc(UserHandler.HandleCreateUser))

	for _, testCase := range testCases {
		HandlerTest(testCase, t, ts, testName)
	}

	ts.Close()
}

func TestSignUpPasswordFail(t *testing.T) {
	testName := "TestSignUpPassword"

	suRepImpl := repositories.NewUserRepImpl()
	sesRepImpl := repositories.NewSessionRepImpl()

	UserHandler := handlers.NewUserHandler(suRepImpl, sesRepImpl)

	testCases := []TestCase{
		TestCase{
			input: map[string]string{
				"username":          "Danaal",
				"email":             "dai@yandaex.ru",
				"password":          "1234pass",
				"repeated_password": "1234passa",
			},
			output: map[string]interface{}{
				"message": "Passwords do not match",
			},
		},
	}

	ts := httptest.NewServer(http.HandlerFunc(UserHandler.HandleCreateUser))

	for _, testCase := range testCases {
		HandlerTest(testCase, t, ts, testName)
	}

	ts.Close()
}

func TestSignUpNoEmailFail(t *testing.T) {
	testName := "TestSignUpNoEmail"

	suRepImpl := repositories.NewUserRepImpl()
	sesRepImpl := repositories.NewSessionRepImpl()

	UserHandler := handlers.NewUserHandler(suRepImpl, sesRepImpl)

	testCases := []TestCase{
		TestCase{
			input: map[string]string{
				"username":          "Danaal",
				"email":             "",
				"password":          "1234pass",
				"repeated_password": "1234passa",
			},
			output: map[string]interface{}{
				"message": handlers.NoEmail,
			},
		},
	}

	ts := httptest.NewServer(http.HandlerFunc(UserHandler.HandleCreateUser))

	for _, testCase := range testCases {
		HandlerTest(testCase, t, ts, testName)
	}

	ts.Close()
}

func TestSignUpNoUsernameFail(t *testing.T) {
	testName := "TestSignUpNoUsername"

	suRepImpl := repositories.NewUserRepImpl()
	sesRepImpl := repositories.NewSessionRepImpl()

	UserHandler := handlers.NewUserHandler(suRepImpl, sesRepImpl)

	testCases := []TestCase{
		TestCase{
			input: map[string]string{
				"username":          "",
				"email":             "random@mail.ru",
				"password":          "1234pass",
				"repeated_password": "1234passa",
			},
			output: map[string]interface{}{
				"message": handlers.NoUsername,
			},
		},
	}

	ts := httptest.NewServer(http.HandlerFunc(UserHandler.HandleCreateUser))

	for _, testCase := range testCases {
		HandlerTest(testCase, t, ts, testName)
	}

	ts.Close()
}

func TestSignUpShortPasswordFail(t *testing.T) {
	testName := "TestSignUpShortPassword"

	suRepImpl := repositories.NewUserRepImpl()
	sesRepImpl := repositories.NewSessionRepImpl()

	UserHandler := handlers.NewUserHandler(suRepImpl, sesRepImpl)

	testCases := []TestCase{
		TestCase{
			input: map[string]string{
				"username":          "Daniil",
				"email":             "random@mail.ru",
				"password":          "12",
				"repeated_password": "12",
			},
			output: map[string]interface{}{
				"message": handlers.PasswordTooShort,
			},
		},
	}

	ts := httptest.NewServer(http.HandlerFunc(UserHandler.HandleCreateUser))

	for _, testCase := range testCases {
		HandlerTest(testCase, t, ts, testName)
	}

	ts.Close()
}

func TestSignUpNoPasswordFail(t *testing.T) {
	testName := "TestSignUpNoPassword"

	suRepImpl := repositories.NewUserRepImpl()
	sesRepImpl := repositories.NewSessionRepImpl()

	UserHandler := handlers.NewUserHandler(suRepImpl, sesRepImpl)

	testCases := []TestCase{
		TestCase{
			input: map[string]string{
				"username":          "Daniil",
				"email":             "random@mail.ru",
				"password":          "",
				"repeated_password": "1212414",
			},
			output: map[string]interface{}{
				"message": handlers.NoPassword,
			},
		},
	}

	ts := httptest.NewServer(http.HandlerFunc(UserHandler.HandleCreateUser))

	for _, testCase := range testCases {
		HandlerTest(testCase, t, ts, testName)
	}

	ts.Close()
}

func TestSignUpNoRepeatedPasswordFail(t *testing.T) {
	testName := "TestSignUpNoRepeatedPassword"

	suRepImpl := repositories.NewUserRepImpl()
	sesRepImpl := repositories.NewSessionRepImpl()

	UserHandler := handlers.NewUserHandler(suRepImpl, sesRepImpl)

	testCases := []TestCase{
		TestCase{
			input: map[string]string{
				"username":          "Daniil",
				"email":             "random@mail.ru",
				"password":          "12124142",
				"repeated_password": "",
			},
			output: map[string]interface{}{
				"message": handlers.NoRepeatedPassword,
			},
		},
	}

	ts := httptest.NewServer(http.HandlerFunc(UserHandler.HandleCreateUser))

	for _, testCase := range testCases {
		HandlerTest(testCase, t, ts, testName)
	}

	ts.Close()
}

func TestLogInSuccess(t *testing.T) {
	testName := "TestLogInSuccess"

	suRepImpl := repositories.NewUserRepImpl()
	sesRepImpl := repositories.NewSessionRepImpl()

	user := &models.User{
		ID:       0,
		Name:     "Daniil",
		Email:    "dai@yandaex.ru",
		Password: "123456qQ",
	}
	suRepImpl.CreateUser(user)

	UserHandler := handlers.NewUserHandler(suRepImpl, sesRepImpl)

	testCases := []TestCase{
		TestCase{
			input: map[string]string{
				"login":    "Daniil",
				"password": "123456qQ",
			},
			output: map[string]interface{}{
				"id":       0.0,
				"username": "Daniil",
				"email":    "dai@yandaex.ru",
				"avatar":   "",
			},
		},
	}

	ts := httptest.NewServer(http.HandlerFunc(UserHandler.HandleLogInUser))

	for _, testCase := range testCases {
		HandlerTest(testCase, t, ts, testName)
	}

	ts.Close()
}

func TestLogInNoUsernameFail(t *testing.T) {
	testName := "TestLogInNoUsername"

	suRepImpl := repositories.NewUserRepImpl()
	sesRepImpl := repositories.NewSessionRepImpl()

	user := &models.User{
		ID:       0,
		Name:     "Danaal",
		Email:    "daai@yandaex.ru",
		Password: "123456qQ",
	}
	suRepImpl.CreateUser(user)

	UserHandler := handlers.NewUserHandler(suRepImpl, sesRepImpl)

	testCases := []TestCase{
		TestCase{
			input: map[string]string{
				"login":    "Daniil",
				"password": "123456qQ",
			},
			output: map[string]interface{}{
				"message": "Incorrect login or password",
			},
		},
	}

	ts := httptest.NewServer(http.HandlerFunc(UserHandler.HandleLogInUser))

	for _, testCase := range testCases {
		HandlerTest(testCase, t, ts, testName)
	}

	ts.Close()
}

func TestLogInPasswordFail(t *testing.T) {
	testName := "TestLogInPassword"

	suRepImpl := repositories.NewUserRepImpl()
	sesRepImpl := repositories.NewSessionRepImpl()

	user := &models.User{
		ID:       0,
		Name:     "Danaal",
		Email:    "daai@yandaex.ru",
		Password: "123456qQ",
	}
	suRepImpl.CreateUser(user)

	UserHandler := handlers.NewUserHandler(suRepImpl, sesRepImpl)

	testCases := []TestCase{
		TestCase{
			input: map[string]string{
				"login":    "Danaal",
				"password": "123456qQQQQ",
			},
			output: map[string]interface{}{
				"message": "Incorrect login or password",
			},
		},
	}

	ts := httptest.NewServer(http.HandlerFunc(UserHandler.HandleLogInUser))

	for _, testCase := range testCases {
		HandlerTest(testCase, t, ts, testName)
	}

	ts.Close()
}

func TestLogInNoFieldUsernameFail(t *testing.T) {
	testName := "TestLogInNoField"

	suRepImpl := repositories.NewUserRepImpl()
	sesRepImpl := repositories.NewSessionRepImpl()

	user := &models.User{
		ID:       0,
		Name:     "Danaal",
		Email:    "daai@yandaex.ru",
		Password: "123456qQ",
	}
	suRepImpl.CreateUser(user)

	UserHandler := handlers.NewUserHandler(suRepImpl, sesRepImpl)

	testCases := []TestCase{
		TestCase{
			input: map[string]string{
				"password": "123456qQ",
			},
			output: map[string]interface{}{
				"message": "no login field",
			},
		},
	}

	ts := httptest.NewServer(http.HandlerFunc(UserHandler.HandleLogInUser))

	for _, testCase := range testCases {
		HandlerTest(testCase, t, ts, testName)
	}

	ts.Close()
}

func TestLogInNoFieldPasswordFail(t *testing.T) {
	testName := "TestLogInNoField"

	suRepImpl := repositories.NewUserRepImpl()
	sesRepImpl := repositories.NewSessionRepImpl()

	user := &models.User{
		ID:       0,
		Name:     "Danaal",
		Email:    "daai@yandaex.ru",
		Password: "123456qQ",
	}
	suRepImpl.CreateUser(user)

	UserHandler := handlers.NewUserHandler(suRepImpl, sesRepImpl)

	testCases := []TestCase{
		TestCase{
			input: map[string]string{
				"login": "Danaal",
			},
			output: map[string]interface{}{
				"message": "no password field",
			},
		},
	}

	ts := httptest.NewServer(http.HandlerFunc(UserHandler.HandleLogInUser))

	for _, testCase := range testCases {
		HandlerTest(testCase, t, ts, testName)
	}

	ts.Close()
}
