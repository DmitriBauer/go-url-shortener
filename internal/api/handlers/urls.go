package handlers

import (
	"context"
	"encoding/json"
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

func HandleURLsDelete(rest *rest.Rest, w http.ResponseWriter, r *http.Request) {
	sessionID := sessionIDFromRequest(rest, w, r)

	defer r.Body.Close()
	var ids []string
	if json.NewDecoder(r.Body).Decode(&ids) != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	go rest.URLRepo.RemoveList(context.Background(), ids, sessionID)

	w.WriteHeader(http.StatusAccepted)
}
