package query

import (
	"pixie/lib/db"
)

func EnsureSchema(database db.IDatabase) error {
	if database.GetDriver() != db.DriverPostgresql {
		return nil
	}

	sql := `
		CREATE IF NOT EXISTS SCHEMA public
	`
	_, err := database.GetConnection().Exec(sql)
	return err
}
