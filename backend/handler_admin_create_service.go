package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"

	"github.com/thom151/rewardsHub/internal/database"
)

type ServiceCode string

const (
	ServicePhoto       ServiceCode = "photo"
	ServiceVideo       ServiceCode = "video"
	ServiceFloorPlan   ServiceCode = "floor_plan"
	ServiceDrone       ServiceCode = "drone"
	ServiceSocialReel  ServiceCode = "social_reel"
	ServiceVirtualTour ServiceCode = "virtual_tour"
	ServiceEditing     ServiceCode = "editing"
	ServiceOther       ServiceCode = "other"
)

func (c ServiceCode) Valid() bool {
	switch c {
	case ServicePhoto, ServiceVideo, ServiceFloorPlan, ServiceDrone, ServiceSocialReel, ServiceVirtualTour, ServiceEditing, ServiceOther:
		return true
	default:
		return false
	}
}

func (cfg *apiConfig) handlerAdminCreateService(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Name              string      `json:"name"`
		Code              ServiceCode `json:"code"`
		Description       string      `json:"description"`
		BasePrice         string      `json:"base_price"`
		BasePointsRewards int32       `json:"base_points_rewards"`
	}

	user, ok := r.Context().Value(userKey).(database.User)
	if !ok {
		respondWithError(w, http.StatusUnauthorized, "unauthorized", nil)
		return
	}

	decoder := json.NewDecoder(r.Body)
	var serviceParams parameters
	err := decoder.Decode(&serviceParams)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "couldn't decode service params", err)
		return
	}

	if serviceParams.Name == "" {
		respondWithError(w, http.StatusBadRequest, "missing name", err)
		return
	}

	if serviceParams.Description == "" {
		respondWithError(w, http.StatusBadRequest, "missing description", err)
		return
	}

	if !serviceParams.Code.Valid() {
		respondWithError(w, http.StatusBadRequest, "invalid code", err)
		return
	}

	service, err := cfg.db.CreateService(r.Context(), database.CreateServiceParams{
		Name: serviceParams.Name,
		Code: string(serviceParams.Code),
		Description: sql.NullString{
			String: serviceParams.Description,
			Valid:  serviceParams.Description != "",
		},
		BasePrice:         serviceParams.BasePrice,
		BasePointsRewards: serviceParams.BasePointsRewards,
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "couldn't create service", err)
		return
	}

	log.Printf("Service %s, with service_id: %s was successfully made by admin %s, with user_id: %s",
		service.Name,
		service.ServiceID,
		user.FirstName,
		user.UserID,
	)

	respondWithJSON(w, http.StatusOK, service)

}
