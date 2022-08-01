package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/asaskevich/govalidator"

	"github.com/dmitribauer/go-url-shortener/internal/api/rest"
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

	urlID, err := rest.URLRepo.Save(reqBody.URL)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	resBody := shortenResBody{
		Result: rest.ShortURL(urlID),
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(resBody)
}
