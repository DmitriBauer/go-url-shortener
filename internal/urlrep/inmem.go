package urlrep

import (
	"sync"

	"github.com/google/uuid"
)

type inMemURLRepo struct {
	urls           sync.Map
	urlIDGenerator func(url string) string
}

// NewInMemory returns new in-memory URL repository.
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

func (r *inMemURLRepo) URLByID(id string) (string, bool) {
	url, ok := r.urls.Load(id)
	if !ok {
		return "", false
	}
	return url.(string), true
}

func (r *inMemURLRepo) Save(url string) (string, error) {
	id := r.GenerateID(url)
	if _, ok := r.urls.Load(id); ok {
		return r.Save(url)
	}
	r.urls.Store(id, url)
	return id, nil
}

func (r *inMemURLRepo) GenerateID(url string) string {
	return r.urlIDGenerator(url)
}
