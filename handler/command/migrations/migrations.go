package migrations

import (
	"fmt"
	"pixie/lib/console"
	"pixie/lib/db/migrations"
	migrations_map "pixie/migrations"
)

func Generate(args map[string]string) int {
	resolver := migrations.New()
	err := resolver.Generate()

	if err != nil {
		fmt.Println(err.Error())
		return console.Failure
	}

	return console.Success
}

func Migrate(args map[string]string) int {
	resolver := migrations.New()
	migrations_map.RegisterMigrations(resolver)

	return console.Success
}

func Rollback(args map[string]string) int {
	return console.Success
}
