package main

import (
	"encoding/json"
	"net/http"
)

func handlerChirpValidate(w http.ResponseWriter, r *http.Request) {
	type reqParameters struct {
		Body string `json:"body"`
	}

	type resBody struct {
		Valid bool `json:"valid"`
	}

	decoder := json.NewDecoder(r.Body)
	params := reqParameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters", err)
		return
	}

	const maxChirpLength = 140
	if len(params.Body) > maxChirpLength {
		respondWithError(w, http.StatusBadRequest, "Chirp is too long", nil)
		return
	}

	respondWithJSON(w, http.StatusOK, resBody{
		Valid: true,
	})
}
