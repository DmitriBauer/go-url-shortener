package api

import (
	"compress/gzip"
	"net/http"

	"github.com/go-chi/chi/v5"
	chimiddleware "github.com/go-chi/chi/v5/middleware"

	"github.com/dmitribauer/go-url-shortener/internal/api/handlers"
	"github.com/dmitribauer/go-url-shortener/internal/api/middleware"
	"github.com/dmitribauer/go-url-shortener/internal/api/rest"
	"github.com/dmitribauer/go-url-shortener/internal/auth"
	"github.com/dmitribauer/go-url-shortener/internal/reqrep"
	"github.com/dmitribauer/go-url-shortener/internal/urlrep"
)

func NewRest(
	urlRepo urlrep.URLRepo,
	reqRepo reqrep.ReqRepo,
	authService *auth.Service,
) *rest.Rest {
	return &rest.Rest{
		URLRepo:     urlRepo,
		ReqRepo:     reqRepo,
		AuthService: authService,
	}
}

func Run(rest *rest.Rest, address string, port int, path string) error {
	return rest.Run(address, port, path, func(mux *chi.Mux) {
		mux.Use(chimiddleware.NewCompressor(gzip.DefaultCompression).Handler)
		mux.Use(middleware.DecompressGZipHandler)

		mux.Get(rest.Path+"/{id}", func(writer http.ResponseWriter, request *http.Request) {
			handlers.HandleRoot(rest, writer, request)
		})
		mux.Get(rest.Path+"/ping", func(writer http.ResponseWriter, request *http.Request) {
			handlers.HandlePingGet(rest, writer, request)
		})
		mux.Get(rest.Path+"/api/user/urls", func(writer http.ResponseWriter, request *http.Request) {
			handlers.HandleURLsGet(rest, writer, request)
		})

		mux.Post(rest.Path+"/", func(writer http.ResponseWriter, request *http.Request) {
			handlers.HandleRoot(rest, writer, request)
		})
		mux.Post(rest.Path+"/api/shorten", func(writer http.ResponseWriter, request *http.Request) {
			handlers.HandleShortenPost(rest, writer, request)
		})
		mux.Post(rest.Path+"/api/shorten/batch", func(writer http.ResponseWriter, request *http.Request) {
			handlers.HandleShortenBatchPost(rest, writer, request)
		})

		mux.Delete(rest.Path+"/api/user/urls", func(writer http.ResponseWriter, request *http.Request) {
			handlers.HandleURLsDelete(rest, writer, request)
		})
	})
}
