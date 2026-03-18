package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/google/uuid"
	"github.com/thom151/rewardsHub/internal/database"
	"github.com/thom151/rewardsHub/internal/dropbox"
)

type AssetType string

const (
	AssetPhoto     AssetType = "photo"
	AssetVideo     AssetType = "video"
	AssetFloorPlan AssetType = "floorplan"
)

func (cfg *apiConfig) handlerSyncJobAssets(w http.ResponseWriter, r *http.Request) {

	user, ok := r.Context().Value(userKey).(database.User)
	if !ok {
		respondWithError(w, http.StatusUnauthorized, "couldn't get user", nil)
		return
	}

	jobIDString := r.PathValue("job_id")
	if jobIDString == "" {
		respondWithError(w, http.StatusBadRequest, "missing job id in path", nil)
		return
	}

	jobID, err := uuid.Parse(jobIDString)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "couldn't parse job id", err)
		return
	}

	job, err := cfg.db.GetJob(r.Context(), jobID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "couldn't get job", err)
		return
	}
	entries, err := dropbox.ListDropboxFolder(r.Context(), job.DropboxFolderPath.String, cfg.dropboxAccToken)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "couldn't list folder", err)
		return
	}

	tx, err := cfg.conn.BeginTx(r.Context(), nil)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "couldn't begin transaction", err)
		return
	}
	defer tx.Rollback()

	qtx := cfg.db.WithTx(tx)

	sortOrder := 0
	syncedCount := 0
	for _, entry := range entries {
		if entry.Tag == "file" {
			assetType, mimeType, err := getAssetType(entry.Name)
			if err != nil {
				continue
			}
			sortOrder += 1
			jobAsset, err := qtx.CreateJobAsset(r.Context(), database.CreateJobAssetParams{
				JobID: job.JobID,
				CreatedByUserID: uuid.NullUUID{
					UUID:  user.UserID,
					Valid: true,
				},
				AssetType:   string(assetType),
				MimeType:    mimeType,
				FileName:    entry.Name,
				SizeBytes:   int64(entry.Size),
				DropboxPath: entry.PathDisplay,
				DropboxFileID: sql.NullString{
					String: entry.ID,
					Valid:  entry.ID != "",
				},
				IsPrimary: false,
				ChecksumSha256: sql.NullString{
					String: "",
					Valid:  false,
				},
				SortOrder: int32(sortOrder),
			})

			if err != nil {
				log.Printf("couldn't create job asset for %s: %v", entry.PathDisplay, err)
				continue
			}

			syncedCount++
			log.Printf("created job asset for %s with asset_id %s", entry.PathDisplay, jobAsset.AssetID)
		}
	}

	if syncedCount == 0 {
		respondWithError(w, http.StatusBadRequest, "no valid assets were synced", nil)
		return
	}

	jobCompleted, err := qtx.CompleteJob(r.Context(), job.JobID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "couldnt set job completed", err)
		return
	}

	booking, err := cfg.db.GetBooking(r.Context(), job.BookingID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "coulnd't get booking", err)
		return
	}
	bookingItems, err := cfg.db.GetBookingItems(r.Context(), booking.BookingID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "couldn't get booking item", err)
		return
	}

	points := 0
	for _, item := range bookingItems {
		service, err := cfg.db.GetService(r.Context(), item.ServiceID)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, "couldn't get service", err)
			return
		}

		points += int(service.BasePointsRewards)
	}

	pointLedger, err := qtx.CreatePointLedger(r.Context(), database.CreatePointLedgerParams{
		OrganizationID: booking.OrganizationID,
		UserID: uuid.NullUUID{
			UUID:  booking.RequestedByUserID,
			Valid: true,
		},
		EventType: "job_completed",
		EventID: uuid.NullUUID{
			UUID:  job.JobID,
			Valid: true,
		},
		PointsDelta: int32(points),
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "couldn't create point ledger", err)
		return
	}

	type resp struct {
		Points      database.PointLedger `json:"points"`
		Job         database.Job         `json:"job"`
		SyncedCount int                  `json:"synced_count"`
	}

	respondWithJSON(w, http.StatusOK, resp{
		Points:      pointLedger,
		Job:         jobCompleted,
		SyncedCount: syncedCount,
	})
}

func getAssetType(fileName string) (assetType AssetType, mimeType string, err error) {
	ext := strings.ToLower(filepath.Ext(fileName))

	switch ext {
	case ".jpg", ".jpeg":
		return AssetPhoto, "image/jpeg", nil
	case ".png":
		return AssetPhoto, "image/png", nil
	case ".webp":
		return AssetPhoto, "image/webp", nil

	case ".mp4":
		return AssetVideo, "video/mp4", nil
	default:
		return "", "", fmt.Errorf("unsupported file type: %s", ext)
	}
}
