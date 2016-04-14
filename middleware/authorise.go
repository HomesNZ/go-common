package middleware

import (
	"net/http"

	"github.com/HomesNZ/data-import/config"
	"github.com/Sirupsen/logrus"
)

// KeyAuth is middleware that authorizes request based on a key URL param.
func KeyAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		params := r.URL.Query()
		if params.Get("key") != config.MaintenanceKey {
			logrus.WithField("request", r.URL).Debug("unauthorised request")
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r)
	})
}
