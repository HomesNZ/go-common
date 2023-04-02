package buggsnag

import (
	"github.com/bugsnag/bugsnag-go/v2"
	"github.com/bugsnag/bugsnag-go/v2/device"
	"github.com/gin-gonic/gin"
)

const frameworkName string = "Gin"

func Notify(rawData ...interface{}) gin.HandlerFunc {
	device.AddVersion(frameworkName, gin.Version)
	state := bugsnag.HandledState{
		SeverityReason:   bugsnag.SeverityReasonUnhandledMiddlewareError,
		OriginalSeverity: bugsnag.SeverityError,
		Unhandled:        true,
		Framework:        frameworkName,
	}
	rawData = append(rawData, state)
	return func(c *gin.Context) {
		r := c.Copy().Request
		notifier := bugsnag.New(append(rawData, r)...)
		ctx := bugsnag.AttachRequestData(r.Context(), r)
		if notifier.Config.IsAutoCaptureSessions() {
			ctx = bugsnag.StartSession(ctx)
		}
		c.Request = r.WithContext(ctx)

		notifier.FlushSessionsOnRepanic(false)
		defer notifier.AutoNotify(ctx)
		c.Next()
	}
}
