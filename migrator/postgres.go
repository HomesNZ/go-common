package migrator

const migrationTableName = "gomigrate"

// POSTGRES

type Postgres struct {
	SchemaName string
}

func (p Postgres) tableName() string {
	return "\"" + p.SchemaName + "\".\"" + migrationTableName + "\""
}

func (p Postgres) SelectMigrationTableSql() string {
	return "SELECT tablename FROM pg_catalog.pg_tables WHERE tablename = $1 AND schemaname = '" + p.SchemaName + "'"
}

func (p Postgres) CreateMigrationTableSql() string {
	return `
	CREATE SCHEMA IF NOT EXISTS ` + p.SchemaName + `;
	
	CREATE TABLE ` + p.tableName() + ` (
		id				SERIAL PRIMARY KEY,
		migration_id	BIGINT UNIQUE NOT NULL
	);`
}

func (p Postgres) GetMigrationSql() string {
	return `SELECT migration_id FROM ` + p.tableName() + ` WHERE migration_id = $1`
}

func (p Postgres) MigrationLogInsertSql() string {
	return `INSERT INTO ` + p.tableName() + ` (migration_id) values ($1)`
}

func (p Postgres) MigrationLogDeleteSql() string {
	return `DELETE FROM ` + p.tableName() + ` WHERE migration_id = $1`
}

func (p Postgres) GetMigrationCommands(sql string) []string {
	return []string{sql}
}
