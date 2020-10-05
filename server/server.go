package server

import (
	"log"
	"net/http"

	"github.com/go-park-mail-ru/2020_2_CodeExpress/handlers"
	"github.com/go-park-mail-ru/2020_2_CodeExpress/repositories"
)

func AddCORS(w http.ResponseWriter) {
	w.Header().Add("Access-Control-Allow-Origin", "http://musicexpress.sarafa2n.ru")
	w.Header().Add("Access-Control-Allow-Methods", "DELETE, POST, GET, OPTIONS")
	w.Header().Add("Access-Control-Allow-Credentials", "true")
	w.Header().Add("Access-Control-Allow-Headers", "Content-Type, Access-Control-Allow-Headers, Authorization, X-Requested-With")
}

type HandlerType int

const (
	SignUpHandler = HandlerType(iota)
	LogInHandler
	LogOutHandler
)

func SetHandler(ht HandlerType, UserHandler *handlers.UserHandler, w http.ResponseWriter, r *http.Request) {
	handler := func(w http.ResponseWriter, r *http.Request) {}

	switch ht {
	case SignUpHandler:
		handler = UserHandler.HandleCreateUser
	case LogInHandler:
		handler = UserHandler.HandleLogInUser
	case LogOutHandler:
		handler = UserHandler.HandleLogOutUser
	}

	w.Header().Set("Content-Type", "application/json")
	AddCORS(w)

	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	if r.Method == http.MethodDelete && ht == LogOutHandler {
		handler(w, r)
		return
	}

	if r.Method == http.MethodPost && ht != LogOutHandler {
		handler(w, r)
		return
	}

	http.Error(w, `"error": "Method not allowed"`, http.StatusMethodNotAllowed)
}

func ServerStart(url string, port string) {
	UserRep := repositories.NewUserRepImpl()
	SessionRep := repositories.NewSessionRepImpl()

	UserHandler := handlers.NewUserHandler(UserRep, SessionRep)

	http.HandleFunc("/api/v1/user/register", func(w http.ResponseWriter, r *http.Request) {
		SetHandler(SignUpHandler, UserHandler, w, r)
	})

	http.HandleFunc("/api/v1/user/login", func(w http.ResponseWriter, r *http.Request) {
		SetHandler(LogInHandler, UserHandler, w, r)
	})

	http.HandleFunc("/api/v1/user/logout", func(w http.ResponseWriter, r *http.Request) {
		SetHandler(LogOutHandler, UserHandler, w, r)
	})

	log.Println("Server listening on ", url+port)
	http.ListenAndServe(url+port, nil)
}
