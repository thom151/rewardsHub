package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"

	"github.com/thom151/rewardsHub/internal/auth"
	"github.com/thom151/rewardsHub/internal/database"
)

func (cfg *apiConfig) handlerCreateUser(w http.ResponseWriter, r *http.Request) {
	type params struct {
		Email      string `json:"email"`
		FirtstName string `json:"first_name"`
		LastName   string `json:"last_name"`
		Phone      string `json:"phone"`
		Password   string `json:"password"`
	}

	decoder := json.NewDecoder(r.Body)
	var userParams params
	err := decoder.Decode(&userParams)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "cannot decode user parameters", err)
		return
	}

	if userParams.Password == "" {
		respondWithError(w, http.StatusBadRequest, "missing password", err)
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

	hashed, err := auth.HashPassword(userParams.Password)
	if err != nil {
		deletedUser, err := cfg.db.DeleteUser(r.Context(), newUser.UserID)
		log.Printf("User %s deleted due to failed haash password", deletedUser.UserID)
		respondWithError(w, http.StatusInternalServerError, "failed to hash password", err)
	}

	authParams := database.SetPasswordForUserParams{
		UserID:   newUser.UserID,
		Provider: string(auth.ProviderEmail),
		ProviderSubject: sql.NullString{
			String: "",
			Valid:  false,
		},
		ProviderHash: sql.NullString{
			String: hashed,
			Valid:  hashed != "",
		},
	}

	authIdentity, err := cfg.db.SetPasswordForUser(r.Context(), authParams)
	if err != nil {
		deletedUser, err := cfg.db.DeleteUser(r.Context(), newUser.UserID)
		log.Printf("User %s deleted due to failed set password", deletedUser.UserID)
		respondWithError(w, http.StatusInternalServerError, "failed to set password", err)
	}

	log.Printf("User %s successfully created", authIdentity.UserID)

	respondWithJSON(w, http.StatusOK, newUser)
}
