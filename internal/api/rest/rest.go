package rest

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/dmitribauer/go-url-shortener/internal/urlrep"
)

type Rest struct {
	Address    string
	Port       int
	Path       string
	URLRepo    urlrep.URLRepo
	httpServer *http.Server
}

func (rest *Rest) ShortURL(urlID string) string {
	return fmt.Sprintf("http://%s:%d%s%s", rest.Address, rest.Port, rest.Path, urlID)
}

func (rest *Rest) Run(address string, port int, path string, prep func(mux *chi.Mux)) error {
	rest.Address = address
	rest.Port = port
	rest.Path = path
	mux := chi.NewMux()
	prep(mux)
	rest.httpServer = &http.Server{
		Addr:    fmt.Sprintf("%s:%d", address, port),
		Handler: mux,
	}
	return rest.httpServer.ListenAndServe()
}
