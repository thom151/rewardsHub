package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/thom151/rewardsHub/internal/database"
	"github.com/thom151/rewardsHub/internal/gCalendar"
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

	user, ok := r.Context().Value(userKey).(database.User)
	if !ok {
		respondWithError(w, http.StatusBadRequest, "couldn'get user", nil)
		return
	}

	accessToken, err := cfg.getValidGoogleCalendarAccessToken(r.Context(), user.UserID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "couldn't get access token", err)
		return
	}

	events, err := gCalendar.GetCalendarEvents(r.Context(), accessToken, eParams.Start, eParams.End)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "couldn't get events", err)
		return
	}

	respondWithJSON(w, http.StatusOK, events)

}
