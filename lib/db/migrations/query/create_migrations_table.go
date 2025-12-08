package query

import "pixie/lib/db"

func CreateMigrationsTable(database db.IDatabase) error {
	var sql string

	switch database.GetDriver() {
	case db.DriverPostgresql:
		sql = `
			CREATE IF NOT EXISTS TABLE public.db_migration_versions (
				COLUMN version VARCHAR(255) NOT NULL,
				COLUMN executed_at TIMESTAMP(0) WITHOUT TIME ZONE NOT NULL DEFAULT NOW(),
				COLUMN execution_time NUMERIC(*, 2) NOT NULL,
				PRIMARY KEY(version)
			)
		`
	case db.DriverMysql:
		sql = `
			CREATE IF NOT EXISTS TABLE db_migration_versions (
				COLUMN version VARCHAR(255) NOT NULL,
				COLUMN executed_at TIMESTAMP(0) WITHOUT TIME ZONE NOT NULL DEFAULT NOW(),
				COLUMN execution_time DECIMAL(5, 2) NOT NULL,
				PRIMARY KEY(version)
			)
		`
	case db.DriverSqlite:
		sql = `
			CREATE TABLE IF NOT EXISTS db_migration_versions (
				version TEXT NOT NULL,
				executed_at TEXT NOT NULL DEFAULT (datetime('now')),
				execution_time REAL NOT NULL,
				PRIMARY KEY (version)
			);
		`
	}

	_, err := database.GetConnection().Exec(sql)
	return err
}
