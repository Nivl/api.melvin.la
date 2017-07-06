package logger

import (
	"fmt"
	"log"

	"github.com/bsphere/le_go"
)

// NewLogEntries creates a new logger that will push a logentries logger
func NewLogEntries(driver *le_go.Logger) *LogEntries {
	return &LogEntries{
		driver: driver,
	}
}

// LogEntries represents a logger that sends data to logentries
type LogEntries struct {
	driver     *le_go.Logger
	staticData []string
}

// AddStaticData is used to add static data to the logs.
// static data will be added to all logs
func (le *LogEntries) AddStaticData(msg string, args ...interface{}) {
	le.staticData = append(le.staticData, fmt.Sprintf(msg, args...))
}

// Errorf is used to log an error with formating
func (le *LogEntries) Errorf(msg string, args ...interface{}) {
	le.Send(formatError(le.staticData, fmt.Sprintf(msg, args...)))
}

// Error is used to log an error
func (le *LogEntries) Error(msg string) {
	le.Send(formatError(le.staticData, msg))
}

// Send is used to log a raw message
func (le *LogEntries) Send(msg string) {
	go le.driver.Println(msg)
	log.Println(msg)
}

// Close closes the connection to logentries
func (le *LogEntries) Close() error {
	return le.driver.Close()
}
