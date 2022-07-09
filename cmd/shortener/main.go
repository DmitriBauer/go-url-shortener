package main

import (
	"github.com/dmitribauer/go-url-shortener/internal/app/server"
	"github.com/dmitribauer/go-url-shortener/internal/app/urlrep"
)

func main() {
	urlRep := urlrep.NewInMemory()
	serv := server.NewDefault(urlRep)
	err := serv.Start()
	if err != nil {
		panic(err)
	}
}
