package db

import (
	"database/sql"
	"pixie/lib/db/driver_enum"
	"pixie/lib/db/dsn"
	"sync"
)

const (
	DriverMysql      string = "mysql"
	DriverPostgresql string = "postgres"
	DriverSqlite     string = "sqlite"
)

var instance *Database
var once sync.Once

type Database struct {
	driver     string
	connection *sql.DB
}

type DBConfig struct {
	Driver             string // Драйвер БД
	MysqlUser          string // Пользователь MySQL
	MysqlPassword      string // Пароль от пользователя MySQL
	MysqlHost          string // Адрес/IP MySQL сервера
	MysqlPort          string // Порт MySQL сервера
	MysqlDBName        string // Имя MySQL базы данных
	PostgresqlUser     string // Пользователь PostgreSQL
	PostgresqlPassword string // Пароль от пользователя PostgreSQL
	PostgresqlHost     string // Адрес/IP PostgreSQL сервера
	PostgresqlPort     string // Порт PostgreSQL сервера
	PostgresqlDBName   string // Имя PostgreSQL базы данных
	PostgresqlSslMode  string // Использовать SSL подключение к PostgreSQL
	SqlitePath         string // Путь к локальной SQLite базе данных
}

// Инициализация подключения к БД
func GetInstance(config *DBConfig) (IDatabase, error) {
	var err error = nil

	once.Do(func() {
		connection, err := newConnection(config)
		instance, _ = &Database{
			connection: connection,
			driver:     config.Driver,
		}, err
	})

	return instance, err
}

func (db *Database) GetConnection() *sql.DB {
	if db.connection != nil {
		return db.connection
	}

	panic("Невозможно получить доступ к экземпляру БД до его инициализации.\n")
}

func (db *Database) GetDriver() string {
	return db.driver
}

func (db *Database) Close() {
	if db.connection != nil {
		db.connection.Close()
	}
}

// Фабричный метод для db.Connection
func newConnection(config *DBConfig) (*sql.DB, error) {
	var dsn string

	switch config.Driver {
	case driver_enum.Mysql:
		dsn = createMysqlDsn(config)
	case driver_enum.Postgresql:
		dsn = createPostgresqlDsn(config)
	case driver_enum.Sqlite:
		dsn = createSqliteDsn(config)
	}

	return sql.Open(config.Driver, dsn)
}

// DSN при выбранном драйвере MySQL
func createMysqlDsn(config *DBConfig) string {
	mysql := dsn.MysqlDsn{
		User:     config.MysqlUser,
		Password: config.MysqlPassword,
		Host:     config.MysqlHost,
		Port:     config.MysqlPort,
		Db:       config.MysqlDBName,
	}
	return mysql.GetConnectionString()
}

// DSN при выбранном драйвере PostgreSQL
func createPostgresqlDsn(config *DBConfig) string {
	postresql := dsn.PostgresqlDsn{
		User:     config.PostgresqlUser,
		Password: config.PostgresqlPassword,
		Host:     config.PostgresqlHost,
		Port:     config.PostgresqlPort,
		Db:       config.PostgresqlDBName,
		SslMode:  config.PostgresqlSslMode,
	}
	return postresql.GetConnectionString()
}

// DSN при выбранном драйвере SQLite
func createSqliteDsn(config *DBConfig) string {
	sqlite := dsn.SqliteDsn{
		Path: config.SqlitePath,
	}
	return sqlite.GetConnectionString()
}
