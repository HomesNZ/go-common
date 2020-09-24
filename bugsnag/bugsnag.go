package bugsnag

import (
	"github.com/HomesNZ/go-common/env"
	"github.com/HomesNZ/go-common/version"
	"github.com/sirupsen/logrus"
	"github.com/bugsnag/bugsnag-go"
)

// InitBugsnag initializes bugsnag to capture panics if BUGSNAG_API_KEY is defined. Note that because bugsnag spawns a
// new process, logs will show some initial duplicate entries.
func InitBugsnag() {
	apiKey := env.Get("BUGSNAG_API_KEY")

	if apiKey != "" {
		bugsnag.Configure(bugsnag.Configuration{
			APIKey:       apiKey,
			ReleaseStage: env.Env(),
			AppVersion:   version.Version,
		})
		logrus.Info("Bugsnag configured to capture panics")
	}
}

//Notify wraps the bugsnag.Notify call
func Notify(err error, rawData ...interface{}) {
	bugsnag.Notify(err, rawData)
}
