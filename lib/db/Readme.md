# Database #
### Работа с реляционными базами данных ###

На текущий момент реализована работа только с PostgreSQL, MySQL и SQLite.
Для работы потребуется драйвер для работы с БД, например `github.com/lib/pq`

## Инициализация БД ##

PostgreSQL:
```golang
// postgresql
DB, err := db.GetInstance(&db.DBConfig{
    Driver:             driver_enum.Postgresql,
    PostgresqlUser:     env.Getenv("POSTGRESQL_USER"),
    PostgresqlPassword: env.Getenv("POSTGRESQL_PASSWORD"),
    PostgresqlHost:     env.Getenv("POSTGRESQL_HOST"),
    PostgresqlPort:     env.Getenv("POSTGRESQL_PORT"),
    PostgresqlDBName:   env.Getenv("POSTGRESQL_DB_NAME"),
    PostgresqlSslMode:  env.Getenv("POSTGRESQL_SSL_MODE"),
})
```

SQLite:
```go
// sqlite
DB, err := db.GetInstance(&db.DBConfig{
    Driver:     driver_enum.Sqlite,
    SqlitePath: env.Getenv("SQLITE_PATH"),
})
```

MySQL:
```go
// mysql
DB, err := db.GetInstance(&db.DBConfig{
    Driver:        driver_enum.Mysql,
    MysqlUser:     env.Getenv("MYSQL_USER"),
    MysqlPassword: env.Getenv("MYSQL_PASSWORD"),
    MysqlHost:     env.Getenv("MYSQL_HOST"),
    MysqlPort:     env.Getenv("MYSQL_PORT"),
    MysqlDBName:   env.Getenv("MYSQL_DB_NAME"),
})

```

В конце необходимо закрыть полученное соединение `DB.Close()`

Параметр `Driver` влияет на формирование строки подключения - dsn и не зависит от выбранного ресурса с драйвером.

Объект DB является singletone, который реализует интерфейс `db.IDatabase`
Для дальнейшего получения инстанса, например, в репозиториях, можно в качестве `&db.DBConfig` указывать `nil`

`DB.GetConnection()` возвращает нативное подключение к БД *sql.DB, используемое в go из коробки.

## Datamapper ##

Если Вы хотите использовать структуры в качестве объектов БД, datamapper предоставляет для этого несколько методов:

```go
datamapper.SingleScalarResult(row *sql.Row, dest any) error
datamapper.SingleColumnResult(rows *sql.Rows, dest []any) error
datamapper.SingleResult(row *sql.Row, columns []string, dest any) error
datamapper.Result(rows *sql.Rows, dest any) error
```

Методы работают с указателями, поэтому вернуть могут только ошибку или `nil`.

Теперь по каждому отдельно:
- `SingleScalarResult` - обрабатывает полученное одиночное значение
```go
var now time.Time

row, err := db.GetConnection().QueryRow("SELECT NOW()")
// обработка ошибки err

datamapper.SingleScalarResult(row, *now)

// now: 2006-01-02 15:04:05.000
```

- `SingleColumnResult` - обрабатывает многострочный результат с одной выбранной коллонкой

```go
var ids []int

col, err := db.GetConnection().Query("SELECT id FROM users")
// обработка ошибки err

datamapper.SingleColumnResult(col, *ids)

// ids: []int{1, 2, 3, ...}
```

- `SingleResult` - обрабатывает однострочный результат с несколькими колонками, формирующими структуру

```go
type User struct {
    id    int
    email string
}

user := User{}

row, err := db.GetConnection().Query("SELECT id, email FROM users WHERE id = 1")
// обработка ошибки err
columns := []string{"id", "email"}

datamapper.SingleResult(row, columns, *user)

// user: {id: 1, email: "example@mail.ru"}
```

- `Result` - обрабатывает многострочный результат в срез структур

```go
type User struct {
    id    int
    email string
}

users := []User{}

rows, err := db.GetConnection().Query("SELECT id, email FROM users")
// обработка ошибки err
columns := []string{"id", "email"}

datamapper.SingleResult(rows, columns, *users)

// users: [{id: 1, email: "foo@mail.ru"}, {id: 2, email: "bar@mail.ru"}]
```

Для более явного сопоставления колонок со свойством структуры используется тэг `column`

```
type User struct {
    id        int       `column:"id"`
    email     string    `column:"email"`
    createdAt time.Time `column:"created_at"`
}
```

здесь переопределно поле createdAt, поскольку в БД не принято использовать camelCase

## Migrations ##

В `pixie/handler/command` зарегистрированы базовые команды для работы с миграциями:

- `go run pixie -command=migations_generate` - создать новую миграцию
- `go run pixie -command=migrations_migrate` - применить последние миграции
- `go run pixie -command=migrations_rollback` - откат последней миграции

migrations_generate - Создает в папке `/pixie/migrations` версионированный файл миграиции `Version...`
с объявленными внутри методами Up() и Down(), сразу же регистрирует миграцию в общей карте миграций
`/pixie/migrations/migrations.go`. Руками туда заносить ничего не нужно, можно только удалить неактуальные миграции.

Применение миграций проверяет наличие в базе таблицы `migration_versions`, и создает ее в случае, если та не нашлась.

В таблице содержится информация о примененных ранее миграциях. Важно сохранять имя миграции после создания, 
тк оно участвует в поиске последней миграции, что завязано на логике их применения и отката.
