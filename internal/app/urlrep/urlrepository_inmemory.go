package urlrep

import (
	"github.com/google/uuid"
	"sync"
)

type inMemoryUrlRepository struct {
	urls sync.Map
}

func NewInMemory() UrlRepository {
	return &inMemoryUrlRepository{}
}

func (r *inMemoryUrlRepository) Set(url string) string {
	id := r.GenerateId(url)
	if _, ok := r.urls.Load(id); ok {
		return r.Set(url)
	}
	r.urls.Store(id, url)
	return id
}

func (r *inMemoryUrlRepository) Get(id string) string {
	url, ok := r.urls.Load(id)
	if ok {
		return url.(string)
	}
	return ""
}

func (*inMemoryUrlRepository) GenerateId(url string) string {
	return uuid.New().String()[:8]
}
