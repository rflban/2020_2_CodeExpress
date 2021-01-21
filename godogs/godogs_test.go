package godogs

import (
    "encoding/json"
    "github.com/go-park-mail-ru/2020_2_CodeExpress/internal/models"
    "github.com/stretchr/testify/assert"
    "io/ioutil"
    "net/http"
    //"strconv"
    "strings"
    "testing"

	flag "github.com/spf13/pflag"
    "github.com/cucumber/godog"
    "github.com/cucumber/godog/colors"
    "fmt"
    "os"
)

var (
    opts = godog.Options{Output: colors.Colored(os.Stdout)}

    client = http.Client{}
    csrfToken = ""
    cookies []*http.Cookie

	expectedUser = &models.User{
		ID:       0,
		Name:     "DaaSsjhdf",
		Email:    "LAKSJJDdks@maail.ru",
		Password: "12345678910",
		Avatar:   "",
	}
)

func init() {
	godog.BindCommandLineFlags("godog.", &opts)
}

func TestMain(m *testing.M) {
	flag.Parse()
	opts.Paths = flag.Args()

	status := godog.TestSuite{
		Name:                "godogs",
		ScenarioInitializer: InitializeScenario,
		Options:             &opts,
	}.Run()

	os.Exit(status)
}

func registrationUser(name string, email string, passwd string) error {
	type Request struct {
		Name             string `json:"username" validate:"required"`
		Email            string `json:"email" validate:"required,email"`
		Password         string `json:"password" validate:"required,gte=8"`
		RepeatedPassword string `json:"repeated_password" validate:"required,eqfield=Password"`
	}

    expectedUser.ID += 1
    expectedUser.Name = name
    expectedUser.Email = email
    expectedUser.Password = passwd

	request := &Request{
		Name: name,
		Email: email,
		Password: passwd,
		RepeatedPassword: passwd,
	}

	jsonReq, err := json.Marshal(request)

    assertActual(assert.Nil, err, "Error")

	body := strings.NewReader(string(jsonReq))

	resp, err := client.Post("http://127.0.0.1:8085/api/v1/user", "application/json", body)

    assertActual(assert.Nil, err, "Error")
    assertExpectedAndActual(assert.Equal, resp.StatusCode, 200, "Error")

	resBody, err := ioutil.ReadAll(resp.Body)

    assertActual(assert.Nil, err, "Error")

	jsonExpectedUser, err := json.Marshal(expectedUser)

    assertActual(assert.Nil, err, "Error")
    assertExpectedAndActual(assert.Equal, string(jsonExpectedUser), string(resBody), "Error")

    cookies = resp.Cookies()
    csrfToken = resp.Header.Get("X-CSRF-TOKEN")

    return nil
}

func manageAccount(name string, email string, passwd string) error {
	type Request struct {
		Name  string `json:"username" validate:"required"`
		Email string `json:"email" validate:"required,email"`
	}

	expectedUser.Name = name
	expectedUser.Email = email

	reqChangeProfile := &Request{
		Name: expectedUser.Name,
		Email: expectedUser.Email,
	}

	jsonReq, err := json.Marshal(reqChangeProfile)

    assertActual(assert.Nil, err, "Error")

	body := strings.NewReader(string(jsonReq))

	request, err := http.NewRequest("PUT", "http://127.0.0.1:8085/api/v1/user/profile", body)

    assertActual(assert.Nil, err, "Error")

	request.Header.Add("Content-Type", "application/json")
	request.Header.Set("X-CSRF-TOKEN", csrfToken)

	resp, err := client.Do(request)

    assertActual(assert.Nil, err, "Error")
    assertExpectedAndActual(assert.Equal, resp.StatusCode, http.StatusOK, "Error")

	resBody, err := ioutil.ReadAll(resp.Body)
    assertActual(assert.Nil, err, "Error")

	jsonExpectedUser, err := json.Marshal(expectedUser)

    assertActual(assert.Nil, err, "Error")
    assertExpectedAndActual(assert.Equal, string(jsonExpectedUser), string(resBody), "Error")

	type Request_ struct {
		OldPassword      string `json:"old_password" validate:"required"`
		Password         string `json:"password" validate:"required,gte=8"`
		RepeatedPassword string `json:"repeated_password" validate:"required,eqfield=Password"`
	}

	newPassowrd := passwd

	reqChangePassword := &Request_{
		OldPassword: expectedUser.Password,
		Password: newPassowrd,
		RepeatedPassword: newPassowrd,
	}

	jsonReq, err = json.Marshal(reqChangePassword)
    assertActual(assert.Nil, err, "Error")

	body = strings.NewReader(string(jsonReq))

	request, err = http.NewRequest("PUT", "http://127.0.0.1:8085/api/v1/user/password", body)
    assertActual(assert.Nil, err, "Error")

	request.Header.Add("Content-Type", "application/json")
	request.Header.Set("X-CSRF-TOKEN", csrfToken)

	resp, err = client.Do(request)

    assertActual(assert.Nil, err, "Error")
    assertExpectedAndActual(assert.Equal, resp.StatusCode, http.StatusOK, "Error")

    return nil
}

func getAndLogout(name string, email string, passwd string) error {
	request, err := http.NewRequest("DELETE", "http://127.0.0.1:8085/api/v1/session", nil)

    assertActual(assert.Nil, err, "Error")

	request.Header.Add("Content-Type", "application/json")
	request.Header.Set("X-CSRF-TOKEN", csrfToken)

	resp, err := client.Do(request)
    assertActual(assert.Nil, err, "Error")

    assertExpectedAndActual(assert.Equal, resp.StatusCode, http.StatusOK, "Error")

    return nil;
}

func InitializeScenario(ctx *godog.ScenarioContext) {
    ctx.BeforeScenario(func(*godog.Scenario) {
        client = http.Client{}
        csrfToken = ""
    })

	ctx.Step(`^registration of ([^"]*) with email ([^"]*) and password ([^"]*)`, registrationUser)
	ctx.Step(`^update profile info to ([^"]*), ([^"]*) and ([^"]*)`, manageAccount)
	ctx.Step(`^get new profile of ([^"]*) with ([^"]*) and password ([^"]*)`,  getAndLogout)
}

func assertExpectedAndActual(a expectedAndActualAssertion, expected, actual interface{}, msgAndArgs ...interface{}) error {
	var t asserter
	a(&t, expected, actual, msgAndArgs...)
	return t.err
}

type expectedAndActualAssertion func(t assert.TestingT, expected, actual interface{}, msgAndArgs ...interface{}) bool

func assertActual(a actualAssertion, actual interface{}, msgAndArgs ...interface{}) error {
	var t asserter
	a(&t, actual, msgAndArgs...)
	return t.err
}

type actualAssertion func(t assert.TestingT, actual interface{}, msgAndArgs ...interface{}) bool

type asserter struct {
	err error
}

func (a *asserter) Errorf(format string, args ...interface{}) {
	a.err = fmt.Errorf(format, args...)
}
