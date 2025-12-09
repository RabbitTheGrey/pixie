# Router #
### Package `http` ###
### Маршрутизация ###

## Использоавание. Контроллеры ##

Контроллеры создаются в директории `handler/controller`

Для структурного разделения рекомендую создавать контроллер в одноименной папке, например `example/example.go`

Пример реализованного контроллера `handler/controller/example/example.go`

Методы в контроллере должны реализовывать абстрактную функцию `router.Action`:

```go
type Action func(w http.ResponseWriter, r *http.Request, params map[string]string)
```
где r,w - стандартные структуры запроса и вывода, params - карта параметров, передаваемых в фигурных скобках маршрута

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
