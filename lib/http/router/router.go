package router

import (
	"pixie/lib/http/middleware"
	"net/http"
	"strings"
)

type Action func(w http.ResponseWriter, r *http.Request, params map[string]string)

type Route struct {
	Path        string                 // Путь запроса
	Method      string                 // HTTP метод (GET, POST, PUT, DELETE, etc.)
	Action      Action                 // Метод контроллера
	Middlewares *middleware.Collection // Миддлвары
}

type Router struct {
	routes            []Route                // Список зарегистрированных роутов
	globalMiddlewares *middleware.Collection // Зарегистрированные глобальные middleware
}

// Создание нового экземпляра роутера
func New() IRouter {
	return &Router{
		routes:            make([]Route, 0),
		globalMiddlewares: middleware.NewCollection(),
	}
}

func (r *Router) WithGlobalMiddlewares(middlewares ...middleware.Middleware) IRouter {
	for _, middleware := range middlewares {
		r.globalMiddlewares.Append(middleware)
	}

	return r
}

func (r *Router) AppendRoute(method string, path string, action Action, middlewares []middleware.Middleware) {
	normalizedPath := strings.Trim(path, "/")

	routeMiddlewares := middleware.NewCollection()
	for _, middleware := range middlewares {
		routeMiddlewares.Append(middleware)
	}

	r.routes = append(r.routes, Route{
		Path:        normalizedPath,
		Method:      method,
		Action:      action,
		Middlewares: routeMiddlewares,
	})
}

func (r *Router) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	path := strings.Trim(request.URL.Path, "/")
	hasMatch := false

	for _, route := range r.routes {
		params, ok := matchPath(route.Path, path)
		if ok {
			hasMatch = true
			if route.Method == request.Method {
				finalActionHandler := http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
					route.Action(writer, request, params)
				})

				// 1. Начинаем с роутовых middleware
				chain := route.Middlewares.BuildChain(finalActionHandler)
				// 2. Оборачиваем их глобальными middleware
				chain = r.globalMiddlewares.BuildChain(chain)

				chain.ServeHTTP(writer, request)
				return
			}
		}
	}

	// Нашлось совпадение пути, но не совпал метод
	if hasMatch {
		http.Error(writer, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	http.NotFound(writer, request)
}

// Проверка соответствия пути зарегистрированному маршруту с подстановкой параметров
func matchPath(pattern string, path string) (map[string]string, bool) {
	patternParts := strings.Split(pattern, "/")
	pathParts := strings.Split(path, "/")

	if len(patternParts) != len(pathParts) {
		return nil, false
	}

	params := make(map[string]string)

	for i := range patternParts {
		if strings.HasPrefix(patternParts[i], "{") && strings.HasSuffix(patternParts[i], "}") {
			paramName := strings.TrimSuffix(strings.TrimPrefix(patternParts[i], "{"), "}")
			params[paramName] = pathParts[i]
		} else if patternParts[i] != pathParts[i] {
			return nil, false
		}
	}

	return params, true
}
