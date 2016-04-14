package middleware

import "net/http"

// LoggedResponseWriter is a custom ResponseWriter that wraps http.ResponseWriter and also includes a status field for
// logging.
type LoggedResponseWriter struct {
	status int
	http.ResponseWriter
}

// NewLoggedResponseWriter returns a new LoggedResponseWriter with the default settings.
func NewLoggedResponseWriter(res http.ResponseWriter) *LoggedResponseWriter {
	// Default the status code to 200
	return &LoggedResponseWriter{200, res}
}

// Status returns the response status
func (w LoggedResponseWriter) Status() int {
	return w.status
}

// Header satisfies the http.ResponseWriter interface
func (w LoggedResponseWriter) Header() http.Header {
	return w.ResponseWriter.Header()
}

// Write satisfies the http.ResponseWriter interface
func (w LoggedResponseWriter) Write(data []byte) (int, error) {
	return w.ResponseWriter.Write(data)
}

// WriteHeader satisfies the http.ResponseWriter interface
func (w LoggedResponseWriter) WriteHeader(statusCode int) {
	// Store the status code
	w.status = statusCode

	// Write the status code to the wrapped ResponseWriter
	w.ResponseWriter.WriteHeader(statusCode)
}
