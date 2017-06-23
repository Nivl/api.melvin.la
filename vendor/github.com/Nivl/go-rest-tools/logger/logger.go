package logger

import (
	"fmt"
	"log"

	"github.com/bsphere/le_go"
)

// LogEntries represents a logger that sends data to logentries
var LogEntries *le_go.Logger

func doLog(msg string) {
	if LogEntries != nil {
		go LogEntries.Println(msg)
	}

	log.Println(msg)
}

// Errorf logs a formated error message
func Errorf(msg string, args ...interface{}) {
	fullMessage := fmt.Sprintf(`level: "ERROR", %s"`, fmt.Sprintf(msg, args...))
	doLog(fullMessage)
}

// Error logs an single error message
func Error(msg string) {
	fullMessage := fmt.Sprintf(`level: "ERROR", %s`, msg)
	doLog(fullMessage)
}
