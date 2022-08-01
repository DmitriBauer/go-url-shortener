package urlrep

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInFileSplitter_scanLinesInReverse(t *testing.T) {
	data := []byte(
		`a661b08 https://youtube.com
af6fae6 https://yandex.ru
wersds2 https://google.com
d23sdfs https://twitter.com
17a28625 https://youtube.com
`)

	lines := []string{
		"17a28625 https://youtube.com",
		"d23sdfs https://twitter.com",
		"wersds2 https://google.com",
		"af6fae6 https://yandex.ru",
		"a661b08 https://youtube.com",
	}

	s := newInFileSplitter()

	for _, line := range lines {
		_, token, err := s.scanLinesInReverse(data, false)
		assert.Equal(t, line, string(token))
		assert.NoError(t, err)
	}
}

func TestInFileSplitter_scanLinesInReverseEmptyData(t *testing.T) {
	data := []byte("")

	s := newInFileSplitter()

	d, token, err := s.scanLinesInReverse(data, false)

	assert.Equal(t, 0, d)
	assert.Equal(t, []byte(nil), token)
	assert.NoError(t, err)
}
