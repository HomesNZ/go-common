package bugsnag

import (
	"context"
	"github.com/HomesNZ/go-common/bugsnag/config"
	bugsnag "github.com/bugsnag/bugsnag-go/v2"
	"net/http"
)

func NewFromEnv() error {
	cnf := config.NewFromEnv()
	if err := cnf.Validate(); err != nil {
		return err
	}
	bugsnag.Configure(bugsnag.Configuration{
		APIKey:       cnf.APIKey,
		ReleaseStage: cnf.Stage,
	})

	return nil
}

// Notify wraps the bugsnag.Notify call
// Usage:
// ctx := context.Background()
// ctx = bugsnag.StartSession(ctx)
// _, err := net.Listen("tcp", ":80")
//
// if err != nil {
// bugsnag.Notify(err, ctx)
// }
func Notify(err error, rawData ...interface{}) error {
	return bugsnag.Notify(err, rawData)
}

// StartSession wraps the bugsnag.StartSession call
func StartSession(ctx context.Context) context.Context {
	return bugsnag.StartSession(ctx)
}

// AutoNotify wraps the bugsnag.AutoNotify call
// Usage:
//
//	 go func() {
//	     ctx := bugsnag.StartSession(context.Background())
//			defer bugsnag.AutoNotify(ctx)
//	     // (possibly crashy code)
//	 }()
func AutoNotify(rawData ...interface{}) {
	bugsnag.AutoNotify(rawData)
}

func AttachRequestData(ctx context.Context, r *http.Request) context.Context {
	return bugsnag.AttachRequestData(ctx, r)
}
