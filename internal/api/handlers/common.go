package handlers

import (
	"net/http"

	"github.com/dmitribauer/go-url-shortener/internal/api/rest"
)

func sessionIDFromRequest(rest *rest.Rest, w http.ResponseWriter, r *http.Request) string {
	name := "session_id"
	cookie, err := r.Cookie(name)
	if err != nil {
		sessionID := rest.AuthService.NewSessionID()
		http.SetCookie(w, &http.Cookie{
			Name:  name,
			Value: sessionID,
		})
		return sessionID
	}
	return cookie.Value
}
