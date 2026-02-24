package middlewares

import (
	"net/http"

	"github.com/thom151/rewardsHub/internal/auth"
)

func AdminPlatformOnly(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, err := auth.GetBearerToken(r.Header, r.Cookies())
		if err != nil {
			respondWithMiddlewareError(w, http.StatusInternalServerError, "cannot get acc token", err)
			return
		}
	})
}
