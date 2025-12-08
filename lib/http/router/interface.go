package router

import (
	"pixie/lib/http/middleware"
	"net/http"
)

type IRouter interface {
	// Добавление глобальных middleware
	WithGlobalMiddlewares(middlewares ...middleware.Middleware) IRouter
	// Добавление нового роута
	AppendRoute(method string, path string, action Action, middlewares []middleware.Middleware)
	// Метод релизует http.Handler, что позволяет использовать Router как обработчик в http.Server
	//
	//  На этом этапе обработаны коды ответа:
	//  * 404 - not fount
	//  * 405 - method not allowed
	ServeHTTP(writer http.ResponseWriter, request *http.Request)
}
