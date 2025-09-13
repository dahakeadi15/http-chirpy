package main

import (
	"encoding/json"
	"net/http"
	"slices"
	"strings"
)

func handlerChirpValidate(w http.ResponseWriter, r *http.Request) {
	type reqParameters struct {
		Body string `json:"body"`
	}

	type resBody struct {
		CleanedBody string `json:"cleaned_body"`
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

	profaneWords := []string{"kerfuffle", "sharbert", "fornax"}
	cleanedBodyArr := []string{}
	for word := range strings.SplitSeq(params.Body, " ") {
		if slices.Contains(profaneWords, strings.ToLower(word)) {
			cleanedBodyArr = append(cleanedBodyArr, "****")
			continue
		}
		cleanedBodyArr = append(cleanedBodyArr, word)
	}

	respondWithJSON(w, http.StatusOK, resBody{
		CleanedBody: strings.Join(cleanedBodyArr, " "),
	})
}
