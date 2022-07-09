package server

import (
	"github.com/dmitribauer/go-url-shortener/internal/app/urlrep"
	"io/ioutil"
	"net/http"
	urls "net/url"
)

type defaultServer struct {
	urlRepository urlrep.URLRepository
}

func NewDefault(urlRepository urlrep.URLRepository) Server {
	return &defaultServer{
		urlRepository: urlRepository,
	}
}

func (s *defaultServer) Start() error {
	http.HandleFunc("/", s.handleRoot)
	httpServ := &http.Server{
		Addr: ":8080",
	}
	return httpServ.ListenAndServe()
}

func (s *defaultServer) handleRoot(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		s.postURL(w, r)
	case http.MethodGet:
		s.getURL(w, r)
	default:
		w.WriteHeader(http.StatusBadRequest)
	}
}

func (s *defaultServer) postURL(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	body, err := ioutil.ReadAll(r.Body)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	url := string(body)
	if url == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	_, err = urls.ParseRequestURI(url)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	id := s.urlRepository.Set(url)
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(id))
}

func (s *defaultServer) getURL(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Path[len("/"):]
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
