package main

import "net/http"

func (cfg *apiConfig) handlerGetAllServices(w http.ResponseWriter, r *http.Request) {
	services, err := cfg.db.GetAllServices(r.Context())
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "couldn't get services", err)
		return
	}

	respondWithJSON(w, http.StatusOK, services)
}
