// Package log contains code for logging messages on the server CLI for
// debugging purposes.
package log

import "log"

// Errorer describes a type that can be used to log an error-level message to
// the console.
type Errorer interface{ Error(string) }

// Log can be used to log messages of different log levels across the project.
type Log struct{}

// New creates and returns a new Log.
func New() Log { return Log{} }

// Info logs an information-level message to the console.
func (l Log) Info(message string) { log.Println("--[INFO]-- " + message) }

// Info logs an error-level message to the console.
func (l Log) Error(message string) { log.Println("--[ERROR]-- " + message) }

// Info logs an fatal-level message to the console.
func (l Log) Fatal(message string) { log.Println("--[FATAL]-- " + message) }
