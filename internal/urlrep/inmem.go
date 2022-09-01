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

type inmemRecord struct {
	url       string
	sessionID string
	removed   bool
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
	v, ok := r.urls.Load(id)
	if !ok {
		return "", false
	}

	rec := v.(inmemRecord)

	return rec.url, rec.removed
}

func (r *inMemURLRepo) Save(ctx context.Context, url string, sessionID string) (string, error) {
	id := r.GenerateID(url)
	if _, ok := r.urls.Load(id); ok {
		return r.Save(ctx, url, sessionID)
	}
	r.urls.Store(id, inmemRecord{url: url, sessionID: sessionID})
	return id, nil
}

func (r *inMemURLRepo) SaveList(ctx context.Context, urls []string, sessionID string) ([]string, error) {
	idxs := make([]string, len(urls))
	for i, url := range urls {
		idx, err := r.Save(ctx, url, sessionID)
		if err != nil {
			return nil, err
		}
		idxs[i] = idx
	}
	return idxs, nil
}

func (r *inMemURLRepo) RemoveList(ctx context.Context, ids []string, sessionID string) error {
	for _, id := range ids {
		v, ok := r.urls.LoadAndDelete(id)
		if !ok {
			continue
		}

		rec := v.(inmemRecord)
		if rec.sessionID == sessionID {
			rec.removed = true
		}

		r.urls.Store(id, rec)
	}

	return nil
}

func (r *inMemURLRepo) GenerateID(url string) string {
	return r.urlIDGenerator(url)
}
