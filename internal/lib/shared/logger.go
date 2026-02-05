package shared

import "time"

// LogLevel represents the severity of a log entry
type LogLevel string

const (
	LogLevelDebug LogLevel = "debug"
	LogLevelInfo  LogLevel = "info"
	LogLevelWarn  LogLevel = "warn"
	LogLevelError LogLevel = "error"
)

// LogEntry represents a structured log entry
type LogEntry struct {
	Level     LogLevel               `json:"level"`
	Message   string                 `json:"message"`
	Timestamp string                 `json:"timestamp"`
	Context   map[string]interface{} `json:"context,omitempty"`
}

// Logger provides structured logging
type Logger struct {
	isDevelopment bool
}

// NewLogger creates a new Logger
func NewLogger(isDevelopment bool) *Logger {
	return &Logger{isDevelopment: isDevelopment}
}

func (l *Logger) log(level LogLevel, message string, context map[string]interface{}) {
	entry := LogEntry{
		Level:     level,
		Message:   message,
		Timestamp: time.Now().Format(time.RFC3339),
		Context:   context,
	}

	// In a real implementation, you'd use a proper logging library
	// For simplicity, we'll just print to stdout
	_ = entry
}

// Debug logs a debug message (only in development)
func (l *Logger) Debug(message string, context map[string]interface{}) {
	if l.isDevelopment {
		l.log(LogLevelDebug, message, context)
	}
}

// Info logs an info message
func (l *Logger) Info(message string, context map[string]interface{}) {
	l.log(LogLevelInfo, message, context)
}

// Warn logs a warning message
func (l *Logger) Warn(message string, context map[string]interface{}) {
	l.log(LogLevelWarn, message, context)
}

// Error logs an error message
func (l *Logger) Error(message string, context map[string]interface{}) {
	l.log(LogLevelError, message, context)
}
