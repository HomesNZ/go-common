package db

import (
	"database/sql"
	"fmt"
	"strings"
	"sync"

	"github.com/HomesNZ/go-common/env"

	// SQL driver
	_ "github.com/lib/pq"

	log "github.com/Sirupsen/logrus"
	"github.com/cenkalti/backoff"
)

var (

	// connRedhshift is the current redshift connection
	connRedhshift *sql.DB

	// onceRedshift prevents InitConnection from being called more than once in Conn
	onceRedshift = sync.Once{}
)

// InitConnectionRedhshift creates a new new connection to the database and verifies that it succeeds.
func InitConnectionRedhshift() {
	db := RS{}
	db.Open()
	connRedhshift = db.Conn

}

// SetConnection manually sets the connection.
func SetConnectionRedhsift(db *sql.DB) {
	// This stops InitConnection from being called again and clobbering the connection..
	onceRedshift.Do(func() {})

	connRedhshift = db

}

// RS is a concrete implementation of a Redshift connection
type RS PG

// Conn is the SQL database connection accessor. If the connection is nil, it will be initialized.
func ConnRedshift() *sql.DB {
	if connRedhshift == nil {
		onceRedshift.Do(InitConnection)
	}
	return conn
}

// Open will initialize the database connection or raise an error.
func (db *RS) Open() {
	c, err := sql.Open("postgres", db.connectionString())
	if err != nil {
		log.Error(err)
		log.Fatal(ErrUnableToParseDBConnection)
	}

	db.Conn = c

	err = db.verifyConnection()
	if err != nil {
		log.Error(err)
		log.Fatal(ErrUnableToConnectToDB)
	}
}

// verifyConnection pings the database to verify a connection is established. If the connection cannot be established,
// it will retry with an exponential back off.
func (db RS) verifyConnection() error {
	log.Infof("Attempting to connect to database: %s", db.logSafeConnectionString())

	pingDB := func() error {
		return db.Conn.Ping()
	}

	expBackoff := backoff.NewExponentialBackOff()
	expBackoff.MaxElapsedTime = ConnBackoffTimeout

	err := backoff.Retry(pingDB, expBackoff)
	if err != nil {
		log.Warning(err)
		return ErrUnableToConnectToDB
	}

	log.Info("Connected to database")

	return nil
}

// connectionString returns the database connection string.
func (db RS) connectionString() string {
	password := env.GetString("DB_PASSWORD", "")
	if password != "" {
		password = ":" + password
	}

	connString := fmt.Sprintf(
		"postgres://%s%s@%s:%s/%s?sslmode=%s",
		env.GetString("REDSHIFT_USER", "postgres"),
		password,
		env.GetString("REDSHIFT_HOST", "localhost"),
		env.GetString("REDSHIFT_PORT", "5439"),
		env.MustGetString("REDSHIFT_NAME"),
		env.GetString("REDSHIFT_SSL_MODE", "disable"),
	)

	searchPath := env.GetString("REDSHIFT_SEARCH_PATH", "")
	if len(searchPath) > 0 {
		connString = fmt.Sprintf("%s&search_path=%s", connString, searchPath)
	}
	return connString
}

// logSafeConnectionString is the database connection string with the password replace with `****` so it can be logged
// without revealing the password.
func (db RS) logSafeConnectionString() string {
	c := db.connectionString()

	password := env.GetString("REDSHIFT_PASSWORD", "")
	if password != "" {
		c = strings.Replace(c, password, "****", 1)
	}

	return c
}
