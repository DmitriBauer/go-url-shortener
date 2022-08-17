package handlers

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/asaskevich/govalidator"

	"github.com/dmitribauer/go-url-shortener/internal/api/rest"
	"github.com/dmitribauer/go-url-shortener/internal/reqrep"
	"github.com/dmitribauer/go-url-shortener/internal/urlrep"
)

type shortenReqBody struct {
	URL string `json:"url" valid:"required,url"`
}

type shortenResBody struct {
	Result string `json:"result"`
}

func HandleShortenPost(rest *rest.Rest, w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("Content-Type") != "application/json" {
		w.WriteHeader(http.StatusUnsupportedMediaType)
		return
	}

	defer r.Body.Close()
	var reqBody shortenReqBody
	if json.NewDecoder(r.Body).Decode(&reqBody) != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	ok, _ := govalidator.ValidateStruct(reqBody)
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var statusCode int

	urlID, err := rest.URLRepo.Save(r.Context(), reqBody.URL)
	if err == nil {
		statusCode = http.StatusCreated
	} else if errors.Is(err, urlrep.ErrDuplicateURL) {
		statusCode = http.StatusConflict
	} else {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	resBody := shortenResBody{
		Result: rest.ShortURL(urlID),
	}

	err = rest.ReqRepo.Save(reqrep.Req{
		SessionID:   sessionIDFromRequest(rest, w, r),
		ShortURL:    resBody.Result,
		OriginalURL: reqBody.URL,
	})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(resBody)
}
