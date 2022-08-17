package reqrep

import (
	"encoding/json"
	"os"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewInFile(t *testing.T) {
	dir := "/tmp/TestNewInFileReqRep/"
	defer os.RemoveAll(dir)

	r, err := NewInFile(dir)
	require.NoError(t, err)
	assert.Equal(t, reflect.TypeOf(&inFileReqRepo{}), reflect.TypeOf(r))
}

func TestInFileReqRepo_DataByEncryptedUserID_Save(t *testing.T) {
	req := Req{
		SessionID:   "41734162-d37e-4936-ab4c-b808386a34c9",
		ShortURL:    "http://localhost:8080/sdfsaa",
		OriginalURL: "https://yandex.ru",
	}
	entries := []ReqRecord{
		{
			ShortURL:    req.ShortURL,
			OriginalURL: req.OriginalURL,
		},
	}
	data, _ := json.Marshal(entries)

	dir := "/tmp/"
	defer os.Remove(dir + req.SessionID)

	r, err := NewInFile(dir)
	require.NoError(t, err)

	d, err := r.DataBySessionID(req.SessionID)
	require.Error(t, err)
	require.Nil(t, d)

	err = r.Save(req)
	require.NoError(t, err)

	d, err = r.DataBySessionID(req.SessionID)
	require.NoError(t, err)
	assert.Equal(t, data, d)
}
