# Pixie #
### Мини фреймворк, реализующий базовый функционал для обработки http запросов ###

На текущий момент реализовано:
- парсер .env файла
- запуск http сервера
- роутер
- middleware

В разработке:
- ORM (миграции, data mapper)
- валидаторы
- авторизация и аутентификация
- сессии
- глобальные middleware (добавление заголовков, логирование, передача текущего пользователя в запрос)
- работа с токенами CSRF, JWT

#### author @rabbitthegrey ####

## Технические требования ##
go ^1.25.2

## Использование ##
### Окружение ###

Для объявления переменных окружения создайте файл `.env` в корневой директории и перенесите содержимое `.env.example`.
В нем объявлены необходимые переменные для запуска сервера, замените
```
SERVER_HOST=<Ваш адрес сервера>
SERVER_PORT=<Открытый tcp порт, например 80>
```
на ip и port Вашей машины.

Таймауты по умолчанию составляют 1 минуту. Можно их оставить без изменений, если нет потребности в более длительном ожидании ответа сервера.

Остальные переменные добавляются на Ваше усмотрение.

```go
// Пример получения переменной в коде программы:
import "gofman/lib/dotenv"

func main() {
    token := dotenv.Getenv("JWT_TOKEN")
}
```

### Создание контроллера ###

Контроллеры создаются в директории `handler/controller`.

Для структурного разделения рекомендую создавать контроллер в одноименной папке, например `example/example.go`

Пример реализованного контроллера `handler/controller/example/example.go`

Методы в контроллере должны реализовывать абстрактную функцию `Action`:

```go
type Action func(w http.ResponseWriter, r *http.Request, params map[string]string)
```
где r,w - стандартные структуры запроса и вывода, params - карта параметров, передаваемых в фигурных скобках маршрута

Маршрутизация настраивается в `handler/routes.go` в методе `RegisterRoutes`

Пример маршрута:

```go
import (
	"gofman/handler/controller/example"
	mw "gofman/lib/http/middleware"
	"gofman/lib/http/router"
)

func RegisterRoutes(r *router.Router) {
	r.AppendRoute("GET", "/example", example.List, []mw.Middleware{})
    // Остальные маршруты...
}
```
Здесь первым аргументом передается http метод, вторым - url запроса, третьим - метод Вашего контроллера, последний - используемые middleware (по умолчанию пустой массив mw.Middleware)

### Подключение middleware ###

Middleware условно разделены на *глобальные* и *роутовые*.

Глобальные middleware автоматически распространяются на все маршруты роутера.

Роутовые подключаются при объявлении маршрута, как в примере выше:
```go
import (
	"gofman/handler/controller/example"
	"gofman/handler/middleware/routes_middleware"
	mw "gofman/lib/http/middleware"
	"gofman/lib/http/router"
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
Middleware должна иметь уникальное имя и конструктор

Аналогично создаются и глобальные middleware, но объявлять их нужно в главной функции `main.go` в роутере сразу после его инициализации, для этого используется строительный метод `WithGlobalMiddleware`.

Выглядит это примерно так:
```go
router.
	WithGlobalMiddleware(MyFirstMiddleware).
	WithGlobalMiddleware(MySecondMiddleware).
    WithGlobalMiddleware(MyThirdMiddleware).
    // и так далее
```

Все middleware выполняются в порядке их объявления

Перед обработкой запроса методом контроллера выполняются сначала глобальные middlware, затем роутовые. После обработки завершаются в обратном порядке.

### ORM (миграции, data mapper) ###
### валидаторы ###
### авторизация и аутентификация ###
### сессии ###
### глобальные middleware (добавление заголовков, логирование, передача текущего пользователя в запрос) ###
### работа с токенами CSRF, JWT ###