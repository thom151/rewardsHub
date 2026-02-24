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

func (cfg *apiConfig) handlerApproveOrganizationMembership(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		UserIDToApprove uuid.UUID `json:"user_id_to_approve"`
		OrgRole         string    `json:"org_role"`
	}
	user, ok := r.Context().Value(userKey).(database.User)
	if !ok {
		respondWithError(w, http.StatusForbidden, "forbidden", nil)
		return
	}

	membership, err := cfg.db.GetOrgMembershipFromUserID(r.Context(), user.UserID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "cannot get membership", err)
		return
	}

	if membership.OrgRole != string(OrgRoleAdmin) && !user.IsAdmin {
		respondWithError(w, http.StatusUnauthorized, "unahotrized to approve", nil)
		return
	}

	decoder := json.NewDecoder(r.Body)
	var orgMemberParams parameters
	err = decoder.Decode(&orgMemberParams)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "couldn't decode params", err)
		return
	}

	if orgMemberParams.OrgRole == "" || orgMemberParams.UserIDToApprove == uuid.Nil {
		respondWithError(w, http.StatusBadRequest, "missing parameters", err)
		return
	}

	userApproval, err := cfg.db.GetUserByID(r.Context(), orgMemberParams.UserIDToApprove)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "user not existing", err)
		return
	}

	orgMembership, err := cfg.db.CreateOrgMembership(r.Context(), database.CreateOrgMembershipParams{
		OrganizationID: membership.OrganizationID,
		UserID:         orgMemberParams.UserIDToApprove,
		OrgRole:        orgMemberParams.OrgRole,
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "couldn't approve membership", err)
		return
	}

	log.Printf("The membership of % with user_id %s, was successfully approved by org admin %s, with user_id %s",
		userApproval.FirstName,
		userApproval.UserID,
		user.FirstName,
		user.UserID,
	)

	respondWithJSON(w, http.StatusOK, orgMembership)

}
