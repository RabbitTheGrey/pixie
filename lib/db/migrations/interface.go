package migrations

type IMigrationsResolver interface {
	// Зарегистрировать миграцию в резолвере
	AppendMigration(version string, up Function, down Function)
	// Применить новые миграции
	Migrate() error
	// Откат последней миграции
	Rollback() error
	// Сгенерировать новый файл с миграцией
	Generate() error
}
