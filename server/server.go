package server

import (
	"log"
	"net/http"

	"github.com/go-park-mail-ru/2020_2_CodeExpress/handlers"
	"github.com/go-park-mail-ru/2020_2_CodeExpress/repositories"
)

func ServerStart() {
	SignUpRep := repositories.NewSignUpRepImpl()
	SessionRep := repositories.NewSessionRepImpl()

	SignUpHandler := handlers.NewSignUpHandler(SignUpRep, SessionRep)

	go http.HandleFunc("/api/v1/user/register", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		if r.Method == http.MethodPost {
			SignUpHandler.HandleCreateUser(w, r)
			return
		}
		http.Error(w, `"error": "Method not allowed"`, http.StatusMethodNotAllowed)
	})

	port := ":8080"
	log.Println("Server listening on ", port)
	http.ListenAndServe(port, nil)
}
