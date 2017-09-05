// This file handles how debug and trace messages get passed to stdout.

package echo

import (
	"log"
	"strings"
)

// Levels for implementing the debug and trace message functionality.
const (
	Trace uint8 = iota
	Debug
	Info
	Status
	Warn
	Silent
)

// These variables are initialized in init()
var logLevel = Debug
var logger *log.Logger
var logLevelStrings = [...]string{"trace", "debug", "info", "status", "warn", "silent"}

//===========================================================================
// Interact with debug output
//===========================================================================

// LogLevel returns a string representation of the current level
func LogLevel() string {
	return logLevelStrings[logLevel]
}

// SetLogLevel modifies the log level for messages at runtime. Ensures that
// the highest level that can be set is the trace level.
func SetLogLevel(level uint8) {
	if level > Silent {
		level = Silent
	}

	logLevel = level
}

//===========================================================================
// Debugging output functions
//===========================================================================

// Print to the standard logger at the specified level. Arguments are handled
// in the manner of log.Printf, but a newline is appended.
func print(level uint8, msg string, a ...interface{}) {
	if level >= logLevel {
		if !strings.HasSuffix(msg, "\n") {
			msg += "\n"
		}

		logger.Printf(msg, a...)
	}
}

// Prints to the standard logger if level is warn or greater; arguments are
// handled in the manner of log.Printf, but a newline is appended.
func warn(msg string, a ...interface{}) {
	print(Warn, msg, a...)
}

// Helper function to simply warn about an error received.
func warne(err error) {
	warn(err.Error())
}

// Prints to the standard logger if level is status or greater; arguments are
// handled in the manner of log.Printf, but a newline is appended.
func status(msg string, a ...interface{}) {
	print(Status, msg, a...)
}

// Prints to the standard logger if level is info or greater; arguments are
// handled in the manner of log.Printf, but a newline is appended.
func info(msg string, a ...interface{}) {
	print(Info, msg, a...)
}

// Prints to the standard logger if level is debug or greater; arguments are
// handled in the manner of log.Printf, but a newline is appended.
func debug(msg string, a ...interface{}) {
	print(Debug, msg, a...)
}

// Prints to the standard logger if level is trace or greater; arguments are
// handled in the manner of log.Printf, but a newline is appended.
func trace(msg string, a ...interface{}) {
	print(Trace, msg, a...)
}
