package logger

import (
	"fmt"
	"runtime/debug"

	"github.com/Nivl/api.melvin.la/api/app"
)

func log(msg string) {
	context := app.GetContext()

	if context.LogEntries != nil {
		context.LogEntries.Println(msg)
	}

	fmt.Println(msg)
}

func Errorf(msg string, args ...interface{}) {
	fullMessage := fmt.Sprintf("%s | \"level\": \"ERROR\", %s", debug.Stack(), fmt.Sprintf(msg, args...))

	log(fullMessage)
}

func Error(msg string) {
	fullMessage := fmt.Sprintf("%s | \"level\": \"ERROR\", %s", debug.Stack(), msg)
	log(fullMessage)
}
