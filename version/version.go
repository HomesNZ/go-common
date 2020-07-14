package version

import (
	"fmt"
	"net/http"

	"github.com/sirupsen/logrus"
)

var (
	unknown = "unknown"

	// Version is the service version number. This is intended to be passed in as an ldflag at run/build.
	// For example: go run -ldflags="-X version.Version 123" main.go
	Version = unknown

	// runCmd is the command to run so the version number is correctly provided to the service.
	runCmd = "go run -ldflags \"-X github.com/HomesNZ/go-common/version.Version=`git rev-parse HEAD`\" main.go"
)

// Log logs the version number.
func Log() {
	if Version == unknown {
		logrus.Error("Version unknown. Please run service with:\n\t\t", runCmd)
	} else {
		logrus.Infof("Version: %s", Version)
	}
}

// Version is a http handler that prints out the version number
var Handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	v := []byte(Version)

	w.Header().Set("Content-Type", "text/plain")
	w.Header().Set("Content-Length", fmt.Sprintf("%v", len(v)))
	w.WriteHeader(http.StatusOK)
	w.Write(v)
})
