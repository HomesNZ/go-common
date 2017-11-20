package newrelic

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"sync"

	"github.com/HomesNZ/go-common/env"
	"github.com/Sirupsen/logrus"
	"github.com/gorilla/mux"
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

// NewContext returns a new context with txn added as a value
func NewContext(ctx context.Context, txn newrelic.Transaction) context.Context {
	return context.WithValue(ctx, transactionKey, txn)
}

// FromContext returns the stored Transaction if it exists. Returned bool will be false if not found
func FromContext(ctx context.Context) (newrelic.Transaction, bool) {
	txn, ok := ctx.Value(transactionKey).(newrelic.Transaction)
	return txn, ok
}

// StartTransaction begins a Transaction.
// * The Transaction should only be used in a single goroutine.
// * This method never returns nil.
// * If an http.Request is provided then the Transaction is considered
//   a web transaction.
// * If an http.ResponseWriter is provided then the Transaction can be
//   used in its place.  This allows instrumentation of the response
//   code and response headers.
func StartTransaction(name string, w http.ResponseWriter, r *http.Request) newrelic.Transaction {
	return app.StartTransaction(name, w, r)
}

// Middleware is an easy way to implement NewRelic as middleware in an Alice
// chain.
func Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if app != nil {
			name := routeName(r)
			txn := StartTransaction(name, w, r)
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

func routeName(r *http.Request) string {
	route := mux.CurrentRoute(r)
	if nil == route {
		return r.URL.Path
	}
	if n := route.GetName(); n != "" {
		return n
	}
	if n, _ := route.GetPathTemplate(); n != "" {
		return n
	}
	n, _ := route.GetHostTemplate()
	return n
}
