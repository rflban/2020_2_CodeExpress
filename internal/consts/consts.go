package consts

import "net/http"

const (
	ConstSessionName = "code_express_session_id"
	ConstDaysSession = 1
)

var ConstAllowedOrigins = []string{
	"http://musicexpress.sarafa2n.ru/",
	"localhost",
}

var ConstAllowedMethods = []string{
	http.MethodDelete,
	http.MethodGet,
	http.MethodOptions,
	http.MethodPost,
	http.MethodPut,
}

var ConstAllowedHeaders = []string{
	"Content-Type",
	"Access-Control-Allow-Headers",
	"Authorization",
	"X-Requested-With",
}
