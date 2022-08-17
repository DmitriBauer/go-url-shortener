package handlers

import (
	"bytes"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_sessionIDFromRequest(t *testing.T) {
	tests := []struct {
		name      string
		sessionID string
		want      string
	}{
		{
			name:      "session_id cookie is set",
			sessionID: "2e0303c1-4175-12bf-86ca-39d2d1f6ae88",
			want:      "2e0303c1-4175-12bf-86ca-39d2d1f6ae88",
		},
		{
			name:      "session_id cookie is not set",
			sessionID: "",
			want:      "41734162-d37e-4936-ab4c-b808386a34c9",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rest := newTestRest(nil)
			req := httptest.NewRequest(
				http.MethodGet,
				fmt.Sprintf("http://%s:%d%s", rest.Address, rest.Port, rest.Path),
				bytes.NewReader([]byte{}),
			)
			w := httptest.NewRecorder()
			if tt.sessionID != "" {
				req.AddCookie(&http.Cookie{
					Name:  "session_id",
					Value: tt.sessionID,
				})
			}

			got := sessionIDFromRequest(rest, w, req)

			assert.Equal(t, tt.want, got)
		})
	}
}
