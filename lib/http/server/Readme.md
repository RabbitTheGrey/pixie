# Server #
### Package `http` ###
### Http сервер ###

## Настройка ##
В главном файле проекта `pixie.go` найдите функцию 

```go
func runServer(env dotenv.IDotenv)
```

В нем Вы увидите создание инстанса `IServer`.
Сам `Server` реализует singletone, который при инициализации принимает в конструкторе функциональные опции-сеттеры

```go
httpServer, err := server.GetInstance(
    server.ServerHost(serverHost),
    server.ServerPort(serverPort),
    server.ReadTimeout(readTimeout),
    server.WriteTimeout(writeTimeout),
    server.IdleTimeout(idleTimeout),
)
```

Рекомендую объявлять переменные `serverHost` и т.д. в окружении - `.env`, для получения использовать пакет `dotenv`

## Запуск ##
```go
httpServer.Up(router)
```

#### author @rabbitthegrey ####
contact: akrytar@gmail.com
