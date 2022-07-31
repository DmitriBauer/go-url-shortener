package api

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/dmitribauer/go-url-shortener/internal/urlrep"
)

func TestAPI_NewRest(t *testing.T) {
	urlRepo := urlrep.URLRepo(nil)
	rest := NewRest(urlRepo)

	assert.Equal(t, urlRepo, rest.URLRepo)
}
