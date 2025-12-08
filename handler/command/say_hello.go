package command

import (
	"fmt"
	"pixie/lib/console"
)

func SayHelloCommand(args map[string]string) int {
	fmt.Println("hello")
	return console.Success
}
