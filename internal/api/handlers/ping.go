package handlers

import (
	"net/http"

	"github.com/dmitribauer/go-url-shortener/internal/api/rest"
)

func HandlePingGet(rest *rest.Rest, w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}
