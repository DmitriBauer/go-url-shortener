package urlrep

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewInMemory(t *testing.T) {
	r := NewInMemory(nil)

	assert.Equal(t, reflect.TypeOf(&inMemURLRepo{}), reflect.TypeOf(r))
}

func TestInMemURLRepo_URLByID(t *testing.T) {
	id := "ID"
	url := "https://www.yandex.ru"
	urlIDGenerator := func(url string) string {
		return id
	}
	r := &inMemURLRepo{
		urlIDGenerator: urlIDGenerator,
	}

	u, ok := r.URLByID(id)
	require.Equal(t, "", u)
	require.Equal(t, false, ok)

	r.Save(url)

	u, ok = r.URLByID(id)
	assert.Equal(t, url, u)
	assert.Equal(t, true, ok)
}

func TestInMemURLRepo_Save(t *testing.T) {
	id := "ID"
	url := "https://www.yandex.ru"
	urlIDGenerator := func(url string) string {
		return id
	}
	r := &inMemURLRepo{
		urlIDGenerator: urlIDGenerator,
	}

	u, ok := r.URLByID(id)
	require.Equal(t, "", u)
	require.Equal(t, false, ok)

	urlID := r.Save(url)

	u, ok = r.URLByID(id)
	assert.Equal(t, id, urlID)
	assert.Equal(t, url, u)
	assert.Equal(t, true, ok)
}

func TestInMemURLRepo_GenerateID(t *testing.T) {
	id := "GenID"
	urlIDGenerator := func(url string) string {
		return id
	}

	r := &inMemURLRepo{
		urlIDGenerator: urlIDGenerator,
	}

	assert.Equal(t, id, r.GenerateID("https://www.yandex.ru"))
}
