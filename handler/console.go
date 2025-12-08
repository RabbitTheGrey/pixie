package handler

import (
	"pixie/handler/command"
	"pixie/handler/command/migrations"
	"pixie/lib/console"
)

// Обертка для объявления команд через c.AppendCommand()
func RegisterCommands(c console.IConsole) {
	// Example
	c.AppendCommand("say_hello", command.SayHelloCommand)

	// Migrations
	//c.AppendCommand("migration_migrate", migrations.Migrate)
	//c.AppendCommand("migration_rollback", migrations.Rollback)
	c.AppendCommand("migration_generate", migrations.Generate)

}
