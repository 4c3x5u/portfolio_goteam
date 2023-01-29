// Package log contains code for logging messages on the server CLI for
// debugging purposes.
package log

// Level represents a specific log level.
type Level int8

const (
	// LevelError is a log Level that should be used when logging runtime errors
	// that are recoverable.
	LevelError Level = iota
	// LevelFatal is a log level that should be used when logging runtime errors
	// that are not recoverable where the app terminates after logging.
	LevelFatal
)

// Logger describes a type that can be used to log messages of different levels.
type Logger interface{ Log(Level, string) }
