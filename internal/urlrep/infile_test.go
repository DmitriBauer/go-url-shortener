package urlrep

import (
	"context"
	"os"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewInFile(t *testing.T) {
	wd, err := os.Getwd()
	require.NoError(t, err)
	path := wd + "/TestNewInFile"
	defer os.Remove(path)

	r, _ := NewInFile(path, nil)

	assert.Equal(t, reflect.TypeOf(&inFileURLRepo{}), reflect.TypeOf(r))
}

func TestInFileURLRepo_URLByID(t *testing.T) {
	wd, err := os.Getwd()
	require.NoError(t, err)
	path := wd + "/TestInFileURLRepo_URLByID"
	defer os.Remove(path)

	id := "ID"
	url := "https://www.yandex.ru"

	urlIDGenerator := func(url string) string {
		return id
	}
	r := &inFileURLRepo{
		path:           path,
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

func TestInFileURLRepo_Save(t *testing.T) {
	wd, err := os.Getwd()
	require.NoError(t, err)
	path := wd + "/TestInFileURLRepo_Save"
	defer os.Remove(path)

	id := "ID"
	url := "https://www.yandex.ru"

	urlIDGenerator := func(url string) string {
		return id
	}
	r := &inFileURLRepo{
		path:           path,
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

func TestInFileURLRepo_SaveList(t *testing.T) {
	wd, err := os.Getwd()
	require.NoError(t, err)
	path := wd + "/TestInFileURLRepo_SaveList"
	defer os.Remove(path)

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
	r := &inFileURLRepo{
		path:           path,
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

func TestInFileURLRepo_GenerateID(t *testing.T) {
	id := "GenID"
	urlIDGenerator := func(url string) string {
		return id
	}

	r := &inFileURLRepo{
		urlIDGenerator: urlIDGenerator,
	}

	assert.Equal(t, id, r.GenerateID("https://www.yandex.ru"))
}
