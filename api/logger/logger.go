package logger

import (
	"fmt"
	"log"
	"runtime/debug"

	"github.com/Nivl/api.melvin.la/api/app"
)

func logg(msg string) {
	context := app.GetContext()

	if context != nil && context.LogEntries != nil {
		context.LogEntries.Println(msg)
	}

	log.Println(msg)
}

func Errorf(msg string, args ...interface{}) {
	fullMessage := fmt.Sprintf("%s | \"level\": \"ERROR\", %s", debug.Stack(), fmt.Sprintf(msg, args...))
	logg(fullMessage)
}

func Error(msg string) {
	fullMessage := fmt.Sprintf("%s | \"level\": \"ERROR\", %s", debug.Stack(), msg)
	logg(fullMessage)
}
