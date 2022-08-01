package conf

import (
	"fmt"
	"net/url"
	"strconv"
	"strings"

	"github.com/caarlos0/env/v6"
)

type Config struct {
	ServerURL string `env:"SERVER_ADDRESS" envDefault:"localhost:8080"`
	BaseURL   string `env:"BASE_URL" envDefault:"http://localhost:8080"`
	Address   string
	Port      int
	Path      string
}

func (cfg *Config) Load() error {
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

	if address == "" || port == 0 {
		return fmt.Errorf("invalid SERVER_ADDRESS")
	}

	if path == "" {
		return fmt.Errorf("invalid BASE_URL")
	}

	cfg.Address = address
	cfg.Port = port
	cfg.Path = path

	return nil
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
