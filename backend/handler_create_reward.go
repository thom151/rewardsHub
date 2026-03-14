package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"

	"github.com/google/uuid"
	"github.com/thom151/rewardsHub/internal/database"
)

func (cfg *apiConfig) handlerCreateReward(w http.ResponseWriter, r *http.Request) {
	type params struct {
		Code        string    `json:"code"`
		Name        string    `json:"name"`
		Description string    `json:"description"`
		ServiceID   uuid.UUID `json:"service_id"`
		PointsCost  int32     `json:"points_cost"`
	}

	user, ok := r.Context().Value(userKey).(database.User)
	if !ok {
		respondWithError(w, http.StatusUnauthorized, "couldn't find user", nil)
		return
	}

	if !user.IsAdmin {
		respondWithError(w, http.StatusForbidden, "forbidden", nil)
		return
	}

	decoder := json.NewDecoder(r.Body)
	var rewardsParam params
	err := decoder.Decode(&rewardsParam)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "couldn't decode params", err)
		return
	}
	reward, err := cfg.db.CreateReward(r.Context(), database.CreateRewardParams{
		Code: rewardsParam.Code,
		Name: rewardsParam.Name,
		Description: sql.NullString{
			String: rewardsParam.Description,
			Valid:  rewardsParam.Description != "",
		},
		ServiceID:  rewardsParam.ServiceID,
		PointsCost: rewardsParam.PointsCost,
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "couldn't create reward", err)
		return
	}

	log.Printf("Reward %s with rewar_id %s was successfully made by admin user % with user_id %s", reward.Name, reward.RewardID, user.Email, user.UserID)
	respondWithJSON(w, http.StatusOK, reward)
}
