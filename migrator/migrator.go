package migrator

import (
	"database/sql"
	"time"

	"github.com/gomodule/redigo/redis"
	"github.com/sirupsen/logrus"

	"github.com/DavidHuie/gomigrate"
)

const (
	migrationLock = ":migration-lock"
	rollbackLock  = ":rollback-lock"
)

// Migrator handles DB migration operations
type Migrator struct {
	ServiceName string
	DB          *sql.DB
	Path        string
	Postgres    Postgres
	Redis       *redis.Pool
}

func (m Migrator) Lock(key string, log logrus.FieldLogger) (bool, error) {
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
		go func() {
			t := time.NewTicker(time.Second * 45)
			for range t.C {
				_, err := conn.Do("SETEX", key, int(time.Second*60), true)
				if err != nil {
					log.WithError(err).Fatal()
				}
			}
		}()
		return true, err
	default:
		return false, nil
	}
}

func (m Migrator) Unlock(key string) error {
	conn := m.Redis.Get()
	defer conn.Close()
	_, err := conn.Do("DEL", key)
	return err
}

// Migrate runs any pending migrations
func (m Migrator) Migrate(log logrus.FieldLogger) error {

	if m.Redis != nil {
		free, err := m.Lock(m.ServiceName+migrationLock, log)
		if err != nil {
			return err
		}
		if !free {
			log.Info("locked")
			return nil
		}
		defer func() {
			err := m.Unlock(m.ServiceName + migrationLock)
			if err != nil {
				log.WithError(err).Error()
			}
		}()
	}

	gm, err := gomigrate.NewMigratorWithLogger(
		m.DB,
		m.Postgres,
		m.Path,
		log,
	)
	if err != nil {
		return err
	}

	return gm.Migrate()
}

// Rollback rolls back the last run migration
func (m Migrator) Rollback(log logrus.FieldLogger, steps int) error {

	if m.Redis != nil {
		free, err := m.Lock(m.ServiceName+rollbackLock, log)
		if err != nil {
			return err
		}
		if !free {
			return nil
		}
		defer func() {
			err := m.Unlock(m.ServiceName + rollbackLock)
			if err != nil {
				log.WithError(err).Error()
			}
		}()
	}

	gm, err := gomigrate.NewMigratorWithLogger(
		m.DB,
		m.Postgres,
		m.Path,
		log,
	)
	if err != nil {
		return err
	}
	if steps > len(gm.Migrations(gomigrate.Active)) {
		return gm.RollbackAll()
	}
	return gm.RollbackN(steps)
}
