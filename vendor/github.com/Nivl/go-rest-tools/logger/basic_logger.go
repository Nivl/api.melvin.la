package logger

import (
	"fmt"
	"log"
)

// BasicLogger represents a logger that just uses the default log package
type BasicLogger struct {
	staticData []string
}

// AddStaticData is used to add static data to the logs.
// static data will be added to all logs
func (bl *BasicLogger) AddStaticData(msg string, args ...interface{}) {
	bl.staticData = append(bl.staticData, fmt.Sprintf(msg, args...))
}

// Errorf is used to log an error with formating
func (bl *BasicLogger) Errorf(msg string, args ...interface{}) {
	bl.Send(formatError(bl.staticData, fmt.Sprintf(msg, args...)))
}

// Error is used to log an error
func (bl *BasicLogger) Error(msg string) {
	bl.Send(formatError(bl.staticData, msg))
}

// Send is used to log a raw message
func (bl *BasicLogger) Send(msg string) {
	log.Println(msg)
}

// Close closes the connection to BasicLogger
func (bl *BasicLogger) Close() error {
	return nil
}

// NewBasicLogger creates an instance of BasicLogger
func NewBasicLogger() *BasicLogger {
	return &BasicLogger{}
}
