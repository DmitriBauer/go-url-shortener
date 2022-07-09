package main

import (
	"github.com/dmitribauer/go-url-shortener/internal/app/server"
	"github.com/dmitribauer/go-url-shortener/internal/app/urlrep"
)

func main() {
	urlRep := urlrep.NewInMemory()
	s := server.NewDefault(urlRep)
	err := s.Start()
	if err != nil {
		panic(err)
	}
}
