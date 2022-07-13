package main

import (
	"github.com/dmitribauer/go-url-shortener/internal/api"
	"github.com/dmitribauer/go-url-shortener/internal/urlrep"
)

func main() {
	urlRepo := urlrep.NewInMemory(nil)
	rest := api.NewRest(urlRepo)
	err := rest.Run("localhost", 8080)
	if err != nil {
		panic(err)
	}
}
