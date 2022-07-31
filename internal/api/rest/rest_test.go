package rest

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRest_shortURL(t *testing.T) {
	urlID := "ID"
	address := "localhost"
	port := 8080
	rest := Rest{
		Address:    address,
		Port:       port,
		URLRepo:    nil,
		httpServer: nil,
	}

	assert.Equal(t, fmt.Sprintf("http://%s:%d/%s", address, port, urlID), rest.ShortURL(urlID))
}
