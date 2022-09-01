package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"

	apirest "github.com/dmitribauer/go-url-shortener/internal/api/rest"
	"github.com/dmitribauer/go-url-shortener/internal/auth"
	"github.com/dmitribauer/go-url-shortener/internal/reqrep"
)

func TestHandleURLsGet(t *testing.T) {
	sessionID := "41734162-d37e-4936-ab4c-b808386a34c9"
	restRepoDir := "/tmp/testreqrep/"

	record := reqrep.ReqRecord{
		ShortURL:    "http://127.0.0.1:8282/s/uRlId123",
		OriginalURL: "https://yandex.ru",
	}
	records, _ := json.Marshal([]reqrep.ReqRecord{record})

	authService := auth.NewService(func() string {
		return sessionID
	})
	reqRepo, _ := reqrep.NewInFile(restRepoDir)
	reqRepo.Save(reqrep.Req{
		SessionID:   sessionID,
		ShortURL:    record.ShortURL,
		OriginalURL: record.OriginalURL,
	})
	defer os.Remove(restRepoDir + sessionID)

	type args struct {
		sessionID string
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
			name: "sessionID is new",
			args: args{
				sessionID: "2e0303c1-4175-12bf-86ca-39d2d1f6ae88",
			},
			want: want{
				code: http.StatusNoContent,
			},
		},
		{
			name: "sessionID is not new",
			args: args{
				sessionID: sessionID,
			},
			want: want{
				code: http.StatusOK,
				body: records,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rest := newTestRest(&apirest.Rest{
				AuthService: authService,
				ReqRepo:     reqRepo,
			})
			req := httptest.NewRequest(
				http.MethodGet,
				fmt.Sprintf("http://%s:%d%s/api/user/urls", rest.Address, rest.Port, rest.Path),
				bytes.NewReader([]byte{}),
			)
			if tt.args.sessionID != "" {
				req.AddCookie(&http.Cookie{
					Name:  "session_id",
					Value: tt.args.sessionID,
				})
			}
			w := httptest.NewRecorder()

			HandleURLsGet(rest, w, req)

			res := w.Result()
			body, _ := ioutil.ReadAll(res.Body)
			defer res.Body.Close()

			assert.Equal(t, tt.want.code, res.StatusCode)

			if tt.want.body != nil {
				var wantRecords []reqrep.ReqRecord
				json.Unmarshal(tt.want.body, &wantRecords)
				var gotRecords []reqrep.ReqRecord
				json.Unmarshal(body, &gotRecords)

				assert.Equal(t, wantRecords, gotRecords)
			}
		})
	}
}

func TestHandleURLsDelete(t *testing.T) {
	type args struct {
		body []byte
	}
	type want struct {
		code int
	}
	tests := []struct {
		name string
		args args
		want want
	}{
		{
			name: "body is empty",
			want: want{
				code: http.StatusBadRequest,
			},
		},
		{
			name: "body is not empty",
			args: args{
				body: []byte(`["f85516fe", "a23eb735"]`),
			},
			want: want{
				code: http.StatusAccepted,
			},
		},
		{
			name: "body is not json",
			args: args{
				body: []byte(`["f85516fe", "a23eb`),
			},
			want: want{
				code: http.StatusBadRequest,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rest := newTestRest(&apirest.Rest{})
			req := httptest.NewRequest(
				http.MethodDelete,
				fmt.Sprintf("http://%s:%d%s/api/user/urls", rest.Address, rest.Port, rest.Path),
				bytes.NewReader(tt.args.body),
			)
			w := httptest.NewRecorder()

			HandleURLsDelete(rest, w, req)

			res := w.Result()
			defer res.Body.Close() // statictest

			assert.Equal(t, tt.want.code, res.StatusCode)
		})
	}
}
