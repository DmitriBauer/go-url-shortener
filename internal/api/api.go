package api

import (
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/dmitribauer/go-url-shortener/internal/api/handlers"
	"github.com/dmitribauer/go-url-shortener/internal/api/rest"
	"github.com/dmitribauer/go-url-shortener/internal/urlrep"
)

func NewRest(urlRepo urlrep.URLRepo) *rest.Rest {
	return &rest.Rest{
		URLRepo: urlRepo,
	}
}

func Run(rest *rest.Rest, address string, port int) error {
	return rest.Run(address, port, func(mux *chi.Mux) {
		mux.Get("/{id}", func(writer http.ResponseWriter, request *http.Request) {
			handlers.HandleRoot(rest, writer, request)
		})
		mux.Post("/", func(writer http.ResponseWriter, request *http.Request) {
			handlers.HandleRoot(rest, writer, request)
		})
		mux.Post("/api/shorten", func(writer http.ResponseWriter, request *http.Request) {
			handlers.HandleShortenPost(rest, writer, request)
		})
	})
}
