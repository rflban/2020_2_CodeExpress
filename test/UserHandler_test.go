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

	"github.com/go-park-mail-ru/2020_2_CodeExpress/repositories"
)

type TestCase struct {
	input  map[string]string
	output map[string]interface{}
}

func TestSignUpSuccess(t *testing.T) {
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
				"id":       0.0,
				"username": "Danaal",
				"email":    "dai@yandaex.ru",
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
				"id":       1.0,
				"username": "Danaal2",
				"email":    "dai2@yandaex.ru",
			},
		},
	}

	ts := httptest.NewServer(http.HandlerFunc(UserHandler.HandleCreateUser))

	for _, testCase := range testCases {
		requestBody, err := json.Marshal(testCase.input)
		if err != nil {
			t.Fatalf("TestSignUpSuccess failed on error %s", err)
		}

		resp, err := http.Post(ts.URL, "application/json", bytes.NewBuffer(requestBody))
		if err != nil {
			t.Fatalf("TestSignUpSuccess failed on error %s", err)
		}

		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			t.Fatalf("TestSignUpSuccess failed on error %s", err)
		}

		jsonBody := make(map[string]interface{})

		err = json.Unmarshal(body, &jsonBody)
		if err != nil {
			t.Fatalf("TestSignUpSuccess failed on error %s", err)
		}

		if !reflect.DeepEqual(jsonBody, testCase.output) {
			t.Fatalf("TestSignUpSuccess failed, expected: %s, result: %s", testCase.output, jsonBody)
		}

		resp.Body.Close()
	}

	ts.Close()
}

func TestSignUpEmailFail(t *testing.T) {
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
				"id":       0.0,
				"username": "Danaal",
				"email":    "dai@yandaex.ru",
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
				"error": "Email already exists",
			},
		},
	}

	ts := httptest.NewServer(http.HandlerFunc(UserHandler.HandleCreateUser))

	for _, testCase := range testCases {
		requestBody, err := json.Marshal(testCase.input)
		if err != nil {
			t.Fatalf("TestSignUpEmail failed on error %s", err)
		}

		resp, err := http.Post(ts.URL, "application/json", bytes.NewBuffer(requestBody))
		if err != nil {
			t.Fatalf("TestSignUpEmail failed on error %s", err)
		}

		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			t.Fatalf("TestSignUpEmail failed on error %s", err)
		}

		jsonBody := make(map[string]interface{})

		err = json.Unmarshal(body, &jsonBody)
		if err != nil {
			t.Fatalf("TestSignUpEmail failed on error %s", err)
		}

		if !reflect.DeepEqual(jsonBody, testCase.output) {
			t.Fatalf("TestSignUpEmail failed, expected: %s, result: %s", testCase.output, jsonBody)
		}

		resp.Body.Close()
	}

	ts.Close()
}

func TestSignUpUsernameFail(t *testing.T) {
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
				"id":       0.0,
				"username": "Danaal",
				"email":    "dai@yandaex.ru",
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
				"error": "Username already exists",
			},
		},
	}

	ts := httptest.NewServer(http.HandlerFunc(UserHandler.HandleCreateUser))

	for _, testCase := range testCases {
		requestBody, err := json.Marshal(testCase.input)
		if err != nil {
			t.Fatalf("TestSignUpUsername failed on error %s", err)
		}

		resp, err := http.Post(ts.URL, "application/json", bytes.NewBuffer(requestBody))
		if err != nil {
			t.Fatalf("TestSignUpUsername failed on error %s", err)
		}

		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			t.Fatalf("TestSignUpUsername failed on error %s", err)
		}

		jsonBody := make(map[string]interface{})

		err = json.Unmarshal(body, &jsonBody)
		if err != nil {
			t.Fatalf("TestSignUpUsername failed on error %s", err)
		}

		if !reflect.DeepEqual(jsonBody, testCase.output) {
			t.Fatalf("TestSignUpUsername failed, expected: %s, result: %s", testCase.output, jsonBody)
		}

		resp.Body.Close()
	}

	ts.Close()
}

func TestSignUpPasswordFail(t *testing.T) {
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
				"error": "Passwords do not match",
			},
		},
	}

	ts := httptest.NewServer(http.HandlerFunc(UserHandler.HandleCreateUser))

	for _, testCase := range testCases {
		requestBody, err := json.Marshal(testCase.input)
		if err != nil {
			t.Fatalf("TestSignUpPassword failed on error %s", err)
		}

		resp, err := http.Post(ts.URL, "application/json", bytes.NewBuffer(requestBody))
		if err != nil {
			t.Fatalf("TestSignUpPassword failed on error %s", err)
		}

		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			t.Fatalf("TestSignUpPassword failed on error %s", err)
		}

		jsonBody := make(map[string]interface{})

		err = json.Unmarshal(body, &jsonBody)
		if err != nil {
			t.Fatalf("TestSignUpPassword failed on error %s", err)
		}

		if !reflect.DeepEqual(jsonBody, testCase.output) {
			t.Fatalf("TestSignUpPassword failed, expected: %s, result: %s", testCase.output, jsonBody)
		}

		resp.Body.Close()
	}

	ts.Close()
}
