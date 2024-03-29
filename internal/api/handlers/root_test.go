package handlers

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	apirest "github.com/dmitribauer/go-url-shortener/internal/api/rest"
	"github.com/dmitribauer/go-url-shortener/internal/auth"
	"github.com/dmitribauer/go-url-shortener/internal/reqrep"

	"github.com/dmitribauer/go-url-shortener/internal/urlrep"
)

func Test_handleRoot_POST(t *testing.T) {
	type args struct {
		body []byte
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
			name: "POST a correct URL in the body",
			args: args{[]byte("https://yandex.ru")},
			want: want{http.StatusCreated, []byte("http://127.0.0.1:8282/s/uRlId123")},
		},
		{
			name: "POST a wrong URL in the body",
			args: args{[]byte("http//google.com")},
			want: want{http.StatusBadRequest, []byte{}},
		},
		{
			name: "POST with an empty body",
			args: args{[]byte{}},
			want: want{http.StatusBadRequest, []byte{}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rest := newTestRest(nil)
			req := httptest.NewRequest(
				http.MethodPost,
				fmt.Sprintf("http://%s:%d%s", rest.Address, rest.Port, rest.Path),
				bytes.NewReader(tt.args.body),
			)
			w := httptest.NewRecorder()
			HandleRoot(rest, w, req)
			res := w.Result()

			assert.Equal(t, tt.want.code, res.StatusCode)

			defer res.Body.Close()
			body, err := ioutil.ReadAll(res.Body)
			if err != nil {
				assert.Error(t, err)
			}
			assert.Equal(t, string(tt.want.body), string(body))
		})
	}
}

func Test_handleRoot_GET(t *testing.T) {
	id := "ID"
	wrongID := "WRONG_ID"
	url := "https://yandex.ru"
	type args struct {
		id string
	}
	type want struct {
		code    int
		headers map[string]string
	}
	tests := []struct {
		name string
		args args
		want want
	}{
		{
			name: "GET a URL by the correct id",
			args: args{id},
			want: want{http.StatusTemporaryRedirect, map[string]string{"Location": url}},
		},
		{
			name: "GET a URL by the wrong id",
			args: args{wrongID},
			want: want{http.StatusNotFound, map[string]string{}},
		},
		{
			name: "GET a URL by an empty id",
			args: args{""},
			want: want{http.StatusBadRequest, map[string]string{}},
		},
	}

	urlIDGenerator := func(url string) string {
		return id
	}
	urlRepo := urlrep.NewInMemory(urlIDGenerator)
	reqRepo, err := reqrep.NewInFile("/tmp/reqrep/")
	require.NoError(t, err)

	rest := &apirest.Rest{
		Address:     "localhost",
		Port:        8080,
		Path:        "",
		URLRepo:     urlRepo,
		ReqRepo:     reqRepo,
		AuthService: auth.NewService(nil),
	}

	req := httptest.NewRequest(
		http.MethodPost,
		fmt.Sprintf("http://%s:%d", rest.Address, rest.Port),
		bytes.NewReader([]byte(url)),
	)
	w := httptest.NewRecorder()
	HandleRoot(rest, w, req)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(
				http.MethodGet,
				fmt.Sprintf("http://%s:%d/%s", rest.Address, rest.Port, tt.args.id),
				nil,
			)
			w := httptest.NewRecorder()
			HandleRoot(rest, w, req)
			res := w.Result()
			res.Body.Close() // statictest

			assert.Equal(t, tt.want.code, res.StatusCode)

			for k, v := range tt.want.headers {
				assert.Equal(t, v, res.Header.Get(k))
			}
		})
	}
}

func Test_handleRoot_OtherRESTMethods(t *testing.T) {
	type args struct {
		method string
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
			name: "HEAD",
			args: args{http.MethodHead},
			want: want{http.StatusMethodNotAllowed},
		},
		{
			name: "PUT",
			args: args{http.MethodPut},
			want: want{http.StatusMethodNotAllowed},
		},
		{
			name: "DELETE",
			args: args{http.MethodDelete},
			want: want{http.StatusMethodNotAllowed},
		},
		{
			name: "CONNECT",
			args: args{http.MethodConnect},
			want: want{http.StatusMethodNotAllowed},
		},
		{
			name: "OPTIONS",
			args: args{http.MethodOptions},
			want: want{http.StatusMethodNotAllowed},
		},
		{
			name: "TRACE",
			args: args{http.MethodTrace},
			want: want{http.StatusMethodNotAllowed},
		},
		{
			name: "PATCH",
			args: args{http.MethodPatch},
			want: want{http.StatusMethodNotAllowed},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			urlIDGenerator := func(url string) string {
				return "ID"
			}
			urlRepo := urlrep.NewInMemory(urlIDGenerator)
			rest := &apirest.Rest{
				Address: "localhost",
				Port:    8080,
				Path:    "/",
				URLRepo: urlRepo,
			}

			req := httptest.NewRequest(
				tt.args.method,
				fmt.Sprintf("http://%s:%d", rest.Address, rest.Port),
				nil,
			)
			w := httptest.NewRecorder()
			HandleRoot(rest, w, req)
			res := w.Result()
			res.Body.Close() // statictest

			assert.Equal(t, tt.want.code, res.StatusCode)
		})
	}
}
