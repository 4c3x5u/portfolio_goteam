// Package log contains code for logging messages on the server CLI for
// debugging purposes.
package log

import "log"

// Log can be used to log messages of different log levels across the project.
type Log struct{}

// New creates and returns a new Log.
func New() Log { return Log{} }

// Info logs an information-level message to the console.
func (l Log) Info(msg string) { log.Println("--[INFO]-- " + msg) }

// Info logs an error-level message to the console.
func (l Log) Error(msg string) { log.Println("--[ERROR]-- " + msg) }

// Info logs an fatal-level message to the console.
func (l Log) Fatal(msg string) { log.Println("--[FATAL]-- " + msg) }
