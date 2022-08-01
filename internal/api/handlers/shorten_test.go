package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"

	apirest "github.com/dmitribauer/go-url-shortener/internal/api/rest"
	"github.com/dmitribauer/go-url-shortener/internal/urlrep"
)

func Test_HandleShorten_POST(t *testing.T) {
	urlID := "ID"
	type args struct {
		contentType string
		body        string
	}
	type want struct {
		code int
		body shortenResBody
	}
	tests := []struct {
		name string
		args args
		want want
	}{
		{
			name: "POST a correct URL in the body and the application/json as content type",
			args: args{"application/json", `{"url": "https://yandex.ru"}`},
			want: want{http.StatusCreated, shortenResBody{Result: fmt.Sprintf("http://localhost:8080/%s", urlID)}},
		},
		{
			name: "POST a correct URL in the body and a wrong content-type (text/plain)",
			args: args{"text/plain", `{"url": "https://yandex.ru"}`},
			want: want{http.StatusUnsupportedMediaType, shortenResBody{}},
		},
		{
			name: "POST a wrong URL in the body and the application/json as content type",
			args: args{"application/json", `{"url": "httpss://yandex.ru"}`},
			want: want{http.StatusBadRequest, shortenResBody{}},
		},
		{
			name: "POST a wrong body and the application/json as content type",
			args: args{"application/json", `{}`},
			want: want{http.StatusBadRequest, shortenResBody{}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			urlIDGenerator := func(url string) string {
				return urlID
			}
			urlRepo := urlrep.NewInMemory(urlIDGenerator)
			rest := &apirest.Rest{
				Address: "localhost",
				Port:    8080,
				Path:    "/",
				URLRepo: urlRepo,
			}

			req := httptest.NewRequest(
				http.MethodPost,
				fmt.Sprintf("http://%s:%d", rest.Address, rest.Port),
				bytes.NewReader(json.RawMessage(tt.args.body)),
			)
			req.Header.Set("Content-Type", tt.args.contentType)
			w := httptest.NewRecorder()

			HandleShortenPost(rest, w, req)

			res := w.Result()

			assert.Equal(t, tt.want.code, res.StatusCode)

			defer res.Body.Close()
			var body shortenResBody
			err := json.NewDecoder(res.Body).Decode(&body)
			if err != nil {
				assert.Error(t, err)
			}

			assert.Equal(t, tt.want.body, body)
		})
	}
}
