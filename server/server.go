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
<<<<<<< HEAD
	EditProfileHandler
	EditPasswordHandler
=======
	SetAvatarHandler
	CurrentProfile
>>>>>>> CP-21: Реализовать механизм загрузки изображений профиля и просмотра профиля
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
<<<<<<< HEAD
	case EditProfileHandler:
		handler = UserHandler.HandleUpdateProfile
	case EditPasswordHandler:
		handler = UserHandler.HandleUpdatePassword
=======
	case SetAvatarHandler:
		handler = UserHandler.HandleSetAvatar
	case CurrentProfile:
		handler = UserHandler.HandleCurrentUser
>>>>>>> CP-21: Реализовать механизм загрузки изображений профиля и просмотра профиля
	}

	w.Header().Set("Content-Type", "application/json")
	AddCORS(w)

	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	if r.Method == http.MethodGet && ht == CurrentProfile {
		handler(w, r)
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

<<<<<<< HEAD
	http.HandleFunc("/api/v1/user/change/profile", func(w http.ResponseWriter, r *http.Request) {
		SetHandler(EditProfileHandler, UserHandler, w, r)
	})

	http.HandleFunc("/api/v1/user/change/password", func(w http.ResponseWriter, r *http.Request) {
		SetHandler(EditPasswordHandler, UserHandler, w, r)
	})

=======
	http.HandleFunc("/api/v1/user/change/avatar", func(w http.ResponseWriter, r *http.Request) {
		SetHandler(SetAvatarHandler, UserHandler, w, r)
	})

	http.HandleFunc("/api/v1/user/current", func(w http.ResponseWriter, r *http.Request) {
		SetHandler(CurrentProfile, UserHandler, w, r)
	})

	http.Handle("/avatars/", http.StripPrefix("/avatars/", http.FileServer(http.Dir("./avatars"))))

>>>>>>> CP-21: Реализовать механизм загрузки изображений профиля и просмотра профиля
	log.Println("Server listening on ", url+port)
	http.ListenAndServe(url+port, nil)
}
