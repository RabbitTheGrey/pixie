package dotenv

import (
	"bufio"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

var instance *Dotenv
var once sync.Once

type Dotenv struct {
	variables map[string]string
}

// Инициализация dotenv
func GetInstance() IDotenv {
	once.Do(func() {
		instance = &Dotenv{variables: make(map[string]string)}
		instance.initVars()
	})

	return instance
}

func (de *Dotenv) Getenv(key string) string {
	value := de.variables[key]
	if value == "" {
		return os.Getenv(key)
	}

	return value
}

// Сохранение карты переменных
//
// Выполняется при инициализации!
func (de *Dotenv) initVars() {
	currentDir, err := os.Getwd()
	if err != nil {
		log.Fatalf("Ошибка получения корневой директории: %v\n", err)
		return
	}

	file, err := os.Open(filepath.Join(currentDir, "/.env"))
	if err != nil {
		log.Fatalf("Файл .env не найден: %v\n", err)
		return
	}

	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()

		// Отсекаем комментарии
		if strings.TrimSpace(line) == "" || strings.HasPrefix(line, "#") {
			continue
		}

		parts := strings.SplitN(line, "=", 2)
		if len(parts) == 2 {
			envKey := strings.TrimSpace(parts[0])
			envValue := strings.TrimSpace(parts[1])
			envValue = strings.Trim(envValue, `"'`)
			de.variables[envKey] = envValue
		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatalf("Ошибка чтения файла .env: %v\n", err)
		return
	}
}
