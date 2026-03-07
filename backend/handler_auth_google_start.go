package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
)

func (cfg *apiConfig) handlerAuthGoogleStart(w http.ResponseWriter, r *http.Request) {

	redirect := "https://accounts.google.com/o/oauth2/v2/auth?" +
		"client_id=" + cfg.googleClientID +
		"&redirect_uri=" + cfg.googleRedirectUri +
		"&response_type=code" +
		"&scope=https://www.googleapis.com/auth/calendar.readonly"

	http.Redirect(w, r, redirect, http.StatusFound)
}

func (cfg *apiConfig) handlerAuthGoogleCallback(w http.ResponseWriter, r *http.Request) {

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

	respondWithJSON(w, http.StatusOK, tr)

}
