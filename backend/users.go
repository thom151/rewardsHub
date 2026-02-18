package main

import (
	"encoding/json"
	"net/http"
)

func (cfg *apiConfig) handlerCreateUser(w http.ResponseWriter, r *http.Request) {
	type params struct {
		Email string `json:"email"`
	}

	decoder := json.NewDecoder(r.Body)
	userParams := params{}

	err := decoder.Decode(&userParams)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "couldn't decode params", err)
		return
	}

	user, err := cfg.db.CreateUser(r.Context(), userParams.Email)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "couldn't crease user", err)
		return
	}

	respondWithJSON(w, http.StatusOK, user)
}
