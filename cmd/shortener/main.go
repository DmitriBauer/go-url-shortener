package main

import (
	"github.com/dmitribauer/go-url-shortener/internal/api"
	"github.com/dmitribauer/go-url-shortener/internal/conf"
	"github.com/dmitribauer/go-url-shortener/internal/urlrep"
)

func main() {
	var cfg conf.Config
	err := cfg.Load()
	if err != nil {
		panic(err)
	}

	var urlRepo urlrep.URLRepo
	if cfg.FileStoragePath != "" {
		urlRepo, err = urlrep.NewInFile(cfg.FileStoragePath, nil)
		if err != nil {
			panic(err)
		}
	} else {
		urlRepo = urlrep.NewInMemory(nil)
	}

	rest := api.NewRest(urlRepo)
	err = api.Run(rest, cfg.Address, cfg.Port, cfg.Path)
	if err != nil {
		panic(err)
	}
}
