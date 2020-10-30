package builder

import (
	"net/http"

	"github.com/go-park-mail-ru/2020_2_CodeExpress/internal/models"
)

func BuildCookie(session *models.Session) *http.Cookie {
	return &http.Cookie{
		Value:    session.ID,
		Name:     session.Name,
		Expires:  session.Expire,
		HttpOnly: true,
		Path:     "/",
		SameSite: http.SameSiteStrictMode,
	}
}
