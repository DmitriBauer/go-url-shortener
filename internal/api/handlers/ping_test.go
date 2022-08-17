package handlers

import (
	"bytes"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_HandlePingGet(t *testing.T) {
	rest := newTestRest(nil)
	req := httptest.NewRequest(
		http.MethodGet,
		fmt.Sprintf("http://%s:%d%s/ping", rest.Address, rest.Port, rest.Path),
		bytes.NewReader([]byte{}),
	)
	w := httptest.NewRecorder()

	HandlePingGet(rest, w, req)

	res := w.Result()
	res.Body.Close() // statictest

	assert.Equal(t, http.StatusOK, res.StatusCode)
}
