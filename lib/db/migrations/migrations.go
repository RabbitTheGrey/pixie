package migrations

import (
	"bufio"
	"database/sql"
	"fmt"
	"pixie/lib/db"
	"pixie/lib/db/datamapper"
	"pixie/lib/db/migrations/query"
	"io"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
)

const templatePath = "/lib/db/migrations/template/template.tpl"
const registerMigrationsPath = "/migrations/migrations.go"

const (
	Schema = "public"
	Table  = "db_migration_version"
)

const (
	insertMigrationSql = "INSERT INTO %s (version, executed_at, execution_time) VALUES ($1, $2, $3)"
	deleteMigrationSql = "DELETE FROM %s WHERE version = $1"
)

type Function func(transaction *sql.Tx) error

type MigrationsResolver struct {
	database   db.IDatabase
	migrations []Migration
}

type Migration struct {
	Version       string    `column:"version"`        // Версия миграции. Например Version20251206202851
	ExecutedAt    time.Time `column:"executed_at"`    // Временная отметка применения миграции
	ExecutionTime float64   `column:"execution_time"` // Время выполнения SQL миграции
	Migrate       Function  // Функция применения миграции
	Rollback      Function  // Функция отмены миграции
}

func New() IMigrationsResolver {
	database, err := db.GetInstance(nil)

	if err != nil {
		panic(err)
	}

	return &MigrationsResolver{
		database:   database,
		migrations: make([]Migration, 0),
	}
}

func Destroy(mr *MigrationsResolver) {
	if mr != nil {
		mr = nil
	}
}

func (mr *MigrationsResolver) AppendMigration(version string, up Function, down Function) {
	mr.migrations = append(mr.migrations, Migration{
		Version:  version,
		Migrate:  up,
		Rollback: down,
	})
}

func (mr *MigrationsResolver) Migrate() error {
	newMigrations := mr.lookup()
	conn := mr.database.GetConnection()

	for _, migration := range newMigrations {
		startTime := time.Now()

		transaction, err := conn.Begin()
		if err != nil {
			fmt.Println("[ERROR] failed to start transaction:")
			return err
		}

		func() {
			defer func() {
				if r := recover(); r != nil {
					if rollbackErr := transaction.Rollback(); rollbackErr != nil {
						log.Printf("[PANIC] failed to rollback after panic: %v", rollbackErr)
					}
					log.Printf("[PANIC] recovered: %v", r)
				}
			}()

			// Выполнение миграции
			if err := migration.Migrate(transaction); err != nil {
				transaction.Rollback()
				fmt.Printf("[ERROR] failed to apply migration version %s: %v\n", migration.Version, err)
				panic(err)
			}

			duration := time.Since(startTime).Seconds()

			migration.ExecutedAt = time.Now().UTC()
			migration.ExecutionTime = duration

			table := mr.schemaPrefix() + "migration_versions"
			executionTimeFormatted := strconv.FormatFloat(migration.ExecutionTime, 'f', 2, 64)
			executedAtFormatted := migration.ExecutedAt.Format("2006-01-02 15:04:05.000")

			// Добавление миграции в таблицу
			_, err := transaction.Exec(insertMigrationSql, table, migration.Version, executedAtFormatted, executionTimeFormatted)
			if err != nil {
				transaction.Rollback()
				fmt.Printf("[ERROR] failed to insert migration record for version %s: %v\n", migration.Version, err)
				panic(err)
			}

			// Завершение транзакции
			if err := transaction.Commit(); err != nil {
				fmt.Println("[ERROR] failed to commit transaction:")
				panic(err)
			}

			fmt.Printf("[SUCCESS] applied migration version %s.\n", migration.Version)
		}()
	}

	return nil
}

func (mr *MigrationsResolver) Rollback() error {
	lastVersion := mr.getLastVersion()
	conn := mr.database.GetConnection()

	if lastVersion == "" {
		fmt.Print("Nothing to rollback\n")
		return nil
	}

	for _, migration := range mr.migrations {
		if migration.Version == lastVersion {
			transaction, err := conn.Begin()
			if err != nil {
				fmt.Printf("[ERROR] failed to start transaction for migration version %s: %v\n", migration.Version, err)
				return err
			}

			// Обработка паники и откат транзакции
			func() {
				defer func() {
					if r := recover(); r != nil {
						if rollbackErr := transaction.Rollback(); rollbackErr != nil {
							log.Printf("[PANIC] failed to rollback after panic: %v", rollbackErr)
						}
						log.Printf("[PANIC] recovered: %v", r)
					}
				}()

				// Выполнение rollback миграции
				if err := migration.Rollback(transaction); err != nil {
					transaction.Rollback()
					fmt.Printf("[ERROR] migration version %s rollback failed: %v\n", migration.Version, err)
					panic(err)
				}

				// Удаление миграции из таблицы
				table := mr.schemaPrefix() + "migration_versions"
				_, err := transaction.Exec(deleteMigrationSql, table, migration.Version)
				if err != nil {
					transaction.Rollback()
					fmt.Printf("[ERROR] failed to delete migration version %s: %v\n", migration.Version, err)
					panic(err)
				}

				// Завершение транзакции
				if err := transaction.Commit(); err != nil {
					fmt.Printf("[ERROR] failed to commit transaction for migration version %s: %v\n", migration.Version, err)
					panic(err)
				}

				fmt.Printf("[SUCCESS] rollback migration version %s.\n", migration.Version)
			}()

			return nil
		}
	}

	fmt.Print("... no migrations to rollback\n")
	return nil
}

func (mr *MigrationsResolver) Generate() error {
	root, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	templateFile := root + templatePath
	template, err := os.Open(templateFile)
	if err != nil {
		fmt.Println("Ошибка при открытии шаблона.")
		return err
	}
	defer template.Close()

	filepath := root + "/migrations/"
	version := "Version" + time.Now().Format("20060102150405")
	fullFilename := filepath + version + ".go"

	file, err := os.Create(fullFilename)
	if err != nil {
		fmt.Println("Ошибка при создании файла.")
		return err
	}
	defer file.Close()

	reader := bufio.NewReader(template)
	writer := bufio.NewWriter(file)

	for {
		line, readErr := reader.ReadString('\n')

		if readErr != nil && readErr != io.EOF {
			fmt.Println("Ошибка чтения шаблона.")
			return readErr
		}

		modifiedLine := strings.ReplaceAll(line, "{version}", version)
		_, writeErr := writer.WriteString(modifiedLine)
		if writeErr != nil {
			fmt.Println("Ошибка записи в файл.")
			return writeErr
		}

		if readErr == io.EOF {
			break
		}
	}

	if err := writer.Flush(); err != nil {
		fmt.Println("Ошибка при сохранении файла.")
		return err
	}

	mr.registerMigration(version)

	fmt.Println("[SUCCESS] generated new migration", fullFilename)
	return nil
}

// Поиск новых миграций не записанных в бд
func (mr *MigrationsResolver) lookup() []Migration {
	var newMigrations []Migration

	query.EnsureSchema(mr.database)
	query.CreateMigrationsTable(mr.database)

	lastVersion := mr.getLastVersion()

	for _, migration := range mr.migrations {
		if migration.Version > lastVersion {
			newMigrations = append(newMigrations, migration)
		}
	}

	return newMigrations
}

// Получение версии последней примененной миграции
func (mr *MigrationsResolver) getLastVersion() string {
	var lastVersion string

	row := query.GetLastVersion(mr.database)
	err := datamapper.SingleScalarResult(row, lastVersion)

	if err != nil {
		panic(err)
	}

	return lastVersion
}

// Префикс таблицы с использованием схемы
func (mr *MigrationsResolver) schemaPrefix() string {
	var schema string

	if mr.database.GetDriver() == db.DriverPostgresql {
		schema = "public."
	}

	return schema + "."
}

func (mr *MigrationsResolver) registerMigration(version string) error {
	root, err := os.Getwd()
	if err != nil {
		return err
	}
	filename := root + registerMigrationsPath

	codeBytes, err := os.ReadFile(filename)
	if err != nil {
		return err
	}

	code := string(codeBytes)

	newLine := "\tmr.AppendMigration(\"" + version + "\", " + version + "_Up, " + version + "_Down)\n"
	code = strings.Replace(code, "}", newLine+"}", 1)

	err = os.WriteFile(filename, []byte(code), 0644)
	if err != nil {
		return err
	}

	return nil
}
