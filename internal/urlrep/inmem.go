package urlrep

import (
	"sync"

	"github.com/google/uuid"
)

type inMemURLRepo struct {
	urls           sync.Map
	urlIDGenerator func(url string) string
}

func NewInMemory(urlIDGenerator func(url string) string) URLRepo {
	if urlIDGenerator == nil {
		urlIDGenerator = func(url string) string {
			return uuid.New().String()[:8]
		}
	}
	return &inMemURLRepo{
		urlIDGenerator: urlIDGenerator,
	}
}

func (r *inMemURLRepo) URLByID(id string) string {
	url, ok := r.urls.Load(id)
	if ok {
		return url.(string)
	}
	return ""
}

func (r *inMemURLRepo) Save(url string) string {
	id := r.GenerateID(url)
	if _, ok := r.urls.Load(id); ok {
		return r.Save(url)
	}
	r.urls.Store(id, url)
	return id
}

func (r *inMemURLRepo) GenerateID(url string) string {
	return r.urlIDGenerator(url)
}
