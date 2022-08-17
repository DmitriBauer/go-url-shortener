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

func Test_HandleShortenBatchPost(t *testing.T) {
	type args struct {
		headers map[string]string
		body    []byte
	}
	type want struct {
		code int
		body []byte
	}
	tests := []struct {
		name string
		args args
		want want
	}{
		{
			name: "headers do not contain Content-Type: application/json",
			want: want{
				code: http.StatusUnsupportedMediaType,
			},
		},
		{
			name: "headers contain Content-Type: application/json and body is empty",
			args: args{
				headers: map[string]string{
					"Content-Type": "application/json",
				},
			},
			want: want{
				code: http.StatusBadRequest,
			},
		},
		{
			name: "headers contain Content-Type: application/json and body is not empty",
			args: args{
				headers: map[string]string{
					"Content-Type": "application/json",
				},
				body: []byte(`[{"correlation_id": "123", "original_url": "https://yandex.ru"}]`),
			},
			want: want{
				code: http.StatusCreated,
				body: []byte(`[{"correlation_id": "123", "short_url": "http://127.0.0.1:8282/s/uRlId123"}]`),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rest := newTestRest(nil)
			req := httptest.NewRequest(
				http.MethodPost,
				fmt.Sprintf("http://%s:%d%s/api/shorten/batch", rest.Address, rest.Port, rest.Path),
				bytes.NewReader(tt.args.body),
			)
			for k, v := range tt.args.headers {
				req.Header.Set(k, v)
			}
			w := httptest.NewRecorder()

			HandleShortenBatchPost(rest, w, req)

			assert.Equal(t, tt.want.code, w.Code)

			if tt.want.body != nil {
				var wantEntries []batchResBodyEntry
				json.Unmarshal(tt.want.body, &wantEntries)
				var gotEntries []batchResBodyEntry
				json.Unmarshal(w.Body.Bytes(), &gotEntries)

				assert.Equal(t, wantEntries, gotEntries)
			}
		})
	}
}
