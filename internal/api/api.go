package api

import (
	"compress/gzip"
	"net/http"

	"github.com/go-chi/chi/v5"
	chimiddleware "github.com/go-chi/chi/v5/middleware"

	"github.com/dmitribauer/go-url-shortener/internal/api/handlers"
	"github.com/dmitribauer/go-url-shortener/internal/api/middleware"
	"github.com/dmitribauer/go-url-shortener/internal/api/rest"
	"github.com/dmitribauer/go-url-shortener/internal/urlrep"
)

func NewRest(urlRepo urlrep.URLRepo) *rest.Rest {
	return &rest.Rest{
		URLRepo: urlRepo,
	}
}

func Run(rest *rest.Rest, address string, port int, path string) error {
	return rest.Run(address, port, path, func(mux *chi.Mux) {
		mux.Use(chimiddleware.NewCompressor(gzip.DefaultCompression).Handler)
		mux.Use(middleware.DecompressGZipHandler)

		mux.Get(rest.Path+"{id}", func(writer http.ResponseWriter, request *http.Request) {
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
