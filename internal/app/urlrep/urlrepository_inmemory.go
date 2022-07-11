package urlrep

import (
	"github.com/google/uuid"
	"sync"
)

type inMemoryURLRepository struct {
	urls           sync.Map
	urlIDGenerator func(url string) string
}

func NewInMemory(urlIDGenerator func(url string) string) URLRepository {
	if urlIDGenerator == nil {
		urlIDGenerator = func(url string) string {
			return uuid.New().String()[:8]
		}
	}
	return &inMemoryURLRepository{
		urlIDGenerator: urlIDGenerator,
	}
}

func (r *inMemoryURLRepository) Set(url string) string {
	id := r.GenerateID(url)
	if _, ok := r.urls.Load(id); ok {
		return r.Set(url)
	}
	r.urls.Store(id, url)
	return id
}

func (r *inMemoryURLRepository) Get(id string) string {
	url, ok := r.urls.Load(id)
	if ok {
		return url.(string)
	}
	return ""
}

func (r *inMemoryURLRepository) GenerateID(url string) string {
	return r.urlIDGenerator(url)
}
