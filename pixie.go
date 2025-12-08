package main

import (
	"flag"
	"fmt"
	"pixie/handler"
	"pixie/lib/console"
	"pixie/lib/db"
	"pixie/lib/db/driver_enum"
	"pixie/lib/dotenv"
	"pixie/lib/http/router"
	"pixie/lib/http/server"
	"strconv"

	_ "github.com/lib/pq"
)

// Выполнение консольной команды
func runCommand() {
	// Инициализация консольных команд
	console := console.New()
	handler.RegisterCommands(console)

	console.Execute()
}

// Запуск http сервера
func runServer(env dotenv.IDotenv) {
	// Инициализация роутера
	router := router.New()
	handler.RegisterRoutes(router)

	// router.withGlobalMiddlewares(
	// 	middleware,
	// 	middleware,
	// )

	serverHost := env.Getenv("SERVER_HOST")
	serverPort, _ := strconv.Atoi(env.Getenv("SERVER_PORT"))
	readTimeout, _ := strconv.Atoi(env.Getenv("SERVER_READ_TIMEOUT"))
	writeTimeout, _ := strconv.Atoi(env.Getenv("SERVER_WRITE_TIMEOUT"))
	idleTimeout, _ := strconv.Atoi(env.Getenv("SERVER_IDLE_TIMEOUT"))

	// Запуск сервера
	httpServer, err := server.GetInstance(
		server.ServerHost(serverHost),
		server.ServerPort(serverPort),
		server.ReadTimeout(readTimeout),
		server.WriteTimeout(writeTimeout),
		server.IdleTimeout(idleTimeout),
	)

	if err != nil {
		fmt.Printf("Ошибка запуска http сервера %s\n", err.Error())
		return
	}

	httpServer.Up(router)
}

// Инициализация БД
func createDatabaseConnection(env dotenv.IDotenv) db.IDatabase {
	DB, err := db.GetInstance(&db.DBConfig{
		Driver:             driver_enum.Postgresql,
		PostgresqlUser:     env.Getenv("POSTGRESQL_USER"),
		PostgresqlPassword: env.Getenv("POSTGRESQL_PASSWORD"),
		PostgresqlHost:     env.Getenv("POSTGRESQL_HOST"),
		PostgresqlPort:     env.Getenv("POSTGRESQL_PORT"),
		PostgresqlDBName:   env.Getenv("POSTGRESQL_DB_NAME"),
		PostgresqlSslMode:  env.Getenv("POSTGRESQL_SSL_MODE"),
	})

	// DB, err := db.GetInstance(&db.DBConfig{
	// 	Driver:     driver_enum.Sqlite,
	// 	SqlitePath: env.Getenv("SQLITE_PATH"),
	// })

	// DB, err := db.GetInstance(&db.DBConfig{
	// 	Driver:        driver_enum.Mysql,
	// 	MysqlUser:     env.Getenv("MYSQL_USER"),
	// 	MysqlPassword: env.Getenv("MYSQL_PASSWORD"),
	// 	MysqlHost:     env.Getenv("MYSQL_HOST"),
	// 	MysqlPort:     env.Getenv("MYSQL_PORT"),
	// 	MysqlDBName:   env.Getenv("MYSQL_DB_NAME"),
	// })

	if err != nil {
		fmt.Printf("Ошибка подключения к базе данных: %s\n", err)
		return nil
	}

	return DB
}

func main() {
	// Инициализация и парсинг .env
	env := dotenv.GetInstance()

	// Инициализация флагов
	command := flag.String("command", "", "Команда для выполнения")
	flag.Parse()

	// Инициализация БД (убрать если не используется) */
	DB := createDatabaseConnection(env)
	defer DB.Close()

	switch true {

	case *command != "":
		runCommand()
	default:
		runServer(env)
	}
}
