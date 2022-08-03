package conf

import (
	"fmt"
	"net/url"
	"os"
	"strconv"
	"strings"

	"github.com/caarlos0/env/v6"
)

const (
	serverAddressDefault = "localhost:8080"
	baseURLDefault       = "http://localhost:8080"
)

type Config struct {
	ServerAddress   string `env:"SERVER_ADDRESS"`
	BaseURL         string `env:"BASE_URL"`
	FileStoragePath string `env:"FILE_STORAGE_PATH"`
	Address         string
	Port            int
	Path            string
}

// Load loads the config from the app arguments and environment variables.
// If either fails, Load returns an error.
// The arguments take precedence over the environment variables.
func (cfg *Config) Load() error {
	err := cfg.loadEnvs()
	if err != nil {
		return err
	}

	err = cfg.loadArgs()
	if err != nil {
		return err
	}

	cfg.check()
	return nil
}

func (cfg *Config) loadEnvs() error {
	var address string
	var port int
	var path string

	err := env.Parse(cfg, env.Options{
		OnSet: func(tag string, value interface{}, isDefault bool) {
			if tag == "SERVER_ADDRESS" {
				address, port = parseServerAddress(value.(string))
			} else if tag == "BASE_URL" {
				path = parseBaseURL(value.(string))
			}
		},
	})

	if err != nil {
		return err
	}

	cfg.Address = address
	cfg.Port = port
	cfg.Path = path

	return nil
}

func (cfg *Config) loadArgs() error {
	args := map[string]string{}
	for _, arg := range os.Args {
		kv := strings.Split(arg, "=")
		if len(kv) != 2 {
			continue
		}
		args[kv[0]] = kv[1]
	}

	serverAddress, ok := args["-a"]
	if ok {
		address, port := parseServerAddress(serverAddress)
		if address == "" || port == 0 {
			return fmt.Errorf("invalid server address '-a %s'", serverAddress)
		}
		cfg.ServerAddress = serverAddress
		cfg.Address = address
		cfg.Port = port
	}

	baseURL, ok := args["-b"]
	if ok {
		path := parseBaseURL(baseURL)
		if path == "" {
			return fmt.Errorf("invalid base URL '-b %s'", baseURL)
		}
		cfg.BaseURL = baseURL
		cfg.Path = path
	}

	fileStoragePath, ok := args["-f"]
	if ok {
		cfg.FileStoragePath = fileStoragePath
	}

	return nil
}

func (cfg *Config) check() {
	if cfg.ServerAddress == "" || cfg.BaseURL == "" {
		cfg.ServerAddress = "localhost:8080"
		cfg.BaseURL = "http://localhost:8080"
		cfg.Address = "localhost"
		cfg.Port = 8080
		cfg.Path = "/"
	}
}

func parseServerAddress(serverAddress string) (address string, port int) {
	values := strings.Split(serverAddress, ":")
	if len(values) != 2 {
		address, port = "", 0
		return
	}
	addr := values[0]
	p, err := strconv.Atoi(values[1])
	if err != nil {
		address, port = "", 0
		return
	}
	address, port = addr, p
	return
}

func parseBaseURL(baseURL string) (path string) {
	lastChIdx := len(baseURL) - 1
	if lastChIdx < 0 {
		return
	}
	if baseURL[lastChIdx] != '/' {
		baseURL += "/"
	}

	u, err := url.ParseRequestURI(baseURL)

	if err != nil || u.Host == "" {
		u, err = url.ParseRequestURI("http://" + baseURL)
		if err != nil {
			return
		}
		path = u.Path
		return
	}

	if u.Scheme != "http" && u.Scheme != "https" {
		return
	}

	path = u.Path
	return
}
