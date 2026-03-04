package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/thom151/rewardsHub/internal/database"
)

func (cfg *apiConfig) handlerCreateBooking(w http.ResponseWriter, r *http.Request) {
	type bookingItems struct {
		ServiceID uuid.UUID `json:"service_id"`
		Quantity  int32     `json:"quantity"`
	}

	type parameters struct {
		OrganizationID    uuid.UUID      `json:"organization_id"`
		PropertyID        uuid.UUID      `json:"property_id"`
		RequestedByUserID uuid.UUID      `json:"requested_by_user_id"`
		PreferredDate     time.Time      `json:"preferred_date"`
		ScheduleStart     time.Time      `json:"schedule_start"`
		ScheduleEnd       time.Time      `json:"schedule_end"`
		BookingItems      []bookingItems `json:"booking_items"`
	}

	user, ok := r.Context().Value(userKey).(database.User)
	if !ok {
		respondWithError(w, http.StatusForbidden, "user not found", nil)
		return
	}

	//to get organization
	membership, err := cfg.db.GetOrgMembershipFromUserID(r.Context(), user.UserID)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "couldn't find membership", err)
		return
	}

	decoder := json.NewDecoder(r.Body)
	var b parameters
	err = decoder.Decode(&b)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "couldn't decode parameters", err)
		return
	}

	if b.PropertyID == uuid.Nil || b.PreferredDate.IsZero() || b.ScheduleStart.IsZero() || b.ScheduleEnd.IsZero() {
		respondWithError(w, http.StatusBadRequest, "missing user input", nil)
		return
	}

	organizationID := membership.OrganizationID
	requestedByUser := user.UserID
	if user.IsAdmin && b.OrganizationID != uuid.Nil && b.RequestedByUserID != uuid.Nil {
		organizationID = b.OrganizationID
		requestedByUser = b.RequestedByUserID
		_, err := cfg.db.IsUserMemberofOrg(r.Context(), database.IsUserMemberofOrgParams{
			UserID:         requestedByUser,
			OrganizationID: organizationID,
		})
		if err != nil {
			respondWithError(w, http.StatusBadRequest, "agent not a member of organization", err)
			return
		}
	} else {
		respondWithError(w, http.StatusBadRequest, "admin please set requested_by_user_id and organization_id", nil)
		return
	}

	propertyID := b.PropertyID
	if !user.IsAdmin {
		property, err := cfg.db.GetPropertiesOfUser(r.Context(), database.GetPropertiesOfUserParams{
			PropertyID:      b.PropertyID,
			CreatedByUserID: user.UserID,
			OrganizationID:  organizationID,
		})
		if err != nil {
			respondWithError(w, http.StatusBadRequest, "couldn't get properties", err)
			return
		}

		propertyID = property.PropertyID
	}

	if !b.ScheduleEnd.After(b.ScheduleStart) {
		respondWithError(w, http.StatusBadRequest, "schedule_end must be after schedule_start", nil)
		return
	}

	booking, err := cfg.db.CreateBooking(r.Context(), database.CreateBookingParams{
		OrganizationID:    organizationID,
		RequestedByUserID: requestedByUser,
		PropertyID:        propertyID,
		PreferredDate:     b.PreferredDate,
		ScheduleStart:     b.ScheduleStart,
		ScheduleEnd:       b.ScheduleEnd,
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "couldn't create booking", err)
		return
	}

	var services []string
	for _, bookingItem := range b.BookingItems {
		service, err := cfg.db.GetService(r.Context(), bookingItem.ServiceID)
		if err != nil {
			respondWithError(w, http.StatusBadRequest, "service cannot found", err)
			return
		}
		bItem, err := cfg.db.CreateBookingItem(r.Context(), database.CreateBookingItemParams{
			BookingID:      booking.BookingID,
			ServiceID:      service.ServiceID,
			Quantity:       bookingItem.Quantity,
			UnitPriceCents: bookingItem.Quantity * service.BasePriceCents,
			PointsAward:    bookingItem.Quantity * service.BasePointsRewards,
		})
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, "couldn't add booking item", err)
			return
		}
		log.Printf("%s item %s was added to booking_id %s", bItem.BookingItemID, service.Code, booking.BookingID)
		services = append(services, service.Code)
	}

	respondWithJSON(w, http.StatusOK, booking)

}
