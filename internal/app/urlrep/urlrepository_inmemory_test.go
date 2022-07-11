package urlrep

import (
	"github.com/stretchr/testify/require"
	"reflect"
	"testing"
)
import "github.com/stretchr/testify/assert"

func TestNewInMemory(t *testing.T) {
	r := NewInMemory(nil)

	assert.Equal(t, reflect.TypeOf(&inMemoryURLRepository{}), reflect.TypeOf(r))
}

func TestInMemoryURLRepository_Get(t *testing.T) {
	id := "ID"
	urlIDGenerator := func(url string) string {
		return id
	}
	r := &inMemoryURLRepository{
		urlIDGenerator: urlIDGenerator,
	}

	require.Equal(t, "", r.Get(id))

	r.Set("https://www.yandex.ru")

	assert.Equal(t, "https://www.yandex.ru", r.Get(id))
}

func TestInMemoryURLRepository_Set(t *testing.T) {
	id := "ID"
	urlIDGenerator := func(url string) string {
		return id
	}
	r := &inMemoryURLRepository{
		urlIDGenerator: urlIDGenerator,
	}

	require.Equal(t, "", r.Get(id))

	urlID := r.Set("https://www.yandex.ru")

	assert.Equal(t, id, urlID)
	assert.Equal(t, "https://www.yandex.ru", r.Get(urlID))
}

func TestInMemoryURLRepository_GenerateID(t *testing.T) {
	id := "GenID"
	urlIDGenerator := func(url string) string {
		return id
	}

	r := &inMemoryURLRepository{
		urlIDGenerator: urlIDGenerator,
	}

	assert.Equal(t, id, r.GenerateID("https://www.yandex.ru"))
}
