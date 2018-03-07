package logger

import (
	"fmt"
	"net/http"
	"time"

	"github.com/Sirupsen/logrus"
)

// name is the name of the application as recorded in latency metrics
const name = "web"

func Middleware(logger *logrus.Entry) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			logged := loggedResponseWriter{
				ResponseWriter: w,
				Status:         200,
			}
			req := r.WithContext(populateRequestContext(r.Context(), r))

			defer func(begin time.Time) {
				latency := time.Since(begin)
				logger.WithFields(logrus.Fields{
					"request":     r.RequestURI,
					"method":      r.Method,
					"status":      logged.Status,
					"text_status": http.StatusText(logged.Status),
					"took":        latency,
					fmt.Sprintf("measure#%s.latency", name): latency.Nanoseconds(),
				}).
					WithFields(FieldsFromContext(r.Context())).
					Info("Handled request")
			}(time.Now())

			if f, ok := w.(http.Flusher); ok {
				loggedFlusher := loggedResponseWriteFlusher{
					loggedResponseWriter: logged,
					Flusher:              f,
				}
				next.ServeHTTP(&loggedFlusher, req)
			} else {
				next.ServeHTTP(&logged, req)
			}
		})
	}
}

// loggedResponseWriter is a custom ResponseWriter that wraps http.ResponseWriter and also includes a status field for
// logging.
type loggedResponseWriter struct {
	http.ResponseWriter
	Status int
}

func (w *loggedResponseWriter) WriteHeader(statusCode int) {
	w.Status = statusCode
	w.ResponseWriter.WriteHeader(statusCode)
}

type loggedResponseWriteFlusher struct {
	loggedResponseWriter
	http.Flusher
}
