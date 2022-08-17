package handlers

import (
	"net/http"

	"github.com/dmitribauer/go-url-shortener/internal/api/rest"
)

func HandleURLsGet(rest *rest.Rest, w http.ResponseWriter, r *http.Request) {
	sessionID := sessionIDFromRequest(rest, w, r)

	data, err := rest.ReqRepo.DataBySessionID(sessionID)
	if err != nil {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}
