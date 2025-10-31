package utils

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"
)

// LogLevel represents the severity level of a log message
type LogLevel int

const (
	// LogLevelDebug represents debug level messages
	LogLevelDebug LogLevel = iota
	// LogLevelInfo represents info level messages
	LogLevelInfo
	// LogLevelWarn represents warning level messages
	LogLevelWarn
	// LogLevelError represents error level messages
	LogLevelError
)

// String returns the string representation of a LogLevel
func (l LogLevel) String() string {
	switch l {
	case LogLevelDebug:
		return "DEBUG"
	case LogLevelInfo:
		return "INFO"
	case LogLevelWarn:
		return "WARN"
	case LogLevelError:
		return "ERROR"
	default:
		return "UNKNOWN"
	}
}

// Logger represents a logger instance
type Logger struct {
	level     LogLevel
	logger    *log.Logger
	file      *os.File
	timestamp bool
}

// LoggerConfig holds configuration for creating a new logger
type LoggerConfig struct {
	// Level is the minimum log level to output
	Level LogLevel
	// Output is the output destination (default is os.Stdout)
	Output io.Writer
	// File is the file path to write logs to (optional)
	File string
	// Timestamp indicates whether to include timestamps in logs
	Timestamp bool
	// Prefix is an optional prefix for all log messages
	Prefix string
}

// Global logger instance
var globalLogger *Logger

// init initializes the global logger
func init() {
	globalLogger = NewLogger(&LoggerConfig{
		Level:     LogLevelInfo,
		Timestamp: true,
	})
}

// NewLogger creates a new logger instance
func NewLogger(config *LoggerConfig) *Logger {
	if config == nil {
		config = &LoggerConfig{}
	}

	// Set default values
	if config.Output == nil {
		config.Output = os.Stdout
	}

	if config.Timestamp {
		config.Prefix = fmt.Sprintf("%s %s", time.Now().Format("2006-01-02 15:04:05"), config.Prefix)
	}

	// If file is specified, create directory and open file
	var file *os.File
	if config.File != "" {
		// Create directory if it doesn't exist
		dir := filepath.Dir(config.File)
		if err := os.MkdirAll(dir, 0755); err != nil {
			log.Printf("Failed to create log directory: %v", err)
		}

		// Open file for writing (append mode)
		var err error
		file, err = os.OpenFile(config.File, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
		if err != nil {
			log.Printf("Failed to open log file: %v", err)
			file = nil
		} else {
			// Use both stdout and file as output
			config.Output = io.MultiWriter(config.Output, file)
		}
	}

	logger := &Logger{
		level:     config.Level,
		logger:    log.New(config.Output, config.Prefix, 0),
		file:      file,
		timestamp: config.Timestamp,
	}

	return logger
}

// SetGlobalLogger sets the global logger instance
func SetGlobalLogger(logger *Logger) {
	globalLogger = logger
}

// Debug logs a debug message
func (l *Logger) Debug(format string, v ...interface{}) {
	if l.level <= LogLevelDebug {
		l.log(LogLevelDebug, format, v...)
	}
}

// Info logs an info message
func (l *Logger) Info(format string, v ...interface{}) {
	if l.level <= LogLevelInfo {
		l.log(LogLevelInfo, format, v...)
	}
}

// Warn logs a warning message
func (l *Logger) Warn(format string, v ...interface{}) {
	if l.level <= LogLevelWarn {
		l.log(LogLevelWarn, format, v...)
	}
}

// Error logs an error message
func (l *Logger) Error(format string, v ...interface{}) {
	if l.level <= LogLevelError {
		l.log(LogLevelError, format, v...)
	}
}

// Debug logs a debug message using the global logger
func Debug(format string, v ...interface{}) {
	globalLogger.Debug(format, v...)
}

// Info logs an info message using the global logger
func Info(format string, v ...interface{}) {
	globalLogger.Info(format, v...)
}

// Warn logs a warning message using the global logger
func Warn(format string, v ...interface{}) {
	globalLogger.Warn(format, v...)
}

// Error logs an error message using the global logger
func Error(format string, v ...interface{}) {
	globalLogger.Error(format, v...)
}

// log is the internal logging function
func (l *Logger) log(level LogLevel, format string, v ...interface{}) {
	// Get caller information
	_, file, line, ok := runtime.Caller(2)
	if ok {
		// Extract just the filename
		file = filepath.Base(file)
		format = fmt.Sprintf("[%s] %s:%d %s", level, file, line, format)
	} else {
		format = fmt.Sprintf("[%s] %s", level, format)
	}

	// Add timestamp if needed
	if l.timestamp {
		format = fmt.Sprintf("%s %s", time.Now().Format("2006-01-02 15:04:05"), format)
	}

	// Log the message
	l.logger.Printf(format, v...)
}

// Close closes the logger and any associated resources
func (l *Logger) Close() error {
	if l.file != nil {
		return l.file.Close()
	}
	return nil
}

// Rotate creates a new log file and closes the current one
func (l *Logger) Rotate() error {
	if l.file == nil {
		return nil
	}

	// Close current file
	if err := l.file.Close(); err != nil {
		return err
	}

	// Get file name
	name := l.file.Name()

	// Create backup name with timestamp
	backupName := fmt.Sprintf("%s.%s", name, time.Now().Format("20060102-150405"))

	// Rename current file to backup
	if err := os.Rename(name, backupName); err != nil {
		return err
	}

	// Open new file
	newFile, err := os.OpenFile(name, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return err
	}

	// Update logger
	l.file = newFile
	l.logger.SetOutput(newFile)

	return nil
}

// SetLevel sets the logging level
func (l *Logger) SetLevel(level LogLevel) {
	l.level = level
}

// GetLevel returns the current logging level
func (l *Logger) GetLevel() LogLevel {
	return l.level
}
