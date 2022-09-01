package urlrep

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"
	"strings"
	"sync"

	"github.com/google/uuid"
)

type inFileURLRepo struct {
	path           string
	urlIDGenerator func(url string) string
	mu             sync.Mutex
}

// NewInFile returns new in-file URL repository.
func NewInFile(path string, urlIDGenerator func(url string) string) (URLRepo, error) {
	f, err := os.OpenFile(path, os.O_RDONLY|os.O_CREATE, 0777)
	if err != nil {
		return nil, err
	}
	err = f.Close()
	if err != nil {
		return nil, err
	}

	if urlIDGenerator == nil {
		urlIDGenerator = func(url string) string {
			return uuid.New().String()[:8]
		}
	}

	return &inFileURLRepo{
		path:           path,
		urlIDGenerator: urlIDGenerator,
	}, nil
}

func (r *inFileURLRepo) URLByID(ctx context.Context, id string) (string, bool) {
	r.mu.Lock()
	defer r.mu.Unlock()

	f, err := os.OpenFile(r.path, os.O_RDONLY|os.O_CREATE, 0777)
	defer r.closeFile(f)
	if err != nil {
		return "", false
	}

	s := bufio.NewScanner(f)
	s.Split(newInFileSplitter().scanLinesInReverse)
	for s.Scan() {
		line := strings.Split(s.Text(), " ")
		if lineIsNotValid(line) {
			continue
		}
		if line[0] == id {
			return line[1], line[3] == "1"
		}
	}

	return "", false
}

func (r *inFileURLRepo) Save(ctx context.Context, url string, sessionID string) (string, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	f, err := os.OpenFile(r.path, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0777)
	defer r.closeFile(f)
	if err != nil {
		return "", err
	}

	id := r.GenerateID(url)
	_, err = f.WriteString(fmt.Sprintf("%s %s %s %s\n", id, url, sessionID, "0"))
	if err != nil {
		return "", err
	}

	return id, nil
}

func (r *inFileURLRepo) SaveList(ctx context.Context, urls []string, sessionID string) ([]string, error) {
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

func (r *inFileURLRepo) RemoveList(ctx context.Context, ids []string, sessionID string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	f, err := os.OpenFile(r.path, os.O_RDWR|os.O_CREATE, 0777)
	if err != nil {
		return err
	}
	defer r.closeFile(f)

	offset := 0
	s := bufio.NewScanner(f)
	for s.Scan() {
		text := s.Text()
		line := strings.Split(text, " ")
		if lineIsNotValid(line) {
			continue
		}
		got := false

		offset += len(text) + 1

		for _, id := range ids {
			if line[0] == id {
				got = true
			}
		}

		if !got || line[2] != sessionID {
			continue
		}

		f.WriteAt([]byte("1"), int64(offset-2))
	}

	return nil
}

func (r *inFileURLRepo) GenerateID(url string) string {
	return r.urlIDGenerator(url)
}

func (r *inFileURLRepo) closeFile(f *os.File) {
	err := f.Close()
	if err != nil {
		log.Default().Println("failed to close file:", err)
	}
}

func lineIsNotValid(line []string) bool {
	return len(line) != 4
}
