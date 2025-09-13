package main

import (
	"encoding/json"
	"log"
	"net/http"
)

func handlerValidation(w http.ResponseWriter, req *http.Request) {
	// decoding request
	decoder := json.NewDecoder(req.Body)
	c := chirp{}
	err := decoder.Decode(&c)
	if err != nil {
		log.Printf("Error decoding parameters: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// chirp validation
	if len(c.Body) < 141 {
		resBody := validResponse{
			Valid: true,
		}

		// encoding response
		res, err := json.Marshal(resBody)
		if err != nil {
			log.Printf("Error marshalling JSON: %s", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(res)
		return
	} else {
		resBody := errorResponse{
			Error: "Chirp is too long",
		}

		// encoding response
		res, err := json.Marshal(resBody)
		if err != nil {
			log.Printf("Error marshalling JSON: %s", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		w.Write(res)
		return
	}
}

type chirp struct {
	Body string `json:"body"`
}

type errorResponse struct {
	Error string `json:"error"`
}

type validResponse struct {
	Valid bool `json:"valid"`
}
