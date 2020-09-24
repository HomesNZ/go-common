package util

import (
	"os"

	log "github.com/sirupsen/logrus"
)

// MustMkdirAll recursively makes the required directories to complete the supplied path. If an error occurs, it will be
// logged as fatal.
func MustMkdirAll(path string) {
	err := os.MkdirAll(path, 0777)
	if err != nil {
		log.Fatal(err)
	}
}
