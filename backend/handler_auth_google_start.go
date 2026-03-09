package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/thom151/rewardsHub/internal/database"
)

func (cfg *apiConfig) handlerListEvents(w http.ResponseWriter, r *http.Request) {
	type params struct {
		Start      time.Time `json:"start"`
		End        time.Time `json:"end"`
		CalendarID string    `json:"calendar_id"`
	}

	decoder := json.NewDecoder(r.Body)
	var eParams params
	err := decoder.Decode(&eParams)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "couldn't decode params", err)
		return
	}

}

func (cfg *apiConfig) handlerAuthGoogleStart(w http.ResponseWriter, r *http.Request) {

	state := uuid.NewString()

	user, ok := r.Context().Value(userKey).(database.User)
	if !ok {
		respondWithError(w, http.StatusBadRequest, "couldn't get user", nil)
		return
	}
	oauth_state, err := cfg.db.CreateOAuthState(r.Context(), database.CreateOAuthStateParams{
		State:     state,
		ExpiresAt: time.Now().Add(time.Minute * 10),
		UserID:    user.UserID,
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "couldn't create state", err)
		return
	}

	log.Printf("a state for user with id %s has been made", oauth_state.UserID)
	http.SetCookie(w, &http.Cookie{
		Name:     "google_oauth_state",
		Value:    state,
		Path:     "/",
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
		MaxAge:   600, // 10 minutes
	})

	q := url.Values{}
	q.Set("client_id", cfg.googleClientID)
	q.Set("redirect_uri", cfg.googleRedirectUri)
	q.Set("response_type", "code")
	q.Set("scope", "https://www.googleapis.com/auth/calendar")
	q.Set("access_type", "offline")
	q.Set("prompt", "consent")
	q.Set("state", state)

	redirect := "https://accounts.google.com/o/oauth2/v2/auth?" + q.Encode()
	http.Redirect(w, r, redirect, http.StatusFound)
}

func (cfg *apiConfig) handlerAuthGoogleCallback(w http.ResponseWriter, r *http.Request) {

	returnedState := r.URL.Query().Get("state")
	if returnedState == "" {
		http.Error(w, "missing state", http.StatusBadRequest)
		return
	}

	cookie, err := r.Cookie("google_oauth_state")
	if err != nil {
		http.Error(w, "missing oauth state cookie", http.StatusBadRequest)
		return
	}

	if cookie.Value != returnedState {
		http.Error(w, "invalid oauth state", http.StatusBadRequest)
		return
	}

	oauth_state, err := cfg.db.GetUserFromState(r.Context(), returnedState)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "couldn't get state", err)
		return
	}

	user, err := cfg.db.GetUserByID(r.Context(), oauth_state.UserID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "couldn't get user from state", err)
		return
	}

	tokenUrl := "https://oauth2.googleapis.com/token"
	code := r.URL.Query().Get("code")
	if code == "" {
		http.Error(w, "missing code", http.StatusBadRequest)
		return
	}

	form := url.Values{}
	form.Set("client_id", cfg.googleClientID)         // or os.Getenv("GOOGLE_CLIENT_ID")
	form.Set("client_secret", cfg.googleClientSecret) // or os.Getenv("GOOGLE_CLIENT_SECRET")
	form.Set("redirect_uri", cfg.googleRedirectUri)   // must match exactly what you used in /start
	form.Set("grant_type", "authorization_code")
	form.Set("code", code)

	req, err := http.NewRequestWithContext(r.Context(), "POST", tokenUrl, strings.NewReader(form.Encode()))
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "coudln't get tokens", err)
		return
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Accept", "application/json")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		respondWithError(w, http.StatusBadGateway, "token request failed", err)
		return
	}

	defer res.Body.Close()

	body, _ := io.ReadAll(res.Body)
	if res.StatusCode < 200 || res.StatusCode >= 300 {
		http.Error(w, fmt.Sprintf("google token error (%d): %s", res.StatusCode, string(body)), http.StatusBadRequest)
		return
	}

	type tokenResp struct {
		AccessToken  string `json:"access_token"`
		ExpiresIn    int    `json:"expires_in"`
		TokenType    string `json:"token_type"`
		Scope        string `json:"scope"`
		RefreshToken string `json:"refresh_token"`
	}

	var tr tokenResp
	if err := json.Unmarshal(body, &tr); err != nil {
		respondWithError(w, http.StatusBadRequest, "failed to parse token response", err)
		return
	}
	if tr.AccessToken == "" {
		http.Error(w, "no access_token returned", http.StatusBadRequest)
		return
	}

	expiresAt := time.Now().Add(time.Duration(tr.ExpiresIn) * time.Second)
	_, err = cfg.db.CreateOauthConnection(r.Context(), database.CreateOauthConnectionParams{
		UserID:   user.UserID,
		Service:  "calendar",
		Provider: "google",
		ProviderEmail: sql.NullString{
			String: user.Email,
			Valid:  user.Email != "",
		},
		RefreshToken: tr.RefreshToken,
		AccessToken: sql.NullString{
			String: tr.AccessToken,
			Valid:  tr.AccessToken != "",
		},
		ExpiresAt: sql.NullTime{
			Time:  expiresAt,
			Valid: true,
		},
	})

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "couldn't save oauth connection", err)
		return
	}

	log.Printf("google calendar connection successfull for user %s with user_id: %s.\n", user.Email, user.UserID)
	respondWithJSON(w, http.StatusOK, tr)

}
