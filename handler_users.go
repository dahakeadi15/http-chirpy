package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/google/uuid"
)

type User struct {
	Id        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Email     string    `json:"email"`
}

func (cfg *apiConfig) handlerCreateUser(w http.ResponseWriter, r *http.Request) {
	type reqParameters struct {
		Email string `json:"email"`
	}

	type resBody struct {
		User
	}

	decoder := json.NewDecoder(r.Body)
	params := reqParameters{}
	if err := decoder.Decode(&params); err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters", err)
		return
	}

	newUser, err := cfg.db.CreateUser(r.Context(), params.Email)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't create the user", err)
		return
	}

	respondWithJSON(w, http.StatusCreated, resBody{
		User: User{Id: newUser.ID,
			CreatedAt: newUser.CreatedAt,
			UpdatedAt: newUser.UpdatedAt,
			Email:     newUser.Email,
		}})
}
