package middleware

import (
	"fmt"
	"net/http"
	"time"

	"github.com/Sirupsen/logrus"
)

// Logger is a middleware handler that logs the request as it goes in and the response as it goes out.
type Logger struct {
	// Logger is the log.Logger instance used to log messages with the Logger middleware
	Logger *logrus.Logger
	// Name is the name of the application as recorded in latency metrics
	Name string
}

// NewLogger returns a new *Logger
func NewLogger() *Logger {
	return NewCustomLogger(logrus.InfoLevel, &logrus.TextFormatter{}, "web")
}

// NewCustomLogger builds a *Logger with the given level and formatter
func NewCustomLogger(level logrus.Level, formatter logrus.Formatter, name string) *Logger {
	log := logrus.New()
	log.Level = level
	log.Formatter = formatter

	return &Logger{Logger: log, Name: name}
}

// Log is middleware that logs both the request and the response.
func (l *Logger) Log(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Wrap the response writer in a LoggedResponseWriter so we can store additional info on the response.
		rw := NewLoggedResponseWriter(w)

		start := time.Now()

		// Try to get the real IP
		remoteAddr := r.RemoteAddr
		if realIP := r.Header.Get("X-Real-IP"); realIP != "" {
			remoteAddr = realIP
		}

		entry := l.Logger.WithFields(logrus.Fields{
			"request": r.RequestURI,
			"method":  r.Method,
			"remote":  remoteAddr,
		})

		if reqID := r.Header.Get("X-Request-Id"); reqID != "" {
			entry = entry.WithField("request_id", reqID)
		}
		entry.Info("started handling request")

		next.ServeHTTP(w, r)

		latency := time.Since(start)
		entry.WithFields(logrus.Fields{
			"status":      rw.Status(),
			"text_status": http.StatusText(rw.Status()),
			"took":        latency,
			fmt.Sprintf("measure#%s.latency", l.Name): latency.Nanoseconds(),
		}).Info("completed handling request")
	})
}
