package main

import (
	"github.com/dmitribauer/go-url-shortener/internal/app/servers"
	"github.com/dmitribauer/go-url-shortener/internal/app/urlrep"
)

func main() {
	urlRep := urlrep.NewInMemory()
	s := servers.NewDefault(urlRep)
	err := s.Start()
	if err != nil {
		panic(err)
	}
}
