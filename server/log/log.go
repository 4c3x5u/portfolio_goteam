// Package log contains code for logging messages on the server terminal for
// debugging purposes.
package log

import (
	"log"
	"os"
)

// Level represents a specific log level.
type Level int64

// Level can be used to distinguish between different log levels.
const (
	LevelError Level = iota
	LevelFatal
)

// Logger describes a type that can be used to logs messages of different log
// levels.
type Logger interface{ Log(Level, string) }

// AppLogger can be used to logs messages of different log levels across
// the project.
type AppLogger struct{}

// NewAppLogger creates and returns a new BasicLogger.
func NewAppLogger() AppLogger { return AppLogger{} }

// Log formats a message based on LogLevel and logs the resulting string.
func (l AppLogger) Log(level Level, message string) {
	switch level {
	case LevelError:
		log.Println("[ERROR]" + message)
	case LevelFatal:
		log.Println("[FATAL]" + message)
		os.Exit(1)
	default:
		log.Println("[WARNING] log level invalid for message: \n  " + message)
	}
}
