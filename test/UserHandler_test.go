package test

/*import (
	"bytes"
	"encoding/json"
	"github.com/go-park-mail-ru/2020_2_CodeExpress/consts"
	"github.com/go-park-mail-ru/2020_2_CodeExpress/handlers"
	"github.com/go-park-mail-ru/2020_2_CodeExpress/models"
	"github.com/go-park-mail-ru/2020_2_CodeExpress/repositories"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

type TestCase struct {
	input  map[string]string
	output map[string]interface{}
}

func serverStart() *echo.Echo {
	e := echo.New()

	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowCredentials: true,
		AllowHeaders:     []string{"Content-Type", "Access-Control-Allow-Headers", "Authorization", "X-Requested-With"},
		AllowMethods:     []string{http.MethodDelete, http.MethodGet, http.MethodOptions, http.MethodPost},
	}))

	return e
}

func HandlerTest(t *testing.T, e *echo.Echo, url string, handler func(e echo.Context) error, testCase TestCase, testName string) {
	requestBody, err := json.Marshal(testCase.input)
	if err != nil {
		t.Fatalf("%s not failed on error %s", testName, err)
	}

	request := httptest.NewRequest(http.MethodPost, url, bytes.NewReader(requestBody))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	recorder := httptest.NewRecorder()
	context := e.NewContext(request, recorder)

	result := recorder.Result()
	defer func() {
		_ = result.Body.Close()
	}()

	if assert.NoError(t, handler(context)) {
		jsonBody := make(map[string]interface{})
		if err = json.Unmarshal([]byte(recorder.Body.String()), &jsonBody); err != nil {
			t.Fatalf("%s not failed on error %s", testName, err)
		}

		assert.Equal(t, testCase.output, jsonBody)
	}
}

func TestSignUpSuccess(t *testing.T) {
	testName := "TestSignUpSuccess"

	suRepImpl := repositories.NewUserRepImpl()
	sesRepImpl := repositories.NewSessionRepImpl()

	UserHandler := handlers.NewUserHandler(suRepImpl, sesRepImpl)

	testCases := []TestCase{
		{
			input: map[string]string{
				"username":          "Danaal",
				"email":             "dai@yandaex.ru",
				"password":          "1234pass",
				"repeated_password": "1234pass",
			},
			output: map[string]interface{}{
				"message": consts.NoError,
			},
		},
		{
			input: map[string]string{
				"username":          "Danaal2",
				"email":             "dai2@yandaex.ru",
				"password":          "12342pass",
				"repeated_password": "12342pass",
			},
			output: map[string]interface{}{
				"message": consts.NoError,
			},
		},
	}

	e := serverStart()
	for _, testCase := range testCases {
		HandlerTest(t, e, "/api/v1/user/register", UserHandler.HandleCreateUser, testCase, testName)
	}
	_ = e.Close()
}

func TestSignUpEmailFail(t *testing.T) {
	testName := "TestSignUpEmail"

	suRepImpl := repositories.NewUserRepImpl()
	sesRepImpl := repositories.NewSessionRepImpl()

	UserHandler := handlers.NewUserHandler(suRepImpl, sesRepImpl)

	testCases := []TestCase{
		{
			input: map[string]string{
				"username":          "Danaal",
				"email":             "dai@yandaex.ru",
				"password":          "1234pass",
				"repeated_password": "1234pass",
			},
			output: map[string]interface{}{
				"message": consts.NoError,
			},
		},
		{
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

	e := serverStart()
	for _, testCase := range testCases {
		HandlerTest(t, e, "/api/v1/user/register", UserHandler.HandleCreateUser, testCase, testName)
	}
	_ = e.Close()
}

func TestSignUpUsernameFail(t *testing.T) {
	testName := "TestSignUpUsername"

	suRepImpl := repositories.NewUserRepImpl()
	sesRepImpl := repositories.NewSessionRepImpl()

	UserHandler := handlers.NewUserHandler(suRepImpl, sesRepImpl)

	testCases := []TestCase{
		{
			input: map[string]string{
				"username":          "Danaal",
				"email":             "dai@yandaex.ru",
				"password":          "1234pass",
				"repeated_password": "1234pass",
			},
			output: map[string]interface{}{
				"message": consts.NoError,
			},
		},
		{
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

	e := serverStart()
	for _, testCase := range testCases {
		HandlerTest(t, e, "/api/v1/user/register", UserHandler.HandleCreateUser, testCase, testName)
	}
	_ = e.Close()
}

func TestSignUpPasswordFail(t *testing.T) {
	testName := "TestSignUpPassword"

	suRepImpl := repositories.NewUserRepImpl()
	sesRepImpl := repositories.NewSessionRepImpl()

	UserHandler := handlers.NewUserHandler(suRepImpl, sesRepImpl)

	testCases := []TestCase{
		{
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

	e := serverStart()
	for _, testCase := range testCases {
		HandlerTest(t, e, "/api/v1/user/register", UserHandler.HandleCreateUser, testCase, testName)
	}
	_ = e.Close()
}

func TestSignUpNoEmailFail(t *testing.T) {
	testName := "TestSignUpNoEmail"

	suRepImpl := repositories.NewUserRepImpl()
	sesRepImpl := repositories.NewSessionRepImpl()

	UserHandler := handlers.NewUserHandler(suRepImpl, sesRepImpl)

	testCases := []TestCase{
		{
			input: map[string]string{
				"username":          "Danaal",
				"email":             "",
				"password":          "1234pass",
				"repeated_password": "1234pass",
			},
			output: map[string]interface{}{
				"message": consts.NoEmail,
			},
		},
	}

	e := serverStart()
	for _, testCase := range testCases {
		HandlerTest(t, e, "/api/v1/user/register", UserHandler.HandleCreateUser, testCase, testName)
	}
	_ = e.Close()
}

func TestSignUpNoUsernameFail(t *testing.T) {
	testName := "TestSignUpNoUsername"

	suRepImpl := repositories.NewUserRepImpl()
	sesRepImpl := repositories.NewSessionRepImpl()

	UserHandler := handlers.NewUserHandler(suRepImpl, sesRepImpl)

	testCases := []TestCase{
		{
			input: map[string]string{
				"username":          "",
				"email":             "random@mail.ru",
				"password":          "1234pass",
				"repeated_password": "1234passa",
			},
			output: map[string]interface{}{
				"message": consts.NoUsername,
			},
		},
	}

	e := serverStart()
	for _, testCase := range testCases {
		HandlerTest(t, e, "/api/v1/user/register", UserHandler.HandleCreateUser, testCase, testName)
	}
	_ = e.Close()
}

func TestSignUpShortPasswordFail(t *testing.T) {
	testName := "TestSignUpShortPassword"

	suRepImpl := repositories.NewUserRepImpl()
	sesRepImpl := repositories.NewSessionRepImpl()

	UserHandler := handlers.NewUserHandler(suRepImpl, sesRepImpl)

	testCases := []TestCase{
		{
			input: map[string]string{
				"username":          "Daniil",
				"email":             "random@mail.ru",
				"password":          "12",
				"repeated_password": "12",
			},
			output: map[string]interface{}{
				"message": consts.PasswordTooShort,
			},
		},
	}

	e := serverStart()
	for _, testCase := range testCases {
		HandlerTest(t, e, "/api/v1/user/register", UserHandler.HandleCreateUser, testCase, testName)
	}
	_ = e.Close()
}

func TestSignUpNoPasswordFail(t *testing.T) {
	testName := "TestSignUpNoPassword"

	suRepImpl := repositories.NewUserRepImpl()
	sesRepImpl := repositories.NewSessionRepImpl()

	UserHandler := handlers.NewUserHandler(suRepImpl, sesRepImpl)

	testCases := []TestCase{
		{
			input: map[string]string{
				"username":          "Daniil",
				"email":             "random@mail.ru",
				"password":          "",
				"repeated_password": "1212414",
			},
			output: map[string]interface{}{
				"message": consts.NoPassword,
			},
		},
	}

	e := serverStart()
	for _, testCase := range testCases {
		HandlerTest(t, e, "/api/v1/user/register", UserHandler.HandleCreateUser, testCase, testName)
	}
	_ = e.Close()
}

func TestSignUpNoRepeatedPasswordFail(t *testing.T) {
	testName := "TestSignUpNoRepeatedPassword"

	suRepImpl := repositories.NewUserRepImpl()
	sesRepImpl := repositories.NewSessionRepImpl()

	UserHandler := handlers.NewUserHandler(suRepImpl, sesRepImpl)

	testCases := []TestCase{
		{
			input: map[string]string{
				"username":          "Daniil",
				"email":             "random@mail.ru",
				"password":          "12124142",
				"repeated_password": "",
			},
			output: map[string]interface{}{
				"message": consts.NoRepeatedPassword,
			},
		},
	}

	e := serverStart()
	for _, testCase := range testCases {
		HandlerTest(t, e, "/api/v1/user/register", UserHandler.HandleCreateUser, testCase, testName)
	}
	_ = e.Close()
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
	_ = suRepImpl.CreateUser(user)

	UserHandler := handlers.NewUserHandler(suRepImpl, sesRepImpl)

	testCases := []TestCase{
		{
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

	e := serverStart()
	for _, testCase := range testCases {
		HandlerTest(t, e, "/api/v1/user/login", UserHandler.HandleLogInUser, testCase, testName)
	}
	_ = e.Close()
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
	_ = suRepImpl.CreateUser(user)

	UserHandler := handlers.NewUserHandler(suRepImpl, sesRepImpl)

	testCases := []TestCase{
		{
			input: map[string]string{
				"login":    "Daniil",
				"password": "123456qQ",
			},
			output: map[string]interface{}{
				"message": "Incorrect login or password",
			},
		},
	}

	e := serverStart()
	for _, testCase := range testCases {
		HandlerTest(t, e, "/api/v1/user/login", UserHandler.HandleLogInUser, testCase, testName)
	}
	_ = e.Close()
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
	_ = suRepImpl.CreateUser(user)

	UserHandler := handlers.NewUserHandler(suRepImpl, sesRepImpl)

	testCases := []TestCase{
		{
			input: map[string]string{
				"login":    "Danaal",
				"password": "123456qQQQQ",
			},
			output: map[string]interface{}{
				"message": "Incorrect login or password",
			},
		},
	}

	e := serverStart()
	for _, testCase := range testCases {
		HandlerTest(t, e, "/api/v1/user/login", UserHandler.HandleLogInUser, testCase, testName)
	}
	_ = e.Close()
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
	_ = suRepImpl.CreateUser(user)

	UserHandler := handlers.NewUserHandler(suRepImpl, sesRepImpl)

	testCases := []TestCase{
		{
			input: map[string]string{
				"password": "123456qQ",
			},
			output: map[string]interface{}{
				"message": "no username field", //TODO: везде возможно нужно использовать константы из Error...
			},
		},
	}

	e := serverStart()
	for _, testCase := range testCases {
		HandlerTest(t, e, "/api/v1/user/login", UserHandler.HandleLogInUser, testCase, testName)
	}
	_ = e.Close()
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
	_ = suRepImpl.CreateUser(user)

	UserHandler := handlers.NewUserHandler(suRepImpl, sesRepImpl)

	testCases := []TestCase{
		{
			input: map[string]string{
				"login": "Danaal",
			},
			output: map[string]interface{}{
				"message": "no password field",
			},
		},
	}

	e := serverStart()
	for _, testCase := range testCases {
		HandlerTest(t, e, "/api/v1/user/login", UserHandler.HandleLogInUser, testCase, testName)
	}
	_ = e.Close()
}

func TestGetUserSuccess(t *testing.T) {
	testName := "TestGetUserSuccess"

	suRepImpl := repositories.NewUserRepImpl()
	sesRepImpl := repositories.NewSessionRepImpl()

	user := &models.User{
		ID:       0,
		Name:     "Daniil",
		Email:    "dai@yandaex.ru",
		Password: "123456qQ",
	}
	_ = suRepImpl.CreateUser(user)
	userSession := repositories.NewSession(user)
	_ = sesRepImpl.AddSession(userSession)

	userCookie := http.Cookie{
		Name:     userSession.Name,
		Value:    userSession.ID,
		Expires:  userSession.Expire,
		HttpOnly: true,
		Path:     "/",
	}

	UserHandler := handlers.NewUserHandler(suRepImpl, sesRepImpl)

	testCase := TestCase{
		input: map[string]string{
			"code_express_session_id": userCookie.String(),
		},
		output: map[string]interface{}{
			"id":       0.0,
			"username": "Daniil",
			"email":    "dai@yandaex.ru",
			"avatar":   "",
		},
	}

	e := serverStart()

	requestBody, err := json.Marshal(testCase.input)
	if err != nil {
		t.Fatalf("%s not failed on error %s", testName, err)
	}

	request := httptest.NewRequest(http.MethodGet, "/api/v1/user/current", bytes.NewReader(requestBody))
	request.AddCookie(&userCookie)
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	recorder := httptest.NewRecorder()
	context := e.NewContext(request, recorder)

	result := recorder.Result()
	defer func() {
		_ = result.Body.Close()
	}()

	if assert.NoError(t, UserHandler.HandleCurrentUser(context)) {
		jsonBody := make(map[string]interface{})
		if err = json.Unmarshal([]byte(recorder.Body.String()), &jsonBody); err != nil {
			t.Fatalf("%s not failed on error %s", testName, err)
		}

		assert.Equal(t, testCase.output, jsonBody)
	}

	_ = e.Close()
}

func TestLogoutSuccess(t *testing.T) {
	testName := "TestGetUserSuccess"

	suRepImpl := repositories.NewUserRepImpl()
	sesRepImpl := repositories.NewSessionRepImpl()

	user := &models.User{
		ID:       0,
		Name:     "Daniil",
		Email:    "dai@yandaex.ru",
		Password: "123456qQ",
	}
	_ = suRepImpl.CreateUser(user)
	userSession := repositories.NewSession(user)
	_ = sesRepImpl.AddSession(userSession)

	userCookie := http.Cookie{
		Name:     userSession.Name,
		Value:    userSession.ID,
		Expires:  userSession.Expire,
		HttpOnly: true,
		Path:     "/",
	}

	UserHandler := handlers.NewUserHandler(suRepImpl, sesRepImpl)

	testCase := TestCase{
		input: map[string]string{
			"code_express_session_id": userCookie.String(),
		},
		output: map[string]interface{}{
			"message": consts.NoError,
		},
	}

	e := serverStart()

	requestBody, err := json.Marshal(testCase.input)
	if err != nil {
		t.Fatalf("%s not failed on error %s", testName, err)
	}

	request := httptest.NewRequest(http.MethodDelete, "/api/v1/user/logout", bytes.NewReader(requestBody))
	request.AddCookie(&userCookie)
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	recorder := httptest.NewRecorder()
	context := e.NewContext(request, recorder)

	result := recorder.Result()
	defer func() {
		_ = result.Body.Close()
	}()

	if assert.NoError(t, UserHandler.HandleLogOutUser(context)) {
		jsonBody := make(map[string]interface{})
		if err = json.Unmarshal([]byte(recorder.Body.String()), &jsonBody); err != nil {
			t.Fatalf("%s not failed on error %s", testName, err)
		}

		assert.Equal(t, testCase.output, jsonBody)
	}

	_ = e.Close()
}

func TestUpdateProfile(t *testing.T) {
	testName := "TestUpdateProfile"

	suRepImpl := repositories.NewUserRepImpl()
	sesRepImpl := repositories.NewSessionRepImpl()

	user := &models.User{
		ID:       0,
		Name:     "Daniil",
		Email:    "danya@yandaex.ru",
		Password: "123456qQ",
	}
	_ = suRepImpl.CreateUser(user)
	userSession := repositories.NewSession(user)
	_ = sesRepImpl.AddSession(userSession)

	userCookie := http.Cookie{
		Name:     userSession.Name,
		Value:    userSession.ID,
		Expires:  userSession.Expire,
		HttpOnly: true,
		Path:     "/",
	}

	UserHandler := handlers.NewUserHandler(suRepImpl, sesRepImpl)

	testCase := TestCase{
		input: map[string]string{
			"email":    "dai@yandaex.ru",
			"username": "Daniil",
		},
		output: map[string]interface{}{
			"id":       0.0,
			"username": "Daniil",
			"email":    "dai@yandaex.ru",
			"avatar":   "",
		},
	}

	e := serverStart()

	requestBody, err := json.Marshal(testCase.input)
	if err != nil {
		t.Fatalf("%s not failed on error %s", testName, err)
	}

	request := httptest.NewRequest(http.MethodPost, "/api/v1/user/change/profile", bytes.NewReader(requestBody))
	request.AddCookie(&userCookie)
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	recorder := httptest.NewRecorder()
	context := e.NewContext(request, recorder)

	result := recorder.Result()
	defer func() {
		_ = result.Body.Close()
	}()

	if assert.NoError(t, UserHandler.HandleUpdateProfile(context)) {
		jsonBody := make(map[string]interface{})
		if err = json.Unmarshal([]byte(recorder.Body.String()), &jsonBody); err != nil {
			t.Fatalf("%s not failed on error %s", testName, err)
		}

		assert.Equal(t, testCase.output, jsonBody)
	}

	_ = e.Close()
}

func TestUpdatePasswordSuccess(t *testing.T) {
	testName := "TestUpdatePasswordSuccess"

	suRepImpl := repositories.NewUserRepImpl()
	sesRepImpl := repositories.NewSessionRepImpl()

	user := &models.User{
		ID:       0,
		Name:     "Daniil",
		Email:    "dai@yandaex.ru",
		Password: "123456qQ",
	}
	_ = suRepImpl.CreateUser(user)
	userSession := repositories.NewSession(user)
	_ = sesRepImpl.AddSession(userSession)

	userCookie := http.Cookie{
		Name:     userSession.Name,
		Value:    userSession.ID,
		Expires:  userSession.Expire,
		HttpOnly: true,
		Path:     "/",
	}

	UserHandler := handlers.NewUserHandler(suRepImpl, sesRepImpl)

	testCase := TestCase{
		input: map[string]string{
			"password":          "123456789qQ",
			"repeated_password": "123456789qQ",
		},
		output: map[string]interface{}{
			"message": consts.NoError,
		},
	}

	e := serverStart()

	requestBody, err := json.Marshal(testCase.input)
	if err != nil {
		t.Fatalf("%s not failed on error %s", testName, err)
	}

	request := httptest.NewRequest(http.MethodPost, "/api/v1/user/change/password", bytes.NewReader(requestBody))
	request.AddCookie(&userCookie)
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	recorder := httptest.NewRecorder()
	context := e.NewContext(request, recorder)

	result := recorder.Result()
	defer func() {
		_ = result.Body.Close()
	}()

	if assert.NoError(t, UserHandler.HandleUpdatePassword(context)) {
		jsonBody := make(map[string]interface{})
		if err = json.Unmarshal([]byte(recorder.Body.String()), &jsonBody); err != nil {
			t.Fatalf("%s not failed on error %s", testName, err)
		}

		assert.Equal(t, testCase.output, jsonBody)
	}

	_ = e.Close()
}

func TestUpdateShortPassword(t *testing.T) {
	testName := "TestUpdateShortPassword"

	suRepImpl := repositories.NewUserRepImpl()
	sesRepImpl := repositories.NewSessionRepImpl()

	user := &models.User{
		ID:       0,
		Name:     "Daniil",
		Email:    "dai@yandaex.ru",
		Password: "123456qQ",
	}
	_ = suRepImpl.CreateUser(user)
	userSession := repositories.NewSession(user)
	_ = sesRepImpl.AddSession(userSession)

	userCookie := http.Cookie{
		Name:     userSession.Name,
		Value:    userSession.ID,
		Expires:  userSession.Expire,
		HttpOnly: true,
		Path:     "/",
	}

	UserHandler := handlers.NewUserHandler(suRepImpl, sesRepImpl)

	testCase := TestCase{
		input: map[string]string{
			"password":          "123",
			"repeated_password": "123",
		},
		output: map[string]interface{}{
			"message": consts.PasswordTooShort,
		},
	}

	e := serverStart()

	requestBody, err := json.Marshal(testCase.input)
	if err != nil {
		t.Fatalf("%s not failed on error %s", testName, err)
	}

	request := httptest.NewRequest(http.MethodPost, "/api/v1/user/change/password", bytes.NewReader(requestBody))
	request.AddCookie(&userCookie)
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	recorder := httptest.NewRecorder()
	context := e.NewContext(request, recorder)

	result := recorder.Result()
	defer func() {
		_ = result.Body.Close()
	}()

	if assert.NoError(t, UserHandler.HandleUpdatePassword(context)) {
		jsonBody := make(map[string]interface{})
		if err = json.Unmarshal([]byte(recorder.Body.String()), &jsonBody); err != nil {
			t.Fatalf("%s not failed on error %s", testName, err)
		}

		assert.Equal(t, testCase.output, jsonBody)
	}

	_ = e.Close()
}

func TestUpdateNoPassword(t *testing.T) {
	testName := "TestUpdateNoPassword"

	suRepImpl := repositories.NewUserRepImpl()
	sesRepImpl := repositories.NewSessionRepImpl()

	user := &models.User{
		ID:       0,
		Name:     "Daniil",
		Email:    "dai@yandaex.ru",
		Password: "123456qQ",
	}
	_ = suRepImpl.CreateUser(user)
	userSession := repositories.NewSession(user)
	_ = sesRepImpl.AddSession(userSession)

	userCookie := http.Cookie{
		Name:     userSession.Name,
		Value:    userSession.ID,
		Expires:  userSession.Expire,
		HttpOnly: true,
		Path:     "/",
	}

	UserHandler := handlers.NewUserHandler(suRepImpl, sesRepImpl)

	testCase := TestCase{
		input: map[string]string{
			"password":          "",
			"repeated_password": "",
		},
		output: map[string]interface{}{
			"message": consts.NoPassword,
		},
	}

	e := serverStart()

	requestBody, err := json.Marshal(testCase.input)
	if err != nil {
		t.Fatalf("%s not failed on error %s", testName, err)
	}

	request := httptest.NewRequest(http.MethodPost, "/api/v1/user/change/password", bytes.NewReader(requestBody))
	request.AddCookie(&userCookie)
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	recorder := httptest.NewRecorder()
	context := e.NewContext(request, recorder)

	result := recorder.Result()
	defer func() {
		_ = result.Body.Close()
	}()

	if assert.NoError(t, UserHandler.HandleUpdatePassword(context)) {
		jsonBody := make(map[string]interface{})
		if err = json.Unmarshal([]byte(recorder.Body.String()), &jsonBody); err != nil {
			t.Fatalf("%s not failed on error %s", testName, err)
		}

		assert.Equal(t, testCase.output, jsonBody)
	}

	_ = e.Close()
}

func TestUpdateNoSesondPassword(t *testing.T) {
	testName := "TestUpdateNoSesondPassword"

	suRepImpl := repositories.NewUserRepImpl()
	sesRepImpl := repositories.NewSessionRepImpl()

	user := &models.User{
		ID:       0,
		Name:     "Daniil",
		Email:    "dai@yandaex.ru",
		Password: "123456qQ",
	}
	_ = suRepImpl.CreateUser(user)
	userSession := repositories.NewSession(user)
	_ = sesRepImpl.AddSession(userSession)

	userCookie := http.Cookie{
		Name:     userSession.Name,
		Value:    userSession.ID,
		Expires:  userSession.Expire,
		HttpOnly: true,
		Path:     "/",
	}

	UserHandler := handlers.NewUserHandler(suRepImpl, sesRepImpl)

	testCase := TestCase{
		input: map[string]string{
			"password":          "adsf",
			"repeated_password": "",
		},
		output: map[string]interface{}{
			"message": consts.NoRepeatedPassword,
		},
	}

	e := serverStart()

	requestBody, err := json.Marshal(testCase.input)
	if err != nil {
		t.Fatalf("%s not failed on error %s", testName, err)
	}

	request := httptest.NewRequest(http.MethodPost, "/api/v1/user/change/password", bytes.NewReader(requestBody))
	request.AddCookie(&userCookie)
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	recorder := httptest.NewRecorder()
	context := e.NewContext(request, recorder)

	result := recorder.Result()
	defer func() {
		_ = result.Body.Close()
	}()

	if assert.NoError(t, UserHandler.HandleUpdatePassword(context)) {
		jsonBody := make(map[string]interface{})
		if err = json.Unmarshal([]byte(recorder.Body.String()), &jsonBody); err != nil {
			t.Fatalf("%s not failed on error %s", testName, err)
		}
		assert.Equal(t, testCase.output, jsonBody)
	}

	_ = e.Close()
}*/
