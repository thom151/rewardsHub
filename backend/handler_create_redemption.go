package main

import (
	"log"
	"net/http"

	"github.com/google/uuid"
	"github.com/thom151/rewardsHub/internal/database"
)

// users will be using this
func (cfg *apiConfig) handlerCreateRedemption(w http.ResponseWriter, r *http.Request) {
	rewardIDFromPath := r.PathValue("reward_id")
	if rewardIDFromPath == "" {
		respondWithError(w, http.StatusBadRequest, "reward id missing in path", nil)
		return
	}

	rewardID, err := uuid.Parse(rewardIDFromPath)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "couldn't parse reward id", err)
		return
	}

	reward, err := cfg.db.GetRewardFromID(r.Context(), rewardID)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "couldn't find reward", err)
		return
	}

	bookingIDFromPath := r.PathValue("booking_id")
	if bookingIDFromPath == "" {
		respondWithError(w, http.StatusBadRequest, "reward id missing in path", nil)
		return
	}

	bookingID, err := uuid.Parse(bookingIDFromPath)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "couldn't parse reward id", err)
		return
	}

	booking, err := cfg.db.GetBooking(r.Context(), bookingID)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "couldn't find reward", err)
		return
	}

	user, ok := r.Context().Value(userKey).(database.User)
	if !ok {
		respondWithError(w, http.StatusUnauthorized, "couldn't get user", nil)
		return
	}

	orgMembership, err := cfg.db.GetOrgMembershipFromUserID(r.Context(), user.UserID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "couldn't get membership", err)
		return
	}

	redemption, err := cfg.db.CreateRedemption(r.Context(), database.CreateRedemptionParams{
		OrganizationID:    orgMembership.OrganizationID,
		RequestedByUserID: user.UserID,
		RewardID:          reward.RewardID,
		PointsCost:        reward.PointsCost,
		AppliedBookingID:  booking.BookingID,
	})

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "couldn't create redemption", err)
		return
	}

	log.Printf("redemption of service %s with redemption_id of %s was made by user %s with user_id %s", reward.ServiceID, redemption.RedemptionID, user.Email, user.UserID)

	respondWithJSON(w, http.StatusOK, redemption)

}
