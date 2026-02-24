package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/thom151/rewardsHub/internal/database"
)

type OrganizationType string

const (
	OrgTypePlatform OrganizationType = "platform"
	OrgTypeAgency   OrganizationType = "agency"
	OrgTypeInternal OrganizationType = "internal"
)

type OrgRole string

const (
	OrgRoleAdmin  OrgRole = "admin"
	OrgRoleOwner  OrgRole = "owner"
	OrgRoleStaff  OrgRole = "Staff"
	OrgRoleClient OrgRole = "client"
)

type Organization struct {
	OrganizationID uuid.UUID `json:"organization_id"`
	Name           string    `json:"name"`
	Status         string    `json:"status"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

func (cfg *apiConfig) handlerAdminCreateOrganization(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Name             string           `json:"name"`
		OrganizationType OrganizationType `json:"organization_type"`
	}

	type response struct {
		User        User
		Organizaion Organization
	}

	user, ok := r.Context().Value(userKey).(database.User)
	if !ok {
		respondWithError(w, http.StatusUnauthorized, "unauthorized", nil)
		return
	}

	decoder := json.NewDecoder(r.Body)
	var orgParams parameters
	err := decoder.Decode(&orgParams)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "couldn't decode organization parameters", err)
		return
	}

	switch orgParams.OrganizationType {
	case OrgTypePlatform:
	case OrgTypeAgency:
	case OrgTypeInternal:
	default:
		respondWithError(w, http.StatusBadRequest, "invalid organization type", nil)
		return
	}

	if orgParams.OrganizationType == OrgTypePlatform {
		respondWithError(w, http.StatusUnauthorized, "couldn't create platform organization", nil)
		return
	}

	org, err := cfg.db.CreateOrganization(r.Context(), database.CreateOrganizationParams{
		Name:             orgParams.Name,
		OrganizationType: string(orgParams.OrganizationType),
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "coudln't create organization", err)
		return
	}

	/*
		membership, err := cfg.db.CreateOrgMembership(r.Context(), database.CreateOrgMembershipParams{
			OrganizationID: org.OrganizationID,
			UserID:         user.UserID,
		})

	*/

	log.Printf("Organization %s with organization_id: %s, successfully created by %s with user_id: %s",
		org.Name,
		org.OrganizationID,
		user.FirstName,
		user.UserID,
	)

	respondWithJSON(w, http.StatusOK, response{
		User{
			ID:        user.UserID,
			Email:     user.Email,
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
		},

		Organization{
			OrganizationID: org.OrganizationID,
			Name:           org.Name,
			Status:         org.Status,
			CreatedAt:      org.CreatedAt,
			UpdatedAt:      org.UpdatedAt,
		},
	})

}
