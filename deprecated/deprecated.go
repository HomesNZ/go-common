package deprecated

import (
	"net/http"

	"github.com/sirupsen/logrus"
)

// Deprecated is a middleware which logs a message when a Deprecated endpoint is used. It also attaches a header to the request notifying the consumer of the deprecation
// Deprecated is intended to be included in the alice.Chain of any deprecated endpoint
func Deprecated(next http.Handler) http.Handler {
	contextLogger := logrus.WithField("Middleware", "Deprecated")
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		contextLogger.Info("Request recieved", r.URL.String())
		w.Header().Set("Endpoint-Deprecated", "True")
		next.ServeHTTP(w, r)
	})
}
