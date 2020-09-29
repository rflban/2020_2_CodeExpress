package test

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
	"time"

	"github.com/go-park-mail-ru/2020_2_CodeExpress/handlers"

	"github.com/go-park-mail-ru/2020_2_CodeExpress/repositories"
)

type TestCase struct {
	input  map[string]string
	output map[string]interface{}
}

func TestSignUp(t *testing.T) {
	suRepImpl := repositories.NewSignUpRepImpl()
	sesRepImpl := repositories.NewSessionRepImpl()

	signUpHandler := handlers.NewSignUpHandler(suRepImpl, sesRepImpl)

	testCases := []TestCase{
		TestCase{
			input: map[string]string{
				"username": "Danaal",
				"email":    "dai@yandaex.ru",
				"password": "1234pass",
			},
			output: map[string]interface{}{
				"id":       1.0,
				"username": "Danaal",
				"email":    "dai@yandaex.ru",
			},
		},
		TestCase{
			input: map[string]string{
				"username": "Danaal2",
				"email":    "dai@yandaex.ru",
				"password": "1234pass",
			},
			output: map[string]interface{}{
				"error": "Email already exists",
			},
		},
		TestCase{
			input: map[string]string{
				"username": "Danaal",
				"email":    "dai22@yandaex.ru",
				"password": "1234pass",
			},
			output: map[string]interface{}{
				"error": "Username already exists",
			},
		},
	}

	ts := httptest.NewServer(http.HandlerFunc(signUpHandler.HandleCreateUser))

	for idx, testCase := range testCases {
		requestBody, err := json.Marshal(testCase.input)
		if err != nil {
			t.Fatalf("|%d| test case failed on error %s", idx, err)
		}

		resp, err := http.Post(ts.URL, "application/json", bytes.NewBuffer(requestBody))
		if err != nil {
			t.Fatalf("|%d| test case failed on error %s", idx, err)
		}

		cookies := resp.Cookies()

		for _, cookie := range cookies {
			if cookie.Name != "session_id" {
				t.Fatalf("|%d| test case failed on cookie doesn't exist", idx)
			}
			if !cookie.HttpOnly {
				t.Fatalf("|%d| test case failed on cookie has no HttpOnly flag", idx)
			}
			if time.Until(cookie.Expires) < 0 {
				t.Fatalf("|%d| test case failed on cookie expired", idx)
			}
		}

		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			t.Fatalf("|%d| test case failed on error %s", idx, err)
		}

		jsonBody := make(map[string]interface{})

		err = json.Unmarshal(body, &jsonBody)
		if err != nil {
			t.Fatalf("|%d| test case failed on error %s", idx, err)
		}

		if !reflect.DeepEqual(jsonBody, testCase.output) {
			t.Fatalf("|%d| test case failed, expected: %s, result: %s", idx, testCase.output, jsonBody)
		}

		resp.Body.Close()
	}

	ts.Close()
}
