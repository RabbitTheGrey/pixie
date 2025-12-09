# Pixie #
### Мини фреймворк, реализующий базовый функционал для обработки http запросов ###

На текущий момент реализовано:
- <a href="https://github.com/RabbitTheGrey/pixie/blob/master/lib/dotenv/Readme.md">парсер .env файла</a>
- <a href="https://github.com/RabbitTheGrey/pixie/blob/master/lib/http/server/Readme.md">запуск http сервера</a>
- <a href="https://github.com/RabbitTheGrey/pixie/blob/master/lib/http/router/Readme.md">роутер</a>
- <a href="https://github.com/RabbitTheGrey/pixie/blob/master/lib/http/middleware/Readme.md">middleware</a>
- <a href="https://github.com/RabbitTheGrey/pixie/blob/master/lib/console/Readme.md">командная строка</a>
- <a href="https://github.com/RabbitTheGrey/pixie/blob/master/lib/db/Readme.md">базы данных (миграции, datamapper)</a>

В разработке:
- валидаторы
- авторизация и аутентификация
- сессии
- некоторые глобальные middleware (добавление заголовков, логирование, передача текущего пользователя в запрос)
- работа с токенами CSRF, JWT
- поддержка gRPC
- рефакторинг существующих пакетов

## Технические требования ##
go ^1.25.2

DB driver (postgresql/mysql/sqlite)

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

#### author @rabbitthegrey ####
contact: akrytar@gmail.com
