package logger

import (
	"fmt"
	"log"
	"runtime/debug"

	"github.com/melvin-laplanche/ml-api/src/app"
)

func logg(msg string) {
	context := app.GetContext()

	if context != nil && context.LogEntries != nil {
		go context.LogEntries.Println(msg)
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
