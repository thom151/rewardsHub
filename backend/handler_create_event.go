package main

import (
	"net/http"

	"github.com/thom151/rewardsHub/internal/database"
)

func (cfg *apiConfig) handlerCreateEvent(w http.ResponseWriter, r *http.Request) {
	_, ok := r.Context().Value(userKey).(database.User)
	if !ok {
		respondWithError(w, http.StatusUnauthorized, "couldn't get user", nil)
		return
	}

}
