package gCalendar

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
)

func RefreshGoogleAccessToken(
	ctx context.Context,
	clientID string,
	clientSecret string,
	refreshToken string,
) (string, int, error) {
	tokenURL := "https://oauth2.googleapis.com/token"

	form := url.Values{}
	form.Set("client_id", clientID)
	form.Set("client_secret", clientSecret)
	form.Set("refresh_token", refreshToken)
	form.Set("grant_type", "refresh_token")

	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		tokenURL,
		strings.NewReader(form.Encode()),
	)
	if err != nil {
		return "", 0, err
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Accept", "application/json")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", 0, err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return "", 0, err
	}

	if res.StatusCode < 200 || res.StatusCode >= 300 {
		return "", 0, fmt.Errorf("google refresh error (%d): %s", res.StatusCode, string(body))
	}

	type refreshResp struct {
		AccessToken string `json:"access_token"`
		ExpiresIn   int    `json:"expires_in"`
		TokenType   string `json:"token_type"`
		Scope       string `json:"scope"`
	}

	var rr refreshResp
	if err := json.Unmarshal(body, &rr); err != nil {
		return "", 0, err
	}
	if rr.AccessToken == "" {
		return "", 0, fmt.Errorf("no access_token returned")
	}

	return rr.AccessToken, rr.ExpiresIn, nil
}
