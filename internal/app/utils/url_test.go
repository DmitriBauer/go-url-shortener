package utils

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_CheckIsURL(t *testing.T) {
	tests := []struct {
		name string
		url  string
		want bool
	}{
		{
			name: "Check the URL with a wrong format",
			url:  "http//yandex.ru",
			want: false,
		},
		{
			name: "Check the URL without a scheme",
			url:  "yandex.ru",
			want: false,
		},
		{
			name: "Check the URL consists only of the path",
			url:  "/pogoda/moscow",
			want: false,
		},
		{
			name: "Check the URL consists only of the scheme ",
			url:  "https://",
			want: false,
		},
		{
			name: "Check the URL with a wrong scheme ",
			url:  "htt://yandex.ru",
			want: false,
		},
		{
			name: "Check the correct URL with http",
			url:  "http://yandex.ru",
			want: true,
		},
		{
			name: "Check the correct URL with https",
			url:  "https://yandex.ru",
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, CheckIsURL(tt.url))
		})
	}
}
