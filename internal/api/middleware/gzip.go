package middleware

import (
	"compress/gzip"
	"net/http"
)

func DecompressGZipHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Content-Encoding") == "gzip" {
			gr, err := gzip.NewReader(r.Body)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			defer gr.Close()
			r.Body = gr
		}
		next.ServeHTTP(w, r)
	})
}
