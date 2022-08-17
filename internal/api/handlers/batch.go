package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/dmitribauer/go-url-shortener/internal/api/rest"
)

type batchReqBodyEntry struct {
	CorrID string `json:"correlation_id"`
	URL    string `json:"original_url"`
}

type batchResBodyEntry struct {
	CorrID   string `json:"correlation_id"`
	ShortURL string `json:"short_url"`
}

func HandleShortenBatchPost(rest *rest.Rest, w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("Content-Type") != "application/json" {
		w.WriteHeader(http.StatusUnsupportedMediaType)
		return
	}

	defer r.Body.Close()
	var reqBody []batchReqBodyEntry
	if json.NewDecoder(r.Body).Decode(&reqBody) != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	urls := make([]string, len(reqBody))
	for i, entry := range reqBody {
		urls[i] = entry.URL
	}

	urlIDs, err := rest.URLRepo.SaveList(r.Context(), urls)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	resBody := make([]batchResBodyEntry, len(reqBody))
	for i, urlID := range urlIDs {
		resBody[i] = batchResBodyEntry{
			CorrID:   reqBody[i].CorrID,
			ShortURL: rest.ShortURL(urlID),
		}
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(resBody)
}
