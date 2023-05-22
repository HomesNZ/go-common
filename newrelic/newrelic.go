package newrelic

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/pkg/errors"

	"github.com/HomesNZ/go-common/env"
	"github.com/gorilla/mux"
	newrelic "github.com/newrelic/go-agent/v3/newrelic"
)

var (
	App *newrelic.Application
)

type contextKey int

var transactionKey contextKey = 0

// NewFromEnv initializes the NewRelic configuration
func NewFromEnv(appName string) error {
	var err error
	apiKey := env.GetString("NEWRELIC_API_KEY", "")
	if apiKey == "" {
		return errors.New("NEWRELIC_API_KEY is required")
	}
	e := env.Env()
	if e == "" {
		e = "development"
	}

	App, err = newrelic.NewApplication(
		newrelic.ConfigAppName(fmt.Sprintf("%s-%s", appName, e)),
		newrelic.ConfigLicense(apiKey),
	)
	if err != nil {
		return errors.Wrap(err, "failed to initialize New Relic")
	}
	return err
}

// NewContext returns a new context with txn added as a value
func NewContext(ctx context.Context, txn newrelic.Transaction) context.Context {
	return context.WithValue(ctx, transactionKey, txn)
}

// FromContext returns the stored Transaction if it exists. Returned bool will be false if not found
func FromContext(ctx context.Context) (newrelic.Transaction, bool) {
	txn, ok := ctx.Value(transactionKey).(newrelic.Transaction)
	return txn, ok
}

// Middleware is an easy way to implement NewRelic as middleware in an Alice
// chain.
func Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if App != nil {
			name := routeName(r)
			txn := App.StartTransaction(name)
			defer txn.End()
			for k, v := range r.URL.Query() {
				txn.AddAttribute(k, strings.Join(v, ","))
			}
			w = txn.SetWebResponse(w)
			r = newrelic.RequestWithTransactionContext(r, txn)
		}
		next.ServeHTTP(w, r)
	})
}

func Shutdown(duration time.Duration) {
	if App != nil {
		App.Shutdown(duration)
	}
}

func routeName(r *http.Request) string {
	route := mux.CurrentRoute(r)
	if nil == route {
		return r.URL.Path
	}
	if n := route.GetName(); n != "" {
		return n
	}
	if n, _ := route.GetPathTemplate(); n != "" {
		return r.Method + " " + n
	}
	n, _ := route.GetHostTemplate()
	return r.Method + " " + n
}
