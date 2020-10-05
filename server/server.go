package server

import (
	"log"
	"net/http"

	"github.com/go-park-mail-ru/2020_2_CodeExpress/handlers"
	"github.com/go-park-mail-ru/2020_2_CodeExpress/repositories"
)

func AddCORS(w http.ResponseWriter) {
	w.Header().Add("Access-Control-Allow-Origin", "http://localhost:3000/")
	w.Header().Add("Access-Control-Allow-Methods", "DELETE, POST, GET, OPTIONS")
	w.Header().Add("Access-Control-Allow-Credentials", "true")
	w.Header().Add("Access-Control-Allow-Headers", "Content-Type, Access-Control-Allow-Headers, Authorization, X-Requested-With")
}

func ServerStart(url string, port string) {
	UserRep := repositories.NewUserRepImpl()
	SessionRep := repositories.NewSessionRepImpl()

	UserHandler := handlers.NewUserHandler(UserRep, SessionRep)

	http.HandleFunc("/api/v1/user/register", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		AddCORS(w)

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		if r.Method == http.MethodPost {
			UserHandler.HandleCreateUser(w, r)
			return
		}
		http.Error(w, `"error": "Method not allowed"`, http.StatusMethodNotAllowed)
	})

	http.HandleFunc("/api/v1/user/login", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		AddCORS(w)
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		if r.Method == http.MethodPost {
			UserHandler.HandleLogInUser(w, r)
			return
		}
		http.Error(w, `"error": "Method not allowed"`, http.StatusMethodNotAllowed)
	})

	http.HandleFunc("/api/v1/user/logout", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		AddCORS(w)
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		if r.Method == http.MethodPost {
			UserHandler.HandleLogOutUser(w, r)
			return
		}
		http.Error(w, `"error": "Method not allowed"`, http.StatusMethodNotAllowed)
	})

	log.Println("Server listening on ", url+port)
	http.ListenAndServe(url+port, nil)
}
