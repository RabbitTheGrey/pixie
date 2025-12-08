package db

import "database/sql"

type IDatabase interface {
	// Получить используемый драйвер
	GetDriver() string
	// Получить открытое соединение с базой данных
	GetConnection() *sql.DB
	// Закрыть соединение
	Close()
}
