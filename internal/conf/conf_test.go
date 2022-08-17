package conf

import (
	"fmt"
	"os"
	"strings"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestConfig_Load(t *testing.T) {
	wd, err := os.Getwd()
	require.NoError(t, err)
	fileStoragePath := wd + "/TestConfig_Load_Storage"
	defer os.Remove(fileStoragePath)
	type arg struct {
		name  string
		value string
	}
	type want struct {
		config Config
		err    error
	}
	tests := []struct {
		name string
		envs []arg
		args []arg
		want want
	}{
		{
			name: "No envs, no args",
			envs: []arg{},
			want: want{
				config: Config{
					ServerAddress: "localhost:8080",
					BaseURL:       "http://localhost:8080",
					ReqRepoDir:    "/tmp/reqrep/",
					Address:       "localhost",
					Port:          8080,
				},
				err: nil,
			},
		},
		{
			name: "Correct envs, no args",
			envs: []arg{
				{name: "SERVER_ADDRESS", value: "127.0.0.1:8080"},
				{name: "BASE_URL", value: "http://127.0.0.1:8080/short"},
			},
			want: want{
				config: Config{
					ServerAddress: "127.0.0.1:8080",
					BaseURL:       "http://127.0.0.1:8080/short",
					ReqRepoDir:    "/tmp/reqrep/",
					Address:       "127.0.0.1",
					Port:          8080,
					Path:          "/short",
				},
				err: nil,
			},
		},

		{
			name: "No envs, correct args",
			args: []arg{
				{name: "a", value: "127.0.0.1:8181"},
				{name: "b", value: "http://127.0.0.1:8181/short"},
				{name: "f", value: fileStoragePath},
				{name: "d", value: "postgres://postgres:postgres@postgres:5432/praktikum?sslmode=disable"},
			},
			want: want{
				config: Config{
					ServerAddress:   "127.0.0.1:8181",
					BaseURL:         "http://127.0.0.1:8181/short",
					FileStoragePath: fileStoragePath,
					ReqRepoDir:      "/tmp/reqrep/",
					DatabaseAddress: "postgres://postgres:postgres@postgres:5432/praktikum?sslmode=disable",
					Address:         "127.0.0.1",
					Port:            8181,
					Path:            "/short",
				},
				err: nil,
			},
		},
		{
			name: "No envs, missed -a",
			args: []arg{
				{name: "b", value: "http://127.0.0.1:8181/short"},
				{name: "f", value: fileStoragePath},
			},
			want: want{
				config: Config{
					ServerAddress:   "localhost:8080",
					BaseURL:         "http://127.0.0.1:8181/short",
					FileStoragePath: fileStoragePath,
					ReqRepoDir:      "/tmp/reqrep/",
					Address:         "localhost",
					Port:            8080,
					Path:            "/short",
				},
				err: nil,
			},
		},
		{
			name: "No envs, missed -b",
			args: []arg{
				{name: "a", value: "127.0.0.1:8181"},
				{name: "f", value: fileStoragePath},
			},
			want: want{
				config: Config{
					ServerAddress:   "127.0.0.1:8181",
					BaseURL:         "http://localhost:8080",
					FileStoragePath: fileStoragePath,
					ReqRepoDir:      "/tmp/reqrep/",
					Address:         "127.0.0.1",
					Port:            8181,
				},
				err: nil,
			},
		},
		{
			name: "No envs, missed -f",
			args: []arg{
				{name: "a", value: "127.0.0.1:8181"},
				{name: "b", value: "http://127.0.0.1:8181/short"},
			},
			want: want{
				config: Config{
					ServerAddress:   "127.0.0.1:8181",
					BaseURL:         "http://127.0.0.1:8181/short",
					FileStoragePath: "",
					ReqRepoDir:      "/tmp/reqrep/",
					Address:         "127.0.0.1",
					Port:            8181,
					Path:            "/short",
				},
				err: nil,
			},
		},
	}

	defaultArgs := os.Args
	defaultEnvs := os.Environ()
	mu := &sync.Mutex{}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mu.Lock()

			os.Clearenv()
			for _, env := range defaultEnvs {
				kv := strings.Split(env, "=")
				if len(kv) != 2 {
					continue
				}
				os.Setenv(kv[0], kv[1])
			}
			for _, env := range tt.envs {
				os.Setenv(env.name, env.value)
			}

			os.Args = defaultArgs
			for _, arg := range tt.args {
				os.Args = append(os.Args, fmt.Sprintf("-%s=%s", arg.name, arg.value))
			}

			cfg := Config{}
			err := cfg.Load()
			assert.Equal(t, tt.want.config, cfg)
			assert.Equal(t, tt.want.err, err)

			mu.Unlock()
		})
	}
}
