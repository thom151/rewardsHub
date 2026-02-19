package main

import (
	"database/sql"
	"encoding/json"
	"net/http"

	"github.com/thom151/rewardsHub/internal/database"
)

func (cfg *apiConfig) handlerCreateUser(w http.ResponseWriter, r *http.Request) {
	type params struct {
		Email      string `json:"email"`
		FirtstName string `json:"first_name"`
		LastName   string `json:"last_name"`
		Phone      string `json:"phone"`
	}

	decoder := json.NewDecoder(r.Body)
	var userParams params
	err := decoder.Decode(&userParams)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "cannot decode user parameters", err)
		return
	}

	createUserParams := database.CreateUserParams{
		Email:     userParams.Email,
		FirstName: userParams.FirtstName,
		LastName:  userParams.LastName,
		Phone: sql.NullString{
			String: userParams.Phone,
			Valid:  userParams.Phone != "",
		},
	}

	newUser, err := cfg.db.CreateUser(r.Context(), createUserParams)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "failed to create uer", err)
		return
	}

	respondWithJSON(w, http.StatusOK, newUser)
}
