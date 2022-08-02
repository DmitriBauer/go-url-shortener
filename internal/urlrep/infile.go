package urlrep

import (
	"bufio"
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

func (r *inFileURLRepo) URLByID(id string) (string, bool) {
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
		if line[0] == id {
			return line[1], true
		}
	}

	return "", false
}

func (r *inFileURLRepo) Save(url string) (string, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	f, err := os.OpenFile(r.path, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0777)
	defer r.closeFile(f)
	if err != nil {
		return "", err
	}

	id := r.GenerateID(url)
	_, err = f.WriteString(id + " " + url + "\n")
	if err != nil {
		return "", err
	}

	return id, nil
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
