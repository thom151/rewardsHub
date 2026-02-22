package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/thom151/rewardsHub/internal/auth"
	"github.com/thom151/rewardsHub/internal/database"
)

type User struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Email     string    `json:"email"`
	Password  string    `json:"-"`
}

func (cfg *apiConfig) handlerLogin(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Password string `json:"password"`
		Email    string `json:"email"`
	}

	type response struct {
		User
		Token        string `json:"token"`
		RefreshToken string `json:"refresh_token"`
	}

	decoder := json.NewDecoder(r.Body)
	loginParams := parameters{}
	err := decoder.Decode(&loginParams)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "couldn't decode params", err)
		return
	}

	user, err := cfg.db.GetUserByEmail(r.Context(), loginParams.Email)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "incorrect email or password", err)
		return
	}

	userIdentity, err := cfg.db.GetUserPasswordHash(r.Context(), user.UserID)
	match, err := auth.CheckPasswordHash(loginParams.Password, userIdentity.ProviderHash.String)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "incorrect email or password", err)
		return
	}

	if !match {
		respondWithError(w, http.StatusUnauthorized, "incorrect email or password", err)
		return
	}

	accessToken, err := auth.MakeJWT(user.UserID, cfg.tokenSecret, time.Hour)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "couldn't create access jwt", err)
		return
	}

	refreshToken := auth.MakeRefreshToken()

	_, err = cfg.db.CreateRefreshToken(r.Context(), database.CreateRefreshTokenParams{
		UserID:    user.UserID,
		Token:     refreshToken,
		ExpiresAt: time.Now().UTC().Add(time.Hour * 24 * 60),
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "couldn't create refresh token", err)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "acc_token",
		Value:    accessToken,
		HttpOnly: true,
		Secure:   true,
		Expires:  time.Now().Add(24 * time.Hour),
	})

	respondWithJSON(w, http.StatusOK, response{
		User: User{
			ID:        user.UserID,
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
			Email:     user.Email,
		},
		Token:        accessToken,
		RefreshToken: refreshToken,
	})

}
