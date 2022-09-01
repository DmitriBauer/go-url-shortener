package handlers

import (
	"errors"
	"io"
	"net/http"

	"github.com/dmitribauer/go-url-shortener/internal/api/rest"
	"github.com/dmitribauer/go-url-shortener/internal/reqrep"
	"github.com/dmitribauer/go-url-shortener/internal/urlrep"
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

	sessionID := sessionIDFromRequest(rest, w, r)

	var statusCode int

	urlID, err := rest.URLRepo.Save(r.Context(), url, sessionID)
	if err == nil {
		statusCode = http.StatusCreated
	} else if errors.Is(err, urlrep.ErrDuplicateURL) {
		statusCode = http.StatusConflict
	} else {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	shortURL := rest.ShortURL(urlID)

	err = rest.ReqRepo.Save(reqrep.Req{
		SessionID:   sessionID,
		ShortURL:    shortURL,
		OriginalURL: url,
	})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(statusCode)
	w.Write([]byte(shortURL))
}

func handleRootGet(rest *rest.Rest, w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path
	baseLen := len(rest.Path + "/")
	if len(path) <= baseLen {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	id := path[baseLen:]

	url, removed := rest.URLRepo.URLByID(r.Context(), id)
	if url == "" {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	if removed {
		w.WriteHeader(http.StatusGone)
		return
	}

	w.Header().Set("Location", url)
	w.WriteHeader(http.StatusTemporaryRedirect)
}
