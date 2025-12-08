package server

import "net/http"

type IServer interface {
	// Запуск Http сервера, прослушивание входящих соединений
	Up(router http.Handler)
}
