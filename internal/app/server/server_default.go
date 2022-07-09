package server

import (
	"github.com/dmitribauer/go-url-shortener/internal/app/urlrep"
	"io/ioutil"
	"net/http"
	"net/url"
)

type defaultServer struct {
	urlRepository urlrep.UrlRepository
}

func NewDefault(urlRepository urlrep.UrlRepository) Server {
	return &defaultServer{
		urlRepository: urlRepository,
	}
}

func (s *defaultServer) Start() error {
	http.HandleFunc("/post", s.postUrl)
	http.HandleFunc("/get/", s.getUrl)
	httpServ := &http.Server{
		Addr: ":8080",
	}
	return httpServ.ListenAndServe()
}

func (s *defaultServer) postUrl(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	defer r.Body.Close()
	body, err := ioutil.ReadAll(r.Body)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	str := string(body)
	if str == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	_, err = url.ParseRequestURI(str)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	id := s.urlRepository.Set(str)
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(id))
}

func (s *defaultServer) getUrl(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	id := r.URL.Path[len("/get/"):]
	if id == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	url := s.urlRepository.Get(id)
	if url == "" {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	w.Header().Set("Location", url)
	w.WriteHeader(http.StatusTemporaryRedirect)
}
