# Middleware #
### Package `http` ###
### Подключение middleware ###

Middleware условно разделены на *глобальные* и *роутовые*.

Глобальные middleware автоматически распространяются на все маршруты роутера.

Роутовые подключаются при объявлении маршрута, как в примере выше:
```go
import (
	"pixie/handler/controller/example"
	"pixie/handler/middleware/routes_middleware"
	mw "pixie/lib/http/middleware"
	"pixie/lib/http/router"
)

func RegisterRoutes(r *router.Router) {
	r.AppendRoute("GET", "/example", example.List, []mw.Middleware{
        routes_middleware.NewExampleMiddleware(),
        // Остальные middleware...
    })
    // Остальные маршруты...
}
```

В директории `handler/middleware/routes_middleware` можно посмотреть пример готовой middleware

```go
type ExampleMiddleware struct{}

func NewExampleMiddleware() *ExampleMiddleware {
	return &ExampleMiddleware{}
}

func (mw *ExampleMiddleware) Handle(w http.ResponseWriter, r *http.Request, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Логика до
		next.ServeHTTP(w, r)
		// Логика после
	})
}

```
`next.ServeHTTP(w, r)` - пропускает нас к следующей middleware или к контроллеру, если вся цепочка middleware уже пройдена.

Middleware должна иметь уникальное имя и конструктор <b>New</b>Middleware.

Аналогично создаются и глобальные middleware, но объявляются они в `handler/middleware/global_middleware`, а передавать их нужно в главной функции `main.go` сразу после инициализации роутера, для этого используется строительный метод `WithGlobalMiddlewares`, в который передается срез middleware.

Выглядит это примерно так:
```go
router.WithGlobalMiddlewares(
    MyFirstMiddleware,
    MySecondMiddleware,
    // и так далее
)
```

Все middleware выполняются в порядке их объявления

Перед обработкой запроса методом контроллера выполняются сначала глобальные middlware, затем роутовые. После обработки завершаются в обратном порядке.

#### author @rabbitthegrey ####
contact: akrytar@gmail.com
