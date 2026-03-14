package main

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/thom151/rewardsHub/internal/database"
)

func (cfg *apiConfig) handlerAdminConfirmRedemption(w http.ResponseWriter, r *http.Request) {
	redemptionIDFromPath := r.PathValue("redemption_id")
	if redemptionIDFromPath == "" {
		respondWithError(w, http.StatusBadRequest, "missing redemption path", nil)
		return
	}

	redemptionID, err := uuid.Parse(redemptionIDFromPath)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid id", err)
		return
	}

	redemption, err := cfg.db.GetRedemptionFromID(r.Context(), redemptionID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "couldn't get redemption", err)
		return
	}

	user, ok := r.Context().Value(userKey).(database.User)
	if !ok {
		respondWithError(w, http.StatusUnauthorized, "unauthorized", nil)
		return
	}

	if !user.IsAdmin {
		respondWithError(w, http.StatusForbidden, "forbidden", nil)
		return
	}

	approvedRedemption, err := cfg.db.ApproveRedemption(r.Context(), redemption.RedemptionID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "couldn't approve redemption", err)
		return
	}

	//TODO: log the necessary info
	respondWithJSON(w, http.StatusOK, approvedRedemption)

}
