package console

type IConsole interface {
	// Добавление команды в пулл
	AppendCommand(name string, action Action)

	// Выполнение команды, переданной значением аргумента --command в main.go
	//
	// Возвращает статус выполнения команды:
	//  * 0 - Success (Успешное выполнение)
	//  * 1 - Failure (Завершилась с ошибкой)
	//  * 2 - Invalid (Некорректные входные данные)
	Execute() int
}
