package urlrep

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"reflect"
	"testing"
)

func TestNewInMemory(t *testing.T) {
	r := NewInMemory(nil)

	assert.Equal(t, reflect.TypeOf(&inMemURLRepo{}), reflect.TypeOf(r))
}

func TestInMemURLRepo_URLByID(t *testing.T) {
	id := "ID"
	urlIDGenerator := func(url string) string {
		return id
	}
	r := &inMemURLRepo{
		urlIDGenerator: urlIDGenerator,
	}

	require.Equal(t, "", r.URLByID(id))

	r.Save("https://www.yandex.ru")

	assert.Equal(t, "https://www.yandex.ru", r.URLByID(id))
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

	require.Equal(t, "", r.URLByID(id))

	urlID := r.Save(url)

	assert.Equal(t, id, urlID)
	assert.Equal(t, url, r.URLByID(urlID))
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
