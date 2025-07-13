package taskrunner

import (
	"fmt"
)

type logger interface {
	Info(string, ...any)
	Error(string, ...any)
	TaskDescription(string, ...any)
}

const (
	green    = "\033[32m"
	red      = "\033[31m"
	blueBold = "\033[1;34m"
	reset    = "\033[0m"
)

type consoleLogger struct{}

func (consoleLogger) Info(f string, a ...any) {
	fmt.Printf(green+f+reset+"\n", a...)
}

func (consoleLogger) Error(f string, a ...any) {
	fmt.Printf(red+f+reset+"\n", a...)
}

func (consoleLogger) TaskDescription(f string, a ...any) {
	title := fmt.Sprintf(f, a...)
	fmt.Printf("\n"+blueBold+"==== %s ===="+reset+"\n\n", title)
}
