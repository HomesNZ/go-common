package newrelic

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"sync"

	"github.com/HomesNZ/go-common/env"
	"github.com/Sirupsen/logrus"
	newrelic "github.com/newrelic/go-agent"
)

var (
	app      newrelic.Application
	initOnce = sync.Once{}
)

type contextKey int

var transactionKey contextKey = 0

// InitNewRelic initializes the NewRelic configuration and panics if there is an
// error.
func InitNewRelic(appName string) {
	var err error
	apiKey := env.GetString("NEWRELIC_API_KEY", "")
	if apiKey == "" {
		logrus.Info("Skipping New Relic initialization - NEWRELIC_API_KEY is empty")
		return
	}
	e := env.Env()
	if e == "" {
		e = "development"
	}
	config := newrelic.NewConfig(fmt.Sprintf("%s-%s", appName, e), apiKey)
	app, err = newrelic.NewApplication(config)
	if err != nil {
		panic(err)
	}
	logrus.Info("New Relic initialized successfully")
}

func NewContext(ctx context.Context, txn newrelic.Transaction) context.Context {
	return context.WithValue(ctx, transactionKey, txn)
}

func FromContext(ctx context.Context) (newrelic.Transaction, bool) {
	txn, ok := ctx.Value(transactionKey).(newrelic.Transaction)
	return txn, ok
}

// Middleware is an easy way to implement NewRelic as middleware in an Alice
// chain.
func Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if app != nil {
			txn := app.StartTransaction(r.URL.Path, w, r)
			for k, v := range r.URL.Query() {
				txn.AddAttribute(k, strings.Join(v, ","))
			}
			defer txn.End()
			w = txn
			r = r.WithContext(NewContext(r.Context(), txn))
		}
		next.ServeHTTP(w, r)
	})
}
