package migrate

import (
	"context"
	"errors"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"sort"

	"github.com/jackc/pgx"
	"github.com/mna/redisc"
	"github.com/sirupsen/logrus"
)

type migrationType string

const (
	migrationTableName = "gomigrate"
	upMigration        = migrationType("up")
	downMigration      = migrationType("down")
	migrationLock      = ":migration-lock"
	rollbackLock       = ":rollback-lock"
	expirationTime     = 86400 //24 hrs
)

var (
	InvalidMigrationFile  = errors.New("Invalid migration file")
	InvalidMigrationPair  = errors.New("Invalid pair of migration files")
	InvalidMigrationsPath = errors.New("Invalid migrations path")
	InvalidMigrationType  = errors.New("Invalid migration type")
	NoActiveMigrations    = errors.New("No active migrations to rollback")
)

type Migrator struct {
	DB             *pgx.ConnPool
	MigrationsPath string
	dbAdapter      Postgres
	migrations     map[uint64]*Migration
	logger         Logger
	Redis          *redisc.Cluster
}

type Logger interface {
	Print(v ...interface{})
	Printf(format string, v ...interface{})
	Println(v ...interface{})
	Fatalf(format string, v ...interface{})
}

func (m *Migrator) Lock(key string, log logrus.FieldLogger) (bool, error) {
	conn := m.Redis.Get()
	defer conn.Close()

	reply, err := conn.Do("EXISTS", key)
	if err != nil {
		return false, err
	}
	switch reply.(int64) {
	case 0:
		// Keep an expiry key updated while the migration is running, this will by automatically culled in the 60s following the completion of the migration
		// this is intended to allow for migration retries and to prevent us from logging into production to resolve migration lock
		_, err := conn.Do("SETEX", key, expirationTime, true)
		return true, err
	default:
		return false, nil
	}
}

func (m *Migrator) Unlock(key string) error {
	conn := m.Redis.Get()
	defer conn.Close()
	_, err := conn.Do("DEL", key)
	return err
}

// Returns true if the migration table already exists.
func (m *Migrator) MigrationTableExists() (bool, error) {
	row := m.DB.QueryRow(m.dbAdapter.SelectMigrationTableSql(), migrationTableName)
	var tableName string
	err := row.Scan(&tableName)
	if err == pgx.ErrNoRows {
		m.logger.Print("Migrations table not found")
		return false, nil
	}
	if err != nil {
		m.logger.Printf("Error checking for migration table: %v", err)
		return false, err
	}
	m.logger.Print("Migrations table found")
	return true, nil
}

// Creates the migrations table if it doesn't exist.
func (m *Migrator) CreateMigrationsTable() error {
	_, err := m.DB.Exec(m.dbAdapter.CreateSchema())
	if err != nil {
		m.logger.Fatalf("Schema name was not provided for migrations: %v", err)
	}
	_, err = m.DB.Exec(m.dbAdapter.CreateMigrationTableSql())
	if err != nil {
		m.logger.Fatalf("Error creating migrations table: %v", err)
	}

	m.logger.Printf("Created migrations table: %s", migrationTableName)

	return nil
}

// Returns a new migrator.
func NewMigrator(db *pgx.ConnPool, adapter Postgres, migrationsPath string, redis *redisc.Cluster) (*Migrator, error) {
	return NewMigratorWithLogger(db, adapter, migrationsPath, redis, log.New(os.Stderr, "[gomigrate] ", log.LstdFlags))
}

// Returns a new migrator with the specified logger.
func NewMigratorWithLogger(db *pgx.ConnPool, adapter Postgres, migrationsPath string, redis *redisc.Cluster, logger Logger) (*Migrator, error) {
	// Normalize the migrations path.
	path := []byte(migrationsPath)
	pathLength := len(path)
	if path[pathLength-1] != '/' {
		path = append(path, '/')
	}

	logger.Printf("Migrations path: %s", path)

	migrator := Migrator{
		DB:             db,
		MigrationsPath: string(path),
		dbAdapter:      adapter,
		migrations:     make(map[uint64]*Migration),
		Redis:          redis,
		logger:         logger,
	}

	// Create the migrations table if it doesn't exist.
	tableExists, err := migrator.MigrationTableExists()
	if err != nil {
		return nil, err
	}
	if !tableExists {
		if err := migrator.CreateMigrationsTable(); err != nil {
			return nil, err
		}
	}

	// Get all metadata from the database.
	if err := migrator.fetchMigrations(); err != nil {
		return nil, err
	}
	if err := migrator.getMigrationStatuses(); err != nil {
		return nil, err
	}

	return &migrator, nil
}

// Populates a migrator with a sorted list of migrations from the file system.
func (m *Migrator) fetchMigrations() error {
	pathGlob := append([]byte(m.MigrationsPath), []byte("*")...)

	matches, err := filepath.Glob(string(pathGlob))
	if err != nil {
		m.logger.Fatalf("Error while globbing migrations: %v", err)
	}

	for _, match := range matches {
		num, migrationType, name, err := parseMigrationPath(match)
		if err != nil {
			m.logger.Printf("Invalid migration file found: %s", match)
			continue
		}

		m.logger.Printf("Migration file found: %s", match)

		migration, ok := m.migrations[num]
		if !ok {
			migration = &Migration{Id: num, Name: name, Status: Inactive}
			m.migrations[num] = migration
		}
		if migrationType == upMigration {
			migration.UpPath = match
		} else {
			migration.DownPath = match
		}
	}

	// Validate each migration.
	for _, migration := range m.migrations {
		if !migration.valid() {
			path := migration.UpPath
			if path == "" {
				path = migration.DownPath
			}
			m.logger.Printf("Invalid migration pair for path: %s", path)
			return InvalidMigrationPair
		}
	}

	m.logger.Printf("Migrations file pairs found: %v", len(m.migrations))

	return nil
}

// Queries the migration table to determine the status of each
// migration.
func (m *Migrator) getMigrationStatuses() error {
	for _, migration := range m.migrations {
		row := m.DB.QueryRow(m.dbAdapter.GetMigrationSql(), migration.Id)
		var mid uint64
		err := row.Scan(&mid)
		if err == pgx.ErrNoRows {
			continue
		}
		if err != nil {
			m.logger.Printf(
				"Error getting migration status for %s: %v",
				migration.Name,
				err,
			)
			return err
		}
		migration.Status = Active
	}
	return nil
}

// Returns a sorted list of migration ids for a given status. -1 returns
// all migrations.
func (m *Migrator) Migrations(status int) []*Migration {
	// Sort all migration ids.
	ids := make([]uint64, 0)
	for id, _ := range m.migrations {
		ids = append(ids, id)
	}
	sort.Sort(uint64slice(ids))

	// Find ids for the given status.
	migrations := make([]*Migration, 0)
	for _, id := range ids {
		migration := m.migrations[id]
		if status == -1 || migration.Status == status {
			migrations = append(migrations, migration)
		}
	}
	return migrations
}

// Applies a single migration.
func (m *Migrator) ApplyMigration(migration *Migration, mType migrationType) error {
	var path string
	if mType == upMigration {
		path = migration.UpPath
	} else if mType == downMigration {
		path = migration.DownPath
	} else {
		return InvalidMigrationType
	}

	m.logger.Printf("Applying migration: %s", path)

	sql, err := ioutil.ReadFile(path)
	if err != nil {
		m.logger.Printf("Error reading migration: %s", path)
		return err
	}
	transaction, err := m.DB.Begin()
	if err != nil {
		m.logger.Printf("Error opening transaction: %v", err)
		return err
	}

	// Certain adapters can not handle multiple sql commands in one file so we need the adapter to split up the command
	commands := m.dbAdapter.GetMigrationCommands(string(sql))

	// Perform the migration.
	for _, cmd := range commands {
		_, err := transaction.Exec(cmd)
		if err != nil {
			m.logger.Printf("Error executing migration: %v", err)
			if rollbackErr := transaction.Rollback(); rollbackErr != nil {
				m.logger.Printf("Error rolling back transaction: %v", rollbackErr)
				return rollbackErr
			}
			return err
		}
		// if &result != nil {
		// 	if rowsAffected := result.RowsAffected(); rowsAffected > 1 {
		// 		m.logger.Printf("Error getting rows affected: %v", rowsAffected)
		// 		if rollbackErr := transaction.Rollback(); rollbackErr != nil {
		// 			m.logger.Printf("Error rolling back transaction: %v", rollbackErr)
		// 			return rollbackErr
		// 		}
		// 		return err
		// 	}
		// }
	}

	// Log the event.
	if mType == upMigration {
		_, err = transaction.Exec(
			m.dbAdapter.MigrationLogInsertSql(),
			migration.Id,
		)
	} else {
		_, err = transaction.Exec(
			m.dbAdapter.MigrationLogDeleteSql(),
			migration.Id,
		)
	}
	if err != nil {
		m.logger.Printf("Error logging migration: %v", err)
		if rollbackErr := transaction.Rollback(); rollbackErr != nil {
			m.logger.Printf("Error rolling back transaction: %v", rollbackErr)
			return rollbackErr
		}
		return err
	}

	// Commit and update the struct status.
	if err := transaction.Commit(); err != nil {
		m.logger.Printf("Error commiting transaction: %v", err)
		return err
	}
	if mType == upMigration {
		migration.Status = Active
	} else {
		migration.Status = Inactive
	}

	return nil
}

// Applies all inactive migrations.
func (m *Migrator) Migrate(log logrus.FieldLogger) error {

	if m.Redis != nil {
		free, err := m.Lock(m.dbAdapter.SchemaName+migrationLock, log)
		if err != nil {
			return err
		}
		if !free {
			log.Info("locked")
			return nil
		}
		defer func() {
			err := m.Unlock(m.dbAdapter.SchemaName + migrationLock)
			if err != nil {
				log.WithError(err).Error()
			}
		}()
	}

	for _, migration := range m.Migrations(Inactive) {
		if err := m.ApplyMigration(migration, upMigration); err != nil {
			return err
		}
	}
	return nil
}

// Rolls back the last migration.
func (m *Migrator) Rollback(ctx context.Context, log logrus.FieldLogger) error {
	return m.RollbackN(ctx, 1, log)
}

// Rolls back N migrations.
func (m *Migrator) RollbackN(ctx context.Context, n int, log logrus.FieldLogger) error {

	if m.Redis != nil {
		free, err := m.Lock(m.dbAdapter.SchemaName+rollbackLock, log)
		if err != nil {
			return err
		}
		if !free {
			return nil
		}
		defer func() {
			err := m.Unlock(m.dbAdapter.SchemaName + rollbackLock)
			if err != nil {
				log.WithError(err).Error()
			}
		}()
	}

	migrations := m.Migrations(Active)
	if len(migrations) == 0 {
		return nil
	}

	if n > len(migrations) {
		n = len(migrations)
	}

	last_migration := len(migrations) - 1 - n

	for i := len(migrations) - 1; i != last_migration; i-- {
		if err := m.ApplyMigration(migrations[i], downMigration); err != nil {
			return err
		}
	}

	return nil
}

// Rolls back all migrations.
func (m *Migrator) RollbackAll(ctx context.Context, log logrus.FieldLogger) error {
	migrations := m.Migrations(Active)
	return m.RollbackN(ctx, len(migrations), log)
}
