package reqrep

import (
	"encoding/json"
	"errors"
	"io"
	"io/ioutil"
	"log"
	"os"
	"sync"
)

type inFileReqRepo struct {
	dir string
	mu  sync.Mutex
}

func NewInFile(dir string) (ReqRepo, error) {
	if _, err := os.Stat(dir); err != nil {
		err := os.MkdirAll(dir, 0777)
		if err != nil {
			return nil, err
		}
	}
	return &inFileReqRepo{
		dir: dir,
	}, nil
}

func (repo *inFileReqRepo) DataBySessionID(userID string) ([]byte, error) {
	repo.mu.Lock()
	defer repo.mu.Unlock()

	f, err := os.OpenFile(repo.fileNameBySessionID(userID), os.O_RDONLY, 0666)
	if err != nil {
		return nil, err
	}
	defer repo.closeFile(f)

	return io.ReadAll(f)
}

func (repo *inFileReqRepo) Save(req Req) error {
	repo.mu.Lock()
	defer repo.mu.Unlock()

	path := repo.fileNameBySessionID(req.SessionID)

	if _, err := os.Stat(path); errors.Is(err, os.ErrNotExist) {
		rr := []ReqRecord{
			{
				ShortURL:    req.ShortURL,
				OriginalURL: req.OriginalURL,
			},
		}
		return repo.writeRecords(rr, path)
	} else if err != nil {
		return err
	}

	data, err := ioutil.ReadFile(repo.fileNameBySessionID(req.SessionID))
	if err != nil {
		return err
	}

	var rr []ReqRecord

	err = json.Unmarshal(data, &rr)
	if err != nil {
		return err
	}

	rr = append(rr, ReqRecord{
		ShortURL:    req.ShortURL,
		OriginalURL: req.OriginalURL,
	})

	return repo.writeRecords(rr, path)
}

func (repo *inFileReqRepo) writeRecords(rr []ReqRecord, file string) error {
	data, err := json.Marshal(rr)
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(file, data, 0666)
	if err != nil {
		return err
	}

	return nil
}

func (repo *inFileReqRepo) fileNameBySessionID(sessionID string) string {
	return repo.dir + sessionID
}

func (repo *inFileReqRepo) closeFile(f *os.File) {
	err := f.Close()
	if err != nil {
		log.Default().Println("failed to close file:", err)
	}
}
