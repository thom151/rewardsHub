package main

import (
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
	"github.com/thom151/rewardsHub/internal/database"
)

func (cfg *apiConfig) handlerCreateVideoMeta(w http.ResponseWriter, r *http.Request) {

	user, ok := r.Context().Value(userKey).(database.User)
	if !ok {
		respondWithError(w, http.StatusUnauthorized, "couldn't find user", nil)
		return
	}

	type parameters struct {
		JobID uuid.UUID `json:"job_id"`
	}

	decoder := json.NewDecoder(r.Body)
	vMetaParams := parameters{}
	err := decoder.Decode(&vMetaParams)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "coudln't decode parameters", err)
		return
	}

	video, err := cfg.db.CreateJobAsset(r.Context(), database.CreateJobAssetParams{
		JobID: vMetaParams.JobID,
		CreatedByUserID: uuid.NullUUID{
			UUID:  user.UserID,
			Valid: user.UserID != uuid.UUID{},
		},
	})

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "couldn't create video", err)
		return
	}

	respondWithJSON(w, http.StatusCreated, video)

}
