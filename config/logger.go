package config

import (
	"os"

	"github.com/HomesNZ/go-common/env"
	"github.com/Sirupsen/logrus"
)

// InitLogger initializes the logger by setting the log level to the env var LOG_LEVEL, or defaulting to `info`.
func InitLogger() {
	logrus.SetOutput(os.Stdout)

	level, err := logrus.ParseLevel(env.GetString("LOG_LEVEL", "info"))
	// No need to handle the error here, just don't update the log level
	if err == nil {
		logrus.SetLevel(level)
	}

	logrus.Infof("Log level: %s", logrus.GetLevel().String())
}
