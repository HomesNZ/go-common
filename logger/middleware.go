package logger

import (
	"fmt"
	"net/http"
	"time"

	"github.com/sirupsen/logrus"
)

// name is the name of the application as recorded in latency metrics
const name = "web"

func getFirst(h http.Header, names ...string) string {
	for _, name := range names {
		if v := h.Get(name); v != "" {
			return v
		}
	}
	return ""
}

func Middleware(logger *logrus.Entry) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			var status int
			defer func(begin time.Time) {
				// Try to get the real IP
				remoteAddr := r.RemoteAddr
				if realIP := getFirst(r.Header, "X-Real-IP", "X-Forwarded-For"); realIP != "" {
					remoteAddr = realIP
				}
				entry := logger.WithFields(logrus.Fields{
					"request": r.RequestURI,
					"method":  r.Method,
					"remote":  remoteAddr,
				})
				if reqID := getFirst(r.Header, "X-Request-Id", "X-Amzn-Trace-Id"); reqID != "" {
					entry = entry.WithField("request_id", reqID)
				}
				latency := time.Since(begin)
				entry.WithFields(logrus.Fields{
					"status":                                status,
					"text_status":                           http.StatusText(status),
					"took":                                  latency,
					fmt.Sprintf("measure#%s.latency", name): latency.Nanoseconds(),
				}).Debug("Handled request")
			}(time.Now())

			var logged http.ResponseWriter = &loggedResponseWriter{
				ResponseWriter: w,
				Status:         &status,
			}
			if f, ok := w.(http.Flusher); ok {
				logged = responseWriteFlusher{
					ResponseWriter: logged,
					Flusher:        f,
				}
			}
			next.ServeHTTP(logged, r)
		})
	}
}

// loggedResponseWriter is a custom ResponseWriter that wraps http.ResponseWriter and also includes a status field for
// logging.
type loggedResponseWriter struct {
	http.ResponseWriter
	Status *int
}

func (w *loggedResponseWriter) WriteHeader(statusCode int) {
	*w.Status = statusCode
	w.ResponseWriter.WriteHeader(statusCode)
}

type responseWriteFlusher struct {
	http.ResponseWriter
	http.Flusher
}
