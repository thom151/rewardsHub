package main

import (
	"database/sql"
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
	"github.com/thom151/rewardsHub/internal/database"
)

func (cfg *apiConfig) handlerAdminCreateReward(w http.ResponseWriter, r *http.Request) {
	type params struct {
		Description string `json:"description"`
	}

	user, ok := r.Context().Value(userKey).(database.User)

	if !ok {
		respondWithError(w, http.StatusUnauthorized, "unauthorized access", nil)
		return
	}
	if !user.IsAdmin {
		respondWithError(w, http.StatusForbidden, "forbidden", nil)
		return
	}

	decoder := json.NewDecoder(r.Body)
	var rewardParams params
	err := decoder.Decode(&rewardParams)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "couldn't decode params", err)
		return
	}

	serviceIDFromPath := r.PathValue("service_id")
	if serviceIDFromPath == "" {
		respondWithError(w, http.StatusBadRequest, "missing service_id in path", nil)
		return
	}

	serviceID, err := uuid.Parse(serviceIDFromPath)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid service id", nil)
		return
	}

	service, err := cfg.db.GetService(r.Context(), serviceID)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "service doesn't exist", err)
		return
	}

	reward, err := cfg.db.CreateReward(r.Context(), database.CreateRewardParams{
		Code: "LWR" + service.Name,
		Name: "rewards" + service.Name,
		Description: sql.NullString{
			String: rewardParams.Description,
			Valid:  rewardParams.Description != "",
		},
		ServiceID:  service.ServiceID,
		PointsCost: service.BasePointsRewards,
	})

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "couldn't create rewards", err)
		return
	}

	respondWithJSON(w, http.StatusOK, reward)

}
