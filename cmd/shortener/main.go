package main

import (
	"github.com/dmitribauer/go-url-shortener/internal/api"
	"github.com/dmitribauer/go-url-shortener/internal/auth"
	"github.com/dmitribauer/go-url-shortener/internal/conf"
	"github.com/dmitribauer/go-url-shortener/internal/reqrep"
	"github.com/dmitribauer/go-url-shortener/internal/urlrep"
)

func main() {
	var cfg conf.Config
	err := cfg.Load()
	if err != nil {
		panic(err)
	}

	var urlRepo urlrep.URLRepo
	if cfg.DatabaseAddress != "" {
		urlRepo, err = urlrep.NewPostgre(cfg.DatabaseAddress, nil)
		if err != nil {
			panic(err)
		}
	} else if cfg.FileStoragePath != "" {
		urlRepo, err = urlrep.NewInFile(cfg.FileStoragePath, nil)
		if err != nil {
			panic(err)
		}
	} else {
		urlRepo = urlrep.NewInMemory(nil)
	}

	reqRepo, err := reqrep.NewInFile(cfg.ReqRepoDir)
	if err != nil {
		panic(err)
	}

	authService := auth.NewService(nil)

	rest := api.NewRest(
		urlRepo,
		reqRepo,
		authService,
	)

	err = api.Run(rest, cfg.Address, cfg.Port, cfg.Path)
	if err != nil {
		panic(err)
	}
}
