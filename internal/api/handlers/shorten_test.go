package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_HandleShorten_POST(t *testing.T) {
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
			want: want{http.StatusCreated, shortenResBody{Result: "http://127.0.0.1:8282/s/uRlId123"}},
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
			rest := newTestRest(nil)

			req := httptest.NewRequest(
				http.MethodPost,
				fmt.Sprintf("http://%s:%d%s/api/shorten", rest.Address, rest.Port, rest.Path),
				bytes.NewReader(json.RawMessage(tt.args.body)),
			)
			req.Header.Set("Content-Type", tt.args.contentType)
			w := httptest.NewRecorder()

			HandleShortenPost(rest, w, req)

			res := w.Result()

			assert.Equal(t, tt.want.code, res.StatusCode)

			defer res.Body.Close()
			var body shortenResBody
			json.NewDecoder(res.Body).Decode(&body)

			assert.Equal(t, tt.want.body, body)
		})
	}
}
