package api

import (
	"fmt"
	"github.com/dmitribauer/go-url-shortener/internal/app/urlrep"
	"github.com/dmitribauer/go-url-shortener/internal/app/utils"
	"github.com/go-chi/chi/v5"
	"io"
	"net/http"
)

type Rest struct {
	address    string
	port       int
	urlRepo    urlrep.URLRepo
	httpServer *http.Server
}

func NewRest(urlRepo urlrep.URLRepo) *Rest {
	return &Rest{
		address: "localhost",
		port:    8080,
		urlRepo: urlRepo,
	}
}

func (rest *Rest) Run(address string, port int) error {
	mux := chi.NewMux()
	mux.Route("/", func(r chi.Router) {
		r.Post("/", rest.handleRoot)
		r.Get("/{id}", rest.handleRoot)
	})
	rest.httpServer = &http.Server{
		Addr:    fmt.Sprintf("%s:%d", address, port),
		Handler: mux,
	}
	rest.address = address
	rest.port = port
	return rest.httpServer.ListenAndServe()
}

func (rest *Rest) handleRoot(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		rest.handleRootPost(w, r)
	case http.MethodGet:
		rest.handleRootGet(w, r)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func (rest *Rest) handleRootPost(w http.ResponseWriter, r *http.Request) {
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

	if !utils.CheckIsURL(url) {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	id := rest.urlRepo.Save(url)
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(fmt.Sprintf("http://%s:%d/%s", rest.address, rest.port, id)))
}

func (rest *Rest) handleRootGet(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path
	if len(path) < 2 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	id := path[1:]

	url := rest.urlRepo.URLByID(id)
	if url == "" {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	w.Header().Set("Location", url)
	w.WriteHeader(http.StatusTemporaryRedirect)
}
