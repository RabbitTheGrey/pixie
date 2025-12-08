package dotenv

type IDotenv interface {
	// Получение переменной из .env корневой директории по ключу
	//  Если переменная не найдена, функция попытается найти ее в os.Getenv()
	//  Если переменная не найдена в os, вернет пустую строку
	Getenv(key string) string
}
