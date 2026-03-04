package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"path"
	"time"

	"github.com/google/uuid"
	"github.com/thom151/rewardsHub/internal/database"
	"github.com/thom151/rewardsHub/internal/dropbox"
)

func (cfg *apiConfig) handlerAdminConfirmBooking(w http.ResponseWriter, r *http.Request) {
	bookingIDFromPath := r.PathValue("booking_id")
	if bookingIDFromPath == "" {
		respondWithError(w, http.StatusBadRequest, "missing booking id path", nil)
		return
	}

	bookingUUID, err := uuid.Parse(bookingIDFromPath)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid id", err)
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

	booking, err := cfg.db.GetBooking(r.Context(), bookingUUID)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "booking not found", err)
		return
	}

	_, err = cfg.db.ConfirmBooking(r.Context(), booking.BookingID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "couldn't confirm booking", err)
		return
	}

	organization, err := cfg.db.GetOrganizationFromID(r.Context(), booking.OrganizationID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "couldn't find organization", err)
		return
	}

	property, err := cfg.db.GetPropertyFromAdmin(r.Context(), booking.PropertyID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "couldn't find property", err)
		return
	}

	bookingItems, err := cfg.db.GetBookingItems(r.Context(), booking.BookingID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "couldn't get booking items", err)
	}

	var services []string
	for _, item := range bookingItems {
		service, err := cfg.db.GetService(r.Context(), item.ServiceID)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, "couldn't get service", err)
			return
		}
		services = append(services, service.Code)
	}

	userAgent, err := cfg.db.GetUserByID(r.Context(), booking.RequestedByUserID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "user agent can't be found", err)
		return
	}

	if time.Now().After(cfg.dropboxAccTokenExpiresAt) {
		newAccTok, err := dropbox.GetNewAccessToken(cfg.dropboxRefreshToken, cfg.dropboxClientID, cfg.dropboxClientSecret)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, "error getting acces token", err)
			return
		}
		cfg.dropboxAccToken = newAccTok.AccessToken
		fmt.Printf("acc token: %s", cfg.dropboxAccToken)
	}

	err = createDropboxBookingFolder(r.Context(), cfg.dropboxAccToken, organization.Name, userAgent.Email, property.AddressLine1, services)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "dropbox folder couldn't be created", err)
		return
	}

	log.Printf("folders created in dropbox %s/%s/%s/", organization.Name, userAgent.Email, property.AddressLine1)

	job, err := cfg.db.CreateJob(r.Context(), database.CreateJobParams{
		BookingID: booking.BookingID,
		AssignedToUserID: uuid.NullUUID{
			UUID:  user.UserID,
			Valid: true,
		},
	})

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "couldn't create job", err)
		return
	}

	log.Printf("job was successfully created by %s with id with %s.  ", user.FirstName, user.UserID)
	respondWithJSON(w, http.StatusOK, job)
}

const BaseDropboxFolder = "leadway-rewards"

func createDropboxBookingFolder(ctx context.Context, accToken string, organization, agent, address string, services []string) error {
	for i, service := range services {
		folderName := fmt.Sprintf("%02d_%s", i+1, service)
		path := path.Join(BaseDropboxFolder, organization, agent, address, folderName)
		dropboxPath := "/" + path
		err := dropbox.CreateDropboxFolder(ctx, dropboxPath, accToken)
		if err != nil {
			return err
		}
	}

	return nil
}
