package logger

import (
	"fmt"
	"log"

	"github.com/melvin-laplanche/ml-api/src/app"
)

func logg(msg string) {
	context := app.GetContext()

	if context != nil && context.LogEntries != nil {
		go context.LogEntries.Println(msg)
	}

	log.Println(msg)
}

// Errorf logs a formated error message
func Errorf(msg string, args ...interface{}) {
	fullMessage := fmt.Sprintf(`level: "ERROR", %s"`, fmt.Sprintf(msg, args...))
	logg(fullMessage)
}

// Error logs an single error message
func Error(msg string) {
	fullMessage := fmt.Sprintf(`level: "ERROR", %s`, msg)
	logg(fullMessage)
}
