package log

import "log"

// Level represents a specific log level.
type Level int8

const (
	// LevelError is a log Level that should be for logging messages that are not
	// errors.
	LevelInfo Level = iota
	// LevelError is a log Level that should be used when logging runtime errors
	// that are recoverable.
	LevelError
	// LevelFatal is a log level that should be used when logging runtime errors
	// that are not recoverable where the app terminates after logging.
	LevelFatal
)

// Logger describes a type that can be used to log messages of different levels.
type Logger interface{ Log(Level, string) }

// AppLogger can be used to log messages of different log levels across the
// project.
type AppLogger struct{}

// NewAppLogger creates and returns a new AppLogger.
func NewAppLogger() AppLogger { return AppLogger{} }

// Log formats a message based on LogLevel and logs the resulting string.
func (l AppLogger) Log(level Level, message string) {
	switch level {
	case LevelInfo:
		log.Println("--[INFO]-- " + message)
	case LevelError:
		log.Println("--[ERROR]-- " + message)
	case LevelFatal:
		log.Println("--[FATAL]-- " + message)
	default:
		log.Println(
			"--[WARNING]-- invalid log level for message: \n  " + message,
		)
	}
}
