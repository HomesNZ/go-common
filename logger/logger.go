package logger

import (
	"os"

	"github.com/HomesNZ/go-common/env"
	"github.com/Sirupsen/logrus"
)

type Option func(*logrus.Logger)

// Level option sets the log level
func Level(level string) Option {
	return func(logger *logrus.Logger) {
		level, err := logrus.ParseLevel(level)
		// No need to handle the error here, just don't update the log level
		if err == nil {
			logger.SetLevel(level)
		}
	}
}

// Init the global logger
func Init(opts ...Option) *logrus.Logger {
	logger := logrus.StandardLogger()
	// If running in the production environment, output the logs as JSON format for parsing by Logstash.
	if env.IsProd() {
		logger.Formatter = &logrus.JSONFormatter{}
	}
	logger.Out = os.Stdout
	logger.AddHook(NewBugsnagHook())
	for _, opt := range opts {
		opt(logger)
	}
	logrus.Infof("Log level: %s", logrus.GetLevel().String())
	return logger
}
