package migrate

import (
	"context"
	"fmt"
	mockRedis "github.com/HomesNZ/go-common/redis/mock"
	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo"
	"github.com/pkg/errors"

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
	redisCtrl *gomock.Controller
	redisMock     *mockRedis.MockCache
)

func TestNewMigrator(t *testing.T) {
	m := getMigrator()

	if len(m.Migrations(-1)) != 3 {
		t.Errorf("Invalid number of migrations detected")
	}
	migrations := m.Migrations(-1)
	migration := migrations[0]

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
	migrations := m.Migrations(-1)

	if status != Active || migrations[1].Status != Active {
		t.Error("Invalid status for migration")
	}
	if err := m.RollbackN(ctx, len(migrations)+1, logger); err != nil {
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
	if migrations[1].Status != Inactive {
		t.Errorf("Invalid status for migration, expected: %d, got: %v", Inactive, migrations[1].Status)
	}

	cleanup()
}

//should Lock redis
func TestLockSuccess(t *testing.T) {
	m := getMigratorWithRedis()
	redisMock.EXPECT().Exists(gomock.Any()).Return(true, nil)
	redisMock.EXPECT().SetExpiry(gomock.Any(),gomock.Any(),gomock.Any()).Return(nil)
	res, _ := m.Lock("test", logger)
	assert.Equal(t, res, true)
}

//should failure Locking, because key doen't exist
func TestLockFailure(t *testing.T) {
	m := getMigratorWithRedis()
	redisMock.EXPECT().Exists(gomock.Any()).Return(false, errors.New("Error"))
	res, _ := m.Lock("test", logger)

	assert.Equal(t, res, false)
}

func TestUnlock(t *testing.T) {
	m := getMigratorWithRedis()
	redisMock.EXPECT().Delete(gomock.Any()).Return("",nil)
	err := m.Unlock("test")
	assert.Equal(t, err, nil)
}

func getMigrator() Migrator {
	path := fmt.Sprintf("%s", "test_migrations")
	m, err := NewMigrator(ctx, db, adapter, path, nil)
	if err != nil {
		panic(err)
	}
	return m
}

func getMigratorWithRedis() Migrator {
	path := fmt.Sprintf("%s", "test_migrations")
	m, err := NewMigrator(ctx, db, adapter, path, redisMock)
	if err != nil {
		panic(err)
	}
	return m
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
	redisCtrl = gomock.NewController(GinkgoT())
	redisMock = mockRedis.NewMockCache(redisCtrl)

	if err != nil {
		panic(err)
	}
}
