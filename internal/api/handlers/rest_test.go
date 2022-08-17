package handlers

import (
	apirest "github.com/dmitribauer/go-url-shortener/internal/api/rest"
	"github.com/dmitribauer/go-url-shortener/internal/auth"
	"github.com/dmitribauer/go-url-shortener/internal/reqrep"
	"github.com/dmitribauer/go-url-shortener/internal/urlrep"
)

func newTestRest(rest *apirest.Rest) *apirest.Rest {
	if rest == nil {
		rest = &apirest.Rest{}
	}
	if rest.Address == "" {
		rest.Address = "127.0.0.1"
	}
	if rest.Port == 0 {
		rest.Port = 8282
	}
	if rest.Path == "" {
		rest.Path = "/s"
	}
	if rest.URLRepo == nil {
		urlIDGenerator := func(url string) string {
			return "uRlId123"
		}
		rest.URLRepo = urlrep.NewInMemory(urlIDGenerator)
	}
	if rest.ReqRepo == nil {
		reqRepo, _ := reqrep.NewInFile("/tmp/reqrep/")
		rest.ReqRepo = reqRepo
	}
	if rest.AuthService == nil {
		sessionIDGenerator := func() string {
			return "41734162-d37e-4936-ab4c-b808386a34c9"
		}
		rest.AuthService = auth.NewService(sessionIDGenerator)
	}
	return rest
}
