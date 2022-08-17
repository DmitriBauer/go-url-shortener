package api

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/dmitribauer/go-url-shortener/internal/auth"
	"github.com/dmitribauer/go-url-shortener/internal/reqrep"
	"github.com/dmitribauer/go-url-shortener/internal/urlrep"
)

func TestAPI_NewRest(t *testing.T) {
	urlRepo := urlrep.URLRepo(nil)
	reqRepo, err := reqrep.NewInFile("/tmp/reqrep/")
	require.NoError(t, err)
	authService := auth.NewService(nil)

	rest := NewRest(
		urlRepo,
		reqRepo,
		authService,
	)

	assert.Equal(t, urlRepo, rest.URLRepo)
	assert.Equal(t, reqRepo, rest.ReqRepo)
	assert.Equal(t, authService, rest.AuthService)
}
