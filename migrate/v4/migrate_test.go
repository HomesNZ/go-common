package migrate

import (
	"context"
	"fmt"
	"testing"

	"github.com/HomesNZ/go-common/dbclient/v4"
	log "github.com/sirupsen/logrus"
)

func TestCreatingMigratorWhenTableExists(t *testing.T) {
	ctx := context.Background()

	cfg := dbclient.Config{
		MaxConns: 1,
		Host:     "localhost",
		Name:     "postgres",
		User:     "postgres",
		Port:     5432,
	}
	adapter := Postgres{SchemaName: "test"}
	path := fmt.Sprintf("test_migrations/")
	db, _ := dbclient.Conn(ctx, &cfg)
	m, err := NewMigrator(ctx, db, adapter, path, nil)
	if err != nil {
		t.Error(err)
	}
	err = m.Migrate(ctx, log.WithField("startup", "migrate"))
	row := db.QueryRow(ctx, "select migration_id from test.gomigrate")
	var id uint64
	err = row.Scan(&id)
	if err != nil {
		t.Error(err)
	}
	if id != 1 {
		t.Error("Invalid id found in database")
	}
	//cleanup()
}

//func TestMigrationAndRollback(t *testing.T) {
//	m := GetMigrator("test1")
//
//	if err := m.Migrate(); err != nil {
//		t.Error(err)
//	}
//
//	// Ensure that the migration ran.
//	row := db.QueryRow(
//		adapter.SelectMigrationTableSql(),
//		"test",
//	)
//	var tableName string
//	if err := row.Scan(&tableName); err != nil {
//		t.Error(err)
//	}
//	if tableName != "test" {
//		t.Errorf("Migration table not created")
//	}
//	// Ensure that the migrate status is correct.
//	row = db.QueryRow(
//		adapter.GetMigrationSql(),
//		1,
//	)
//	var status int
//	if err := row.Scan(&status); err != nil {
//		t.Error(err)
//	}
//	if status != Active || m.migrations[1].Status != Active {
//		t.Error("Invalid status for migration")
//	}
//	if err := m.RollbackN(len(m.migrations)); err != nil {
//		t.Error(err)
//	}
//
//	// Ensure that the down migration ran.
//	row = db.QueryRow(
//		adapter.SelectMigrationTableSql(),
//		"test",
//	)
//	err := row.Scan(&tableName)
//	if err != nil && err != sql.ErrNoRows {
//		t.Errorf("Migration table should be deleted: %v", err)
//	}
//
//	// Ensure that the migration log is missing.
//	row = db.QueryRow(
//		adapter.GetMigrationSql(),
//		1,
//	)
//	if err := row.Scan(&status); err != nil && err != sql.ErrNoRows {
//		t.Error(err)
//	}
//	if m.migrations[1].Status != Inactive {
//		t.Errorf("Invalid status for migration, expected: %d, got: %v", Inactive, m.migrations[1].Status)
//	}
//
//	cleanup()
//}

//func cleanup() {
//	ctx :=context.Background()
//	_, err := db.Exec(ctx, "drop table gomigrate")
//	if err != nil {
//		panic(err)
//	}
//}

//func init() {
//	var err error
//	dbType = "pg"
//	ctx := context.Background()
//	log.Print("Using postgres")
//	//adapter = Postgres{}
//	//db, err = sql.Open("postgres", "host=localhost dbname=gomigrate sslmode=disable")
//
//	db, err = dbclient.Conn(ctx, &dbclient.Config{
//		MaxConns: 1,
//		Host: "localhost",
//		Name: "postgres",
//		User: "postgres",
//	})
//
//	if err != nil {
//		panic(err)
//	}
//}
