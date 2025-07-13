package tr

import (
	"fmt"
)

type logger interface {
	Info(string, ...any)
	Error(string, ...any)
}

type consoleLogger struct{}

func (consoleLogger) Info(f string, a ...any)  { fmt.Printf("\033[32m"+f+"\033[0m", a...) }
func (consoleLogger) Error(f string, a ...any) { fmt.Printf("\033[31m"+f+"\033[0m", a...) }

var logImpl logger = consoleLogger{}
