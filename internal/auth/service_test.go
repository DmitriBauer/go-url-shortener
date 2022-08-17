package auth

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewService(t *testing.T) {
	sessionIDGenerator := func() string {
		return "41734162-d37e-4936-ab4c-b808386a34c9"
	}

	s := NewService(sessionIDGenerator)
	assert.Equal(t, sessionIDGenerator(), s.sessionIDGenerator())
}
