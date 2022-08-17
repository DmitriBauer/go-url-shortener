package urlrep

import (
	"context"
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

func (r *inMemURLRepo) URLByID(ctx context.Context, id string) (string, bool) {
	url, ok := r.urls.Load(id)
	if !ok {
		return "", false
	}
	return url.(string), true
}

func (r *inMemURLRepo) Save(ctx context.Context, url string) (string, error) {
	id := r.GenerateID(url)
	if _, ok := r.urls.Load(id); ok {
		return r.Save(ctx, url)
	}
	r.urls.Store(id, url)
	return id, nil
}

func (r *inMemURLRepo) SaveList(ctx context.Context, urls []string) ([]string, error) {
	idxs := make([]string, len(urls))
	for i, url := range urls {
		idx, err := r.Save(ctx, url)
		if err != nil {
			return nil, err
		}
		idxs[i] = idx
	}
	return idxs, nil
}

func (r *inMemURLRepo) GenerateID(url string) string {
	return r.urlIDGenerator(url)
}
