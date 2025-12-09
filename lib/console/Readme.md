# Console #
### Командная строка ###

## Использоавание. Обработчики команд ##

Команды создаются в директории `handler/command`.

Для структурного разделения рекомендую создавать команды в одноименной папке, например `example/example.go`

Пример реализованных команд `handler/command/say_hello.go`

Методы-обработчики в команде должны реализовывать абстрактную функцию `console.Action`:

```go
type Action func(args map[string]string) int
```

На выходе команда должна отдавать одно из целочиленных значений:
- 0 - Успешное выполнение команды
- 1 - Команда завершилась с ошибкой
- 2 - Некорректное поведение 

Маршрутизация настраивается в `handler/routes.go` в методе `RegisterRoutes`

Пример маршрута:

```go
import (
	"pixie/handler/controller/example"
	mw "pixie/lib/http/middleware"
	"pixie/lib/http/router"
)

func RegisterRoutes(r *router.Router) {
	r.AppendRoute("GET", "/example", example.List, []mw.Middleware{})
    // Остальные маршруты...
}
```
Здесь первым аргументом передается http метод, вторым - url запроса, третьим - метод Вашего контроллера, последний - используемые middleware (по умолчанию пустой массив mw.Middleware)

## author @rabbitthegrey ####
contact: akrytar@gmail.com
