package log

import "log"

// AppLogger can be used to log messages of different log levels across the
// project.
type AppLogger struct{}

// NewAppLogger creates and returns a new BasicLogger.
func NewAppLogger() AppLogger { return AppLogger{} }

// Log formats a message based on LogLevel and logs the resulting string.
func (l AppLogger) Log(level Level, message string) {
	switch level {
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
