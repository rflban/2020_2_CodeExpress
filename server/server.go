package server

import (
	"log"
	"net/http"

	"github.com/go-park-mail-ru/2020_2_CodeExpress/handlers"
	"github.com/go-park-mail-ru/2020_2_CodeExpress/repositories"
)

func ServerStart(url string, port string) {
	UserRep := repositories.NewUserRepImpl()
	SessionRep := repositories.NewSessionRepImpl()

	SignUpHandler := handlers.NewUserHandler(UserRep, SessionRep)

	http.HandleFunc("/api/v1/user/register", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		if r.Method == http.MethodPost {
			SignUpHandler.HandleCreateUser(w, r)
			return
		}
		http.Error(w, `"error": "Method not allowed"`, http.StatusMethodNotAllowed)
	})

	log.Println("Server listening on ", url+port)
	http.ListenAndServe(url+port, nil)
}
