package console

import (
	"flag"
	"fmt"
)

const (
	Success int = 0
	Failure int = 1
	Invalid int = 2
)

type Action func(args map[string]string) int

type Command struct {
	Name   string // Уникальное имя команды
	Action Action // Метод обработчика
}

type Console struct {
	commands []Command
}

// Создание нового экземпляра Console
func New() IConsole {
	return &Console{commands: make([]Command, 0)}
}

func (c *Console) AppendCommand(name string, action Action) {
	c.commands = append(c.commands, Command{
		Name:   name,
		Action: action,
	})
}

func (c *Console) Execute() int {
	name := flag.Lookup("command").Value.String()

	for _, command := range c.commands {
		if command.Name == name {
			return command.Action(c.getArgs())
		}
	}

	fmt.Printf("Команда \"%s\" не найдена \n", name)
	return Invalid
}

// Извлечение аргументов команды
func (c *Console) getArgs() map[string]string {
	args := make(map[string]string)

	flag.Visit(func(arg *flag.Flag) {
		args[arg.Name] = arg.Value.String()
	})

	return args
}
