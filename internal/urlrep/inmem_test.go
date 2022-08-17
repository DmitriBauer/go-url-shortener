package urlrep

import (
	"context"
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

	u, ok := r.URLByID(context.TODO(), id)
	require.Equal(t, "", u)
	require.Equal(t, false, ok)

	r.Save(context.TODO(), url)

	u, ok = r.URLByID(context.TODO(), id)
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

	u, ok := r.URLByID(context.TODO(), id)
	require.Equal(t, "", u)
	require.Equal(t, false, ok)

	urlID, _ := r.Save(context.TODO(), url)

	u, ok = r.URLByID(context.TODO(), id)
	assert.Equal(t, id, urlID)
	assert.Equal(t, url, u)
	assert.Equal(t, true, ok)
}

func TestInMemURLRepo_SaveList(t *testing.T) {
	ctx := context.TODO()
	urls := []string{
		"https://www.yandex.ru",
		"https://www.google.com",
	}
	urlIDGenerator := func(url string) string {
		switch url {
		case "https://www.yandex.ru":
			return "yaID"
		case "https://www.google.com":
			return "googleID"
		default:
			return "ID"
		}
	}
	r := &inMemURLRepo{
		urlIDGenerator: urlIDGenerator,
	}

	ids, err := r.SaveList(ctx, urls)
	require.NoError(t, err)
	assert.Equal(t, []string{"yaID", "googleID"}, ids)

	u, ok := r.URLByID(ctx, "yaID")
	assert.Equal(t, "https://www.yandex.ru", u)
	assert.Equal(t, true, ok)

	u, ok = r.URLByID(ctx, "googleID")
	assert.Equal(t, "https://www.google.com", u)
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
