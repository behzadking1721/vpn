package logging

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"sync"
	"time"
)

// LogLevel represents the severity level of a log message
type LogLevel int

const (
	DEBUG LogLevel = iota
	INFO
	WARNING
	ERROR
	FATAL
)

// String returns the string representation of a LogLevel
func (l LogLevel) String() string {
	switch l {
	case DEBUG:
		return "DEBUG"
	case INFO:
		return "INFO"
	case WARNING:
		return "WARNING"
	case ERROR:
		return "ERROR"
	case FATAL:
		return "FATAL"
	default:
		return "UNKNOWN"
	}
}

// Logger represents a logger instance
type Logger struct {
	level     LogLevel
	output    io.Writer
	file      *os.File
	mutex     sync.Mutex
	timestamp bool
}

// Config holds the configuration for the logger
type Config struct {
	Level     LogLevel
	Output    string // "stdout", "stderr", or file path
	Timestamp bool
}

// NewLogger creates a new logger instance
func NewLogger(config Config) (*Logger, error) {
	logger := &Logger{
		level:     config.Level,
		timestamp: config.Timestamp,
	}

	// Set output destination
	switch config.Output {
	case "stdout":
		logger.output = os.Stdout
	case "stderr":
		logger.output = os.Stderr
	default:
		// Treat as file path
		if config.Output != "" {
			// Create directory if it doesn't exist
			dir := filepath.Dir(config.Output)
			if err := os.MkdirAll(dir, 0755); err != nil {
				return nil, fmt.Errorf("failed to create log directory: %v", err)
			}

			// Open or create log file
			file, err := os.OpenFile(config.Output, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
			if err != nil {
				return nil, fmt.Errorf("failed to open log file: %v", err)
			}
			logger.file = file
			logger.output = file
		} else {
			// Default to stdout
			logger.output = os.Stdout
		}
	}

	return logger, nil
}

// Close closes the logger and any associated resources
func (l *Logger) Close() error {
	l.mutex.Lock()
	defer l.mutex.Unlock()

	if l.file != nil {
		return l.file.Close()
	}
	return nil
}

// SetLevel sets the logging level
func (l *Logger) SetLevel(level LogLevel) {
	l.mutex.Lock()
	defer l.mutex.Unlock()
	l.level = level
}

// Debug logs a debug message
func (l *Logger) Debug(format string, v ...interface{}) {
	if l.level <= DEBUG {
		l.log(DEBUG, format, v...)
	}
}

// Info logs an info message
func (l *Logger) Info(format string, v ...interface{}) {
	if l.level <= INFO {
		l.log(INFO, format, v...)
	}
}

// Warning logs a warning message
func (l *Logger) Warning(format string, v ...interface{}) {
	if l.level <= WARNING {
		l.log(WARNING, format, v...)
	}
}

// Error logs an error message
func (l *Logger) Error(format string, v ...interface{}) {
	if l.level <= ERROR {
		l.log(ERROR, format, v...)
	}
}

// Fatal logs a fatal message and exits the program
func (l *Logger) Fatal(format string, v ...interface{}) {
	l.log(FATAL, format, v...)
	os.Exit(1)
}

// log writes a log message with the specified level
func (l *Logger) log(level LogLevel, format string, v ...interface{}) {
	l.mutex.Lock()
	defer l.mutex.Unlock()

	// Get caller information
	_, file, line, ok := runtime.Caller(2)
	if !ok {
		file = "unknown"
		line = 0
	} else {
		// Get only the filename, not the full path
		file = filepath.Base(file)
	}

	// Format the message
	message := fmt.Sprintf(format, v...)

	// Add timestamp if enabled
	var timestamp string
	if l.timestamp {
		timestamp = time.Now().Format("2006-01-02 15:04:05")
	}

	// Format the complete log entry
	var logEntry string
	if l.timestamp {
		logEntry = fmt.Sprintf("[%s] [%s] [%s:%d] %s\n", timestamp, level.String(), file, line, message)
	} else {
		logEntry = fmt.Sprintf("[%s] [%s:%d] %s\n", level.String(), file, line, message)
	}

	// Write to output
	_, err := l.output.Write([]byte(logEntry))
	if err != nil {
		// If we can't write to the primary output, fallback to stderr
		log.Printf("Failed to write to log output: %v", err)
	}
}

// LogConnectionEvent logs a connection-related event
func (l *Logger) LogConnectionEvent(event string, serverID string, details map[string]interface{}) {
	if l.level > INFO {
		return
	}

	// Format details
	detailsStr := ""
	for key, value := range details {
		if detailsStr != "" {
			detailsStr += ", "
		}
		detailsStr += fmt.Sprintf("%s=%v", key, value)
	}

	if detailsStr != "" {
		l.Info("Connection Event: %s [Server: %s] (%s)", event, serverID, detailsStr)
	} else {
		l.Info("Connection Event: %s [Server: %s]", event, serverID)
	}
}

// LogServerError logs a server-related error
func (l *Logger) LogServerError(serverID string, err error, context string) {
	l.Error("Server Error [Server: %s] %s: %v", serverID, context, err)
}

// LogSubscriptionEvent logs a subscription-related event
func (l *Logger) LogSubscriptionEvent(event string, subscriptionID string, details map[string]interface{}) {
	if l.level > INFO {
		return
	}

	// Format details
	detailsStr := ""
	for key, value := range details {
		if detailsStr != "" {
			detailsStr += ", "
		}
		detailsStr += fmt.Sprintf("%s=%v", key, value)
	}

	if detailsStr != "" {
		l.Info("Subscription Event: %s [Subscription: %s] (%s)", event, subscriptionID, detailsStr)
	} else {
		l.Info("Subscription Event: %s [Subscription: %s]", event, subscriptionID)
	}
}