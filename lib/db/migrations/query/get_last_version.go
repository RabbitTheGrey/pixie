package query

import (
	"database/sql"
	"pixie/lib/db"
)

func GetLastVersion(database db.IDatabase) *sql.Row {
	var sql string

	switch database.GetDriver() {
	case db.DriverPostgresql:
		sql = `
			SELECT dmv.version
			FROM public.db_migration_versions dmv
			ORDER BY dmv.version DESC
			LIMIT 1
		`
	case db.DriverMysql, db.DriverSqlite:
		sql = `
			SELECT dmv.version
			FROM db_migration_versions dmv
			ORDER BY dmv.version DESC
			LIMIT 1;
		`
	}

	return database.GetConnection().QueryRow(sql)
}
