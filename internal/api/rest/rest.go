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
	URLRepo    urlrep.URLRepo
	httpServer *http.Server
}

func (rest *Rest) ShortURL(urlID string) string {
	return fmt.Sprintf("http://%s:%d/%s", rest.Address, rest.Port, urlID)
}

func (rest *Rest) Run(address string, port int, prep func(mux *chi.Mux)) error {
	mux := chi.NewMux()
	prep(mux)
	rest.httpServer = &http.Server{
		Addr:    fmt.Sprintf("%s:%d", address, port),
		Handler: mux,
	}
	rest.Address = address
	rest.Port = port
	return rest.httpServer.ListenAndServe()
}
