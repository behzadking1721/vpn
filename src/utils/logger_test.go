package utils

import (
	"bytes"
	"io/ioutil"
	"os"
	"strings"
	"testing"
)

func TestLoggerCreation(t *testing.T) {
	// Test creating a logger with default config
	logger := NewLogger(nil)
	if logger == nil {
		t.Error("Expected logger to be created")
	}

	// Test creating a logger with custom config
	config := &LoggerConfig{
		Level:     LogLevelDebug,
		Timestamp: false,
		Prefix:    "TEST",
	}
	logger = NewLogger(config)
	if logger == nil {
		t.Error("Expected logger to be created with custom config")
	}
}

func TestLoggerLevels(t *testing.T) {
	var buf bytes.Buffer

	// Create a logger that writes to buffer
	config := &LoggerConfig{
		Level:     LogLevelInfo,
		Output:    &buf,
		Timestamp: false,
	}
	logger := NewLogger(config)

	// Debug messages should not be logged
	logger.Debug("This is a debug message")
	if buf.Len() > 0 {
		t.Error("Debug message should not be logged when level is INFO")
	}

	// Info messages should be logged
	logger.Info("This is an info message")
	if buf.Len() == 0 {
		t.Error("Info message should be logged when level is INFO")
	}

	// Clear buffer
	buf.Reset()

	// Warn messages should be logged
	logger.Warn("This is a warn message")
	if buf.Len() == 0 {
		t.Error("Warn message should be logged when level is INFO")
	}

	// Clear buffer
	buf.Reset()

	// Error messages should be logged
	logger.Error("This is an error message")
	if buf.Len() == 0 {
		t.Error("Error message should be logged when level is INFO")
	}
}

func TestLoggerWithFile(t *testing.T) {
	// Create a temporary log file
	tmpFile, err := ioutil.TempFile("", "vpn_test_*.log")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())
	defer tmpFile.Close()

	// Create a logger with file output
	config := &LoggerConfig{
		Level:     LogLevelDebug,
		File:      tmpFile.Name(),
		Timestamp: false,
	}
	logger := NewLogger(config)
	if logger == nil {
		t.Error("Expected logger to be created with file output")
	}

	// Log some messages
	logger.Info("Test info message")
	logger.Error("Test error message")

	// Close the logger to flush to file
	logger.Close()

	// Read the file content
	content, err := ioutil.ReadFile(tmpFile.Name())
	if err != nil {
		t.Fatalf("Failed to read log file: %v", err)
	}

	// Check if messages were written to file
	contentStr := string(content)
	if !strings.Contains(contentStr, "Test info message") {
		t.Error("Info message not found in log file")
	}

	if !strings.Contains(contentStr, "Test error message") {
		t.Error("Error message not found in log file")
	}
}

func TestGlobalLogger(t *testing.T) {
	var buf bytes.Buffer

	// Create a logger that writes to buffer
	config := &LoggerConfig{
		Level:     LogLevelDebug,
		Output:    &buf,
		Timestamp: false,
	}
	logger := NewLogger(config)

	// Set as global logger
	SetGlobalLogger(logger)

	// Test global logger functions
	Debug("Global debug message")
	Info("Global info message")
	Warn("Global warn message")
	Error("Global error message")

	// Check if messages were logged
	content := buf.String()
	if !strings.Contains(content, "Global debug message") {
		t.Error("Global debug message not logged")
	}

	if !strings.Contains(content, "Global info message") {
		t.Error("Global info message not logged")
	}

	if !strings.Contains(content, "Global warn message") {
		t.Error("Global warn message not logged")
	}

	if !strings.Contains(content, "Global error message") {
		t.Error("Global error message not logged")
	}
}

func TestLoggerLevelSetting(t *testing.T) {
	var buf bytes.Buffer

	// Create a logger that writes to buffer
	config := &LoggerConfig{
		Level:     LogLevelError,
		Output:    &buf,
		Timestamp: false,
	}
	logger := NewLogger(config)

	// Initially, only error messages should be logged
	logger.Info("This info message should not appear")
	if buf.Len() > 0 {
		t.Error("Info message should not be logged when level is ERROR")
	}

	// Change level to debug
	logger.SetLevel(LogLevelDebug)

	// Now debug messages should be logged
	logger.Debug("This debug message should appear")
	if buf.Len() == 0 {
		t.Error("Debug message should be logged when level is DEBUG")
	}

	// Check current level
	if logger.GetLevel() != LogLevelDebug {
		t.Error("Expected level to be DEBUG")
	}
}
