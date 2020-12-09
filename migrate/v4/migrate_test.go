package migrate

import (
	"context"
	"fmt"
	"github.com/gomodule/redigo/redis"
	"github.com/mna/redisc"
	"time"

	pgx "github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	//"github.com/rafaeljusto/redigomock"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"testing"
)

var (
	logger  logrus.FieldLogger
	adapter Postgres
	db      *pgxpool.Pool
	ctx     context.Context
)

func TestNewMigrator(t *testing.T) {
	m := getMigrator()

	if len(m.migrations) != 3 {
		t.Errorf("Invalid number of migrations detected")
	}
	migration := m.migrations[1]

	if migration.Name != "test" {
		t.Errorf("Invalid migration name detected: %s", migration.Name)
	}
	if migration.Id != 1 {
		t.Errorf("Invalid migration num detected: %d", migration.Id)
	}
	if migration.Status != Inactive {
		t.Errorf("Invalid migration num detected: %d", migration.Status)
	}

	cleanup()
}

func TestCreatingMigratorWhenTableExists(t *testing.T) {
	// Create the table and populate it with a row.
	_, err := db.Exec(ctx, adapter.CreateMigrationTableSql())
	if err != nil {
		t.Error(err)
	}
	_, err = db.Exec(ctx, adapter.MigrationLogInsertSql(), 123)
	if err != nil {
		t.Error(err)
	}

	getMigrator()

	// Check that our row is still present.
	row := db.QueryRow(ctx, "select migration_id from test.gomigrate")
	var id uint64
	err = row.Scan(&id)
	if err != nil {
		t.Error(err)
	}
	if id != 123 {
		t.Error("Invalid id found in database")
	}
	cleanup()
}

func TestMigrationAndRollback(t *testing.T) {
	cleanup()
	m := getMigrator()
	if err := m.Migrate(ctx, logger); err != nil {
		t.Error(err)
	}

	// Ensure that the migration ran.
	row := db.QueryRow(
		ctx,
		adapter.SelectMigrationTableSql(),
		"gomigrate",
	)
	var tableName string
	if err := row.Scan(&tableName); err != nil {
		t.Error(err)
	}
	if tableName != "gomigrate" {
		t.Errorf("Migration table not created")
	}
	// Ensure that the migrate status is correct.
	row = db.QueryRow(
		ctx,
		adapter.GetMigrationSql(),
		1,
	)
	var status int
	if err := row.Scan(&status); err != nil {
		t.Error(err)
	}
	if status != Active || m.migrations[1].Status != Active {
		t.Error("Invalid status for migration")
	}
	if err := m.RollbackN(ctx, len(m.migrations)+1, logger); err != nil {
		t.Error(err)
	}

	// Ensure that the down migration ran.
	row = db.QueryRow(
		ctx,
		adapter.SelectMigrationTableSql(),
		"gomigrate",
	)
	err := row.Scan(&tableName)
	if err != nil && err != pgx.ErrNoRows {
		t.Errorf("Migration table should be deleted: %v", err)
	}

	// Ensure that the migration log is missing.
	row = db.QueryRow(
		ctx,
		adapter.GetMigrationSql(),
		1,
	)
	if err := row.Scan(&status); err != nil && err != pgx.ErrNoRows {
		t.Error(err)
	}
	if m.migrations[1].Status != Inactive {
		t.Errorf("Invalid status for migration, expected: %d, got: %v", Inactive, m.migrations[1].Status)
	}

	cleanup()
}

//should Lock redis
func TestLockSuccess(t *testing.T) {
	m := getMigratorWithRedis()
	res, _ := m.Lock("test", logger)
	assert.Equal(t, res, true)
}

//should failure Locking, because key doen't exist
func TestLockFailure(t *testing.T) {
	m := getMigratorWithRedis()
	res, _ := m.Lock("test", logger)

	assert.Equal(t, res, false)
}

func TestUnlock(t *testing.T) {
	m := getMigratorWithRedis()
	err := m.Unlock("test")
	assert.Equal(t, err, nil)
}

func getMigrator() *Migrator {
	path := fmt.Sprintf("%s", "test_migrations")
	m, err := NewMigrator(ctx, db, adapter, path, nil)
	if err != nil {
		panic(err)
	}
	return m
}

func getMigratorWithRedis() *Migrator {

	redisCluster := &redisc.Cluster{
		StartupNodes: []string{":6379"},
		DialOptions:  []redis.DialOption{redis.DialConnectTimeout(2 * time.Second)},
		CreatePool:   createPool,
	}

	path := fmt.Sprintf("%s", "test_migrations")
	m, err := NewMigrator(ctx, db, adapter, path, redisCluster)
	if err != nil {
		panic(err)
	}
	return m
}

func createPool(addr string, opts ...redis.DialOption) (*redis.Pool, error) {
	//mock := redigomock.NewConn()
	//mock.Command(com, key).Expect(exp)
	//
	//p := &redis.Pool{
	//	Dial: func() (redis.Conn, error) { return mock, nil },
	//}

	return &redis.Pool{
		MaxIdle:     5,
		MaxActive:   10,
		IdleTimeout: time.Minute,
		//Dial: func() (redis.Conn, error) { return mock, nil },
		Dial: func() (redis.Conn, error) {
			return redis.Dial("tcp", addr, opts...)
		},
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			_, err := c.Do("PING")
			return err
		},
	}, nil
}

func cleanup() {
	_, err := db.Exec(ctx, "drop table test.gomigrate")
	_, err = db.Exec(ctx, "drop table if exists test.tt")
	if err != nil {
		panic(err)
	}
}

func init() {
	var err error
	logger = logrus.New()
	adapter = Postgres{SchemaName: "test"}
	ctx = context.Background()
	db, err = pgxpool.Connect(ctx, "host=localhost user=postgres dbname=postgres sslmode=disable")

	if err != nil {
		panic(err)
	}
}
