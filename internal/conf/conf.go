package conf

import (
	"flag"
	"fmt"
	"net/url"
	"strconv"
	"strings"

	"github.com/caarlos0/env/v6"
)

// ...conf.test flag redefined: a [recovered]
// https://groups.google.com/g/golang-nuts/c/1aZmhhSvwWc

const (
	serverAddressDefault = "SERVER_ADDRESS_DEFAULT"
	baseURLDefault       = "BASE_URL_DEFAULT"
)

var (
	serverAddress   *string
	baseURL         *string
	fileStoragePath *string
)

func init() {
	serverAddress = flag.String("a", serverAddressDefault, "server address")
	baseURL = flag.String("b", baseURLDefault, "base URL")
	fileStoragePath = flag.String("f", "", "file storage path")
	flag.Parse()
}

type Config struct {
	ServerURL       string `env:"SERVER_ADDRESS" envDefault:"localhost:8080"`
	BaseURL         string `env:"BASE_URL" envDefault:"http://localhost:8080"`
	FileStoragePath string `env:"FILE_STORAGE_PATH"`
	Address         string
	Port            int
	Path            string
}

// Load loads the config from the app arguments and environment variables.
// If either fails, Load returns an error.
// The arguments take precedence over the environment variables.
func (cfg *Config) Load() error {
	ok, err := cfg.loadArgs()
	if err != nil {
		return err
	}

	if ok {
		return nil
	}

	err = cfg.loadEnvs()
	if err != nil {
		return err
	}

	return nil
}

func (cfg *Config) loadArgs() (bool, error) {
	if *serverAddress == serverAddressDefault && *baseURL == baseURLDefault {
		return false, nil
	} else if *serverAddress == serverAddressDefault || *baseURL == baseURLDefault {
		return false, fmt.Errorf("should specify both server address (-a) and base URL (-b)")
	}

	address, port := parseServerAddress(*serverAddress)
	if address == "" || port == 0 {
		return false, fmt.Errorf("invalid server address '-a %s'", *serverAddress)
	}

	path := parseBaseURL(*baseURL)
	if path == "" {
		return false, fmt.Errorf("invalid base URL '-b %s'", *baseURL)
	}

	cfg.Address = address
	cfg.Port = port
	cfg.Path = path
	cfg.FileStoragePath = *fileStoragePath

	return true, nil
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
