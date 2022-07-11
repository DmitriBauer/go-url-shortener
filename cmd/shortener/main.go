package main

import (
	"github.com/dmitribauer/go-url-shortener/internal/app/api"
	"github.com/dmitribauer/go-url-shortener/internal/app/urlrep"
)

func main() {
	urlRepository := urlrep.NewInMemory(nil)
	rest := api.NewRest(urlRepository)
	err := rest.Run("localhost", 8080)
	if err != nil {
		panic(err)
	}
}
