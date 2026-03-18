package main

import (
	"database/sql"
	"fmt"
	"log"
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

	orgPoints, err := cfg.db.GetOrganizationPointsBalance(r.Context(), redemption.OrganizationID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "couldn't get point", err)
	}

	if int32(orgPoints) < redemption.PointsCost {
		respondWithError(w, http.StatusBadRequest, "not enough balance", err)
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

	pointLedger, err := cfg.db.CreatePointLedger(r.Context(), database.CreatePointLedgerParams{
		OrganizationID: redemption.OrganizationID,
		UserID: uuid.NullUUID{
			UUID:  redemption.RequestedByUserID,
			Valid: true,
		},
		EventType: "reward_redeemed",
		EventID: uuid.NullUUID{
			UUID:  redemption.RedemptionID,
			Valid: true,
		},
		PointsDelta: -redemption.PointsCost,
		Description: sql.NullString{
			String: fmt.Sprintf("user %s from organization %s redeemed service", redemption.RequestedByUserID, redemption.OrganizationID),
			Valid:  true,
		},
	})

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "couldn't redeem point", err)
		return
	}

	log.Printf("user %s approved agent %s to redeem %d points", user.Email, redemption.RequestedByUserID, pointLedger.PointsDelta)

	//TODO: log the necessary info
	respondWithJSON(w, http.StatusOK, approvedRedemption)

}
