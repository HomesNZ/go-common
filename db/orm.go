package db

import (
	"database/sql"
	"sync"

	"github.com/HomesNZ/go-common/env"
	"github.com/jinzhu/gorm"
)

var (
	// UseORM is a flag that controls if the ORM connection should be updated when the normal connection is updated.
	UseORM = false

	// orm is the Gorm wrapped current database connection
	orm *gorm.DB

	// ormOnce guards InitORM from running more than once
	ormOnce = sync.Once{}
)

// InitORM initializes the ORM connection from the existing connection.
func InitORM(service string) {
	g, err := gorm.Open("postgres", Conn(service))
	if err != nil {
		// This shouldnt happen unless our DB settings are malformed?
		panic(err)
	}

	if env.GetBool("LOG_ORM_QUERIES", false) {
		g.LogMode(true)
	}

	orm = g
}

// ORM is the gorm wrapped SQL database connection. If the connection is nil, it will be initialized.
func ORM(service string) *gorm.DB {
	UseORM = true

	once.Do(func() { InitConnection(service) })
	ormOnce.Do(func() { InitORM(service) })
	return orm
}

// SetORMConnection manually sets the ORM connection.
func SetORMConnection(db *sql.DB) {
	g, err := gorm.Open("postgres", conn)
	if err != nil {
		// This shouldnt happen unless our DB settings are malformed?
		panic(err)
	}

	orm = g
}
