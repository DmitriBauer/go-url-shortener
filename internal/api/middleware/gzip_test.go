package middleware

import (
	"bytes"
	"compress/gzip"
	"io"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDecompressGZipHandler(t *testing.T) {
	m := "Hello, World!"

	var b bytes.Buffer
	gw := gzip.NewWriter(&b)
	gw.Write([]byte(m))
	gw.Close()

	r := http.Request{
		Body: io.NopCloser(bytes.NewReader(b.Bytes())),
		Header: map[string][]string{
			"Content-Encoding": {"gzip"},
		},
	}

	check := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		b, _ := io.ReadAll(r.Body)
		r.Body.Close()

		assert.Equal(t, m, string(b))
	})
	decompress := DecompressGZipHandler(check)
	decompress.ServeHTTP(nil, &r)
}
