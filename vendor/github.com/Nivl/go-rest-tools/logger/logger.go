package logger

// Logger is an interface used for all loggers
type Logger interface {
	AddStaticData(msg string, args ...interface{})
	Errorf(msg string, args ...interface{})
	Error(msg string)
	Close() error
}
