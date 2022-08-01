package conf

import (
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConfig_Load(t *testing.T) {
	type env struct {
		name  string
		value string
	}
	type want struct {
		config Config
		err    error
	}
	tests := []struct {
		name string
		envs []env
		want want
	}{
		{
			name: "No envs",
			envs: []env{},
			want: want{
				config: Config{
					ServerURL: "localhost:8080",
					BaseURL:   "http://localhost:8080",
					Address:   "localhost",
					Port:      8080,
					Path:      "/",
				},
				err: nil,
			},
		},
		{
			name: "Correct envs",
			envs: []env{
				{name: "SERVER_ADDRESS", value: "127.0.0.1:8080"},
				{name: "BASE_URL", value: "http://127.0.0.1:8080/short"},
			},
			want: want{
				config: Config{
					ServerURL: "127.0.0.1:8080",
					BaseURL:   "http://127.0.0.1:8080/short",
					Address:   "127.0.0.1",
					Port:      8080,
					Path:      "/short/",
				},
				err: nil,
			},
		},
		{
			name: "Wrong SERVER_ADDRESS",
			envs: []env{
				{name: "SERVER_ADDRESS", value: "127.0.0.1:8080:433"},
				{name: "BASE_URL", value: "http://127.0.0.1:8080/short"},
			},
			want: want{
				config: Config{
					ServerURL: "127.0.0.1:8080:433",
					BaseURL:   "http://127.0.0.1:8080/short",
					Address:   "",
					Port:      0,
					Path:      "",
				},
				err: fmt.Errorf("invalid SERVER_ADDRESS"),
			},
		},
		{
			name: "Wrong BASE_URL",
			envs: []env{
				{name: "SERVER_ADDRESS", value: "127.0.0.1:8080"},
				{name: "BASE_URL", value: "htt://127.0.0.1:8080/short"},
			},
			want: want{
				config: Config{
					ServerURL: "127.0.0.1:8080",
					BaseURL:   "htt://127.0.0.1:8080/short",
					Address:   "",
					Port:      0,
					Path:      "",
				},
				err: fmt.Errorf("invalid BASE_URL"),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			for _, env := range tt.envs {
				os.Setenv(env.name, env.value)
			}
			cfg := Config{}
			err := cfg.Load()
			assert.Equal(t, tt.want.config, cfg)
			assert.Equal(t, tt.want.err, err)
		})
	}
}
