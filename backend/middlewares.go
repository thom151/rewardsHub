package main

import (
	"context"
	"log"
	"net/http"

	"github.com/thom151/rewardsHub/internal/auth"
	"github.com/thom151/rewardsHub/internal/database"
)

type contextKey string

const userKey contextKey = "user"

func middlewareLog(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s %s", r.Method, r.URL.Path)
		next.ServeHTTP(w, r)
	})
}

func (cfg *apiConfig) plaformAdminOnly(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user, ok := r.Context().Value(userKey).(database.User)
		if !ok {
			respondWithError(w, http.StatusUnauthorized, "unauthorize", nil)
			return
		}

		if !user.IsAdmin {
			respondWithError(w, http.StatusForbidden, "forbidden", nil)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func (cfg *apiConfig) authMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token, err := auth.GetBearerToken(r.Header, r.Cookies())
		if err != nil {
			respondWithError(w, http.StatusUnauthorized, "cannot get acc token", err)
			return
		}

		userID, err := auth.ValidateJWT(token, cfg.tokenSecret)
		if err != nil {
			respondWithError(w, http.StatusUnauthorized, "couldn't valid access token", err)
			return
		}

		user, err := cfg.db.GetUserByID(r.Context(), userID)
		if err != nil {
			respondWithError(w, http.StatusUnauthorized, "couldn't find user", err)
			return
		}

		ctx := context.WithValue(r.Context(), userKey, user)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
