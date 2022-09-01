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
	defaultID, yaID, googleID := "ID", "yaID", "googleID"
	yaURL, googleURL := "https://www.yandex.ru", "https://www.google.com"
	sessionID := "1bd118e2-a59b-4110-b315-8e191b16e41c"
	ctx := context.TODO()

	urlIDGenerator := func(url string) string {
		switch url {
		case yaURL:
			return yaID
		case googleURL:
			return googleID
		default:
			return defaultID
		}
	}
	r := &inMemURLRepo{
		urlIDGenerator: urlIDGenerator,
	}

	u, _ := r.URLByID(ctx, yaID)
	require.Equal(t, "", u)

	u, _ = r.URLByID(ctx, googleID)
	require.Equal(t, "", u)

	r.Save(ctx, yaURL, sessionID)
	r.Save(ctx, googleURL, sessionID)
	r.RemoveList(ctx, []string{yaID}, sessionID)

	u, removed := r.URLByID(ctx, yaID)
	assert.Equal(t, yaURL, u)
	assert.Equal(t, true, removed)

	u, removed = r.URLByID(ctx, googleID)
	assert.Equal(t, googleURL, u)
	assert.Equal(t, false, removed)
}

func TestInMemURLRepo_Save(t *testing.T) {
	id := "ID"
	url := "https://www.yandex.ru"
	sessionID := "1bd118e2-a59b-4110-b315-8e191b16e41c"
	urlIDGenerator := func(url string) string {
		return id
	}
	r := &inMemURLRepo{
		urlIDGenerator: urlIDGenerator,
	}

	u, _ := r.URLByID(context.TODO(), id)
	require.Equal(t, "", u)

	urlID, _ := r.Save(context.TODO(), url, sessionID)

	u, removed := r.URLByID(context.TODO(), id)
	assert.Equal(t, id, urlID)
	assert.Equal(t, url, u)
	assert.Equal(t, false, removed)
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

	ids, err := r.SaveList(ctx, urls, "1bd118e2-a59b-4110-b315-8e191b16e41c")
	require.NoError(t, err)
	assert.Equal(t, []string{"yaID", "googleID"}, ids)

	u, removed := r.URLByID(ctx, "yaID")
	assert.Equal(t, "https://www.yandex.ru", u)
	assert.Equal(t, false, removed)

	u, removed = r.URLByID(ctx, "googleID")
	assert.Equal(t, "https://www.google.com", u)
	assert.Equal(t, false, removed)

}

func TestInMemURLRepo_RemoveList(t *testing.T) {
	id := "ID"
	url := "https://www.yandex.ru"
	sessionID := "1bd118e2-a59b-4110-b315-8e191b16e41c"
	ctx := context.TODO()

	urlIDGenerator := func(url string) string {
		return id
	}
	r := &inMemURLRepo{
		urlIDGenerator: urlIDGenerator,
	}

	r.Save(ctx, url, sessionID)

	u, removed := r.URLByID(ctx, id)
	require.Equal(t, url, u)
	require.Equal(t, false, removed)

	r.RemoveList(ctx, []string{id}, sessionID)

	u, removed = r.URLByID(ctx, id)
	assert.Equal(t, url, u)
	assert.Equal(t, true, removed)
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
