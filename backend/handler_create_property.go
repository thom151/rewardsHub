package main

import (
	"database/sql"
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
	"github.com/thom151/rewardsHub/internal/database"
)

func (cfg *apiConfig) handlerCreateProperty(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		AddressLine1 string `json:"address_line1"`
		AddressLine2 string `json:"address_lin2"`
		City         string `json:"city"`
		StateRegion  string `json:"state_region"`
		PostalCode   string `json:"postal_code"`
		ListingUrl   string `json:"listing"`
	}

	user, ok := r.Context().Value(userKey).(database.User)
	if !ok {
		respondWithError(w, http.StatusUnauthorized, "unathorized", nil)
		return
	}

	decoder := json.NewDecoder(r.Body)
	var p parameters
	err := decoder.Decode(&p)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "couldn't decode property parameters", err)
		return
	}

	membership, err := cfg.db.GetOrgMembershipFromUserID(r.Context(), user.UserID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "couldn't get membership", err)
		return
	}

	property, err := cfg.db.CreateProperty(r.Context(), database.CreatePropertyParams{
		OrganizationID: membership.OrganizationID,
		CreatedByUserID: uuid.NullUUID{
			UUID:  membership.UserID,
			Valid: true,
		},
		AddressLine1: p.AddressLine1,
		AddressLine2: sql.NullString{
			String: p.AddressLine2,
			Valid:  p.AddressLine2 != "",
		},
		City:        p.City,
		StateRegion: p.StateRegion,
		PostalCode:  p.PostalCode,
		ListingUrl: sql.NullString{
			String: p.ListingUrl,
			Valid:  p.ListingUrl != "",
		},
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "couldn't create property", err)
		return
	}

	respondWithJSON(w, http.StatusOK, property)
}
