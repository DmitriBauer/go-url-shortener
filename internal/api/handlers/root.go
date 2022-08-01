package handlers

import (
	"io"
	"net/http"

	"github.com/dmitribauer/go-url-shortener/internal/api/rest"
	"github.com/dmitribauer/go-url-shortener/internal/util"
)

func HandleRoot(rest *rest.Rest, w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		handleRootPost(rest, w, r)
	case http.MethodGet:
		handleRootGet(rest, w, r)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func handleRootPost(rest *rest.Rest, w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	body, err := io.ReadAll(r.Body)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	url := string(body)
	if url == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if !util.CheckIsURL(url) {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	urlID := rest.URLRepo.Save(url)
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(rest.ShortURL(urlID)))
}

func handleRootGet(rest *rest.Rest, w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path
	baseLen := len(rest.Path)
	if len(path) <= baseLen {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	id := path[baseLen:]

	url, ok := rest.URLRepo.URLByID(id)
	if !ok {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	w.Header().Set("Location", url)
	w.WriteHeader(http.StatusTemporaryRedirect)
}
