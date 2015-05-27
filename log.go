package log

import (
	"errors"
	"fmt"
	"io"
)

var loggers []Logger

// Logger is an interface that should be implemented by all loggers wishing to participate
// in the logger chain initialized by this package.
type Logger interface {
	Info(string, ...interface{})
	Warning(string, ...interface{})
	Error(string, ...interface{})

	// Writer returns a logger's underlying io.Writer used to write log messages to.
	//
	// It may be, for example, the standard output for a console logger or a socket
	// connection for a UDP logger.
	Writer(Severity) io.Writer

	// FormatMessage constructs and returns a final message that will go a logger's
	// output channel.
	FormatMessage(Severity, *callerInfo, string, ...interface{}) string

	String() string
}

// LogConfig represents a configuration of an individual logger.
type LogConfig struct {
	// Name is a logger's identificator used to instantiate a proper logger type
	// from a config.
	Name string

	// Severity indicates the minimum severity a logger will be logging messages at.
	Severity string
}

func (c LogConfig) String() string {
	return fmt.Sprintf("LogConfig(Name=%v, Severity=%v)", c.Name, c.Severity)
}

// Init initializes the logging package with the provided loggers.
func Init(l ...Logger) {
	for _, logger := range l {
		loggers = append(loggers, logger)
	}
}

// InitWithConfig instantiates loggers based on the provided configs and initializes
// the package with them.
func InitWithConfig(logConfigs ...LogConfig) error {
	for _, config := range logConfigs {
		l, err := NewLogger(config)
		if err != nil {
			return err
		}
		loggers = append(loggers, l)
	}
	return nil
}

// NewLogger makes a proper logger from the given configuration.
func NewLogger(config LogConfig) (Logger, error) {
	switch config.Name {
	case ConsoleLoggerName:
		return NewConsoleLogger(config)
	case SysLoggerName:
		return NewSysLogger(config)
	case UDPLoggerName:
		return NewUDPLogger(config)
	}
	return nil, errors.New(fmt.Sprintf("unknown logger: %v", config))
}

// Info logs to the INFO log.
func Info(format string, args ...interface{}) {
	for _, logger := range loggers {
		writeMessage(logger, 1, SeverityInfo, format, args...)
	}
}

// Warning logs to the WARN and INFO logs.
func Warning(format string, args ...interface{}) {
	for _, logger := range loggers {
		writeMessage(logger, 1, SeverityWarning, format, args...)
	}
}

// Error logs to the ERROR, WARN, and INFO logs.
func Error(format string, args ...interface{}) {
	for _, logger := range loggers {
		writeMessage(logger, 1, SeverityError, format, args...)
	}
}

func writeMessage(logger Logger, callDepth int, sev Severity, format string, args ...interface{}) {
	caller := getCallerInfo(callDepth + 1)
	if w := logger.Writer(sev); w != nil {
		message := logger.FormatMessage(sev, caller, format, args...)
		io.WriteString(w, message)
	}
}
