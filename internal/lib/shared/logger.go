package shared

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"
)

// LogLevel represents the severity of a log entry
type LogLevel string

const (
	LogLevelDebug LogLevel = "debug"
	LogLevelInfo  LogLevel = "info"
	LogLevelWarn  LogLevel = "warn"
	LogLevelError LogLevel = "error"
)

var logLevelPriority = map[LogLevel]int{
	LogLevelDebug: 0,
	LogLevelInfo:  1,
	LogLevelWarn:  2,
	LogLevelError: 3,
}

// LogEntry represents a structured log entry
type LogEntry struct {
	Level     LogLevel               `json:"level"`
	Message   string                 `json:"message"`
	Timestamp string                 `json:"timestamp"`
	Context   map[string]interface{} `json:"context,omitempty"`
}

// Logger provides structured logging with JSON (production) and pretty (development) output.
type Logger struct {
	minLevel    LogLevel
	isJSON      bool
	baseContext map[string]interface{}
}

// NewLogger creates a new Logger.
// In production (GIN_MODE=release), outputs JSON. Otherwise, outputs colorized pretty-print.
// Log level can be set via LOG_LEVEL env var; defaults to "info" in production, "debug" in development.
func NewLogger(isDevelopment bool) *Logger {
	minLevel := getMinLogLevel(isDevelopment)
	return &Logger{
		minLevel:    minLevel,
		isJSON:      !isDevelopment,
		baseContext: make(map[string]interface{}),
	}
}

// Child creates a child logger that inherits and extends the parent's base context.
func (l *Logger) Child(context map[string]interface{}) *Logger {
	merged := make(map[string]interface{}, len(l.baseContext)+len(context))
	for k, v := range l.baseContext {
		merged[k] = v
	}
	for k, v := range context {
		merged[k] = v
	}
	return &Logger{
		minLevel:    l.minLevel,
		isJSON:      l.isJSON,
		baseContext: merged,
	}
}

func getMinLogLevel(isDevelopment bool) LogLevel {
	envLevel := strings.ToLower(os.Getenv("LOG_LEVEL"))
	if _, ok := logLevelPriority[LogLevel(envLevel)]; ok {
		return LogLevel(envLevel)
	}
	if isDevelopment {
		return LogLevelDebug
	}
	return LogLevelInfo
}

func (l *Logger) shouldLog(level LogLevel) bool {
	return logLevelPriority[level] >= logLevelPriority[l.minLevel]
}

func (l *Logger) log(level LogLevel, message string, context map[string]interface{}) {
	if !l.shouldLog(level) {
		return
	}

	merged := make(map[string]interface{}, len(l.baseContext)+len(context))
	for k, v := range l.baseContext {
		merged[k] = v
	}
	for k, v := range context {
		merged[k] = v
	}

	entry := LogEntry{
		Level:     level,
		Message:   message,
		Timestamp: time.Now().Format(time.RFC3339),
		Context:   merged,
	}

	if l.isJSON {
		l.formatJSON(entry)
	} else {
		l.formatPretty(entry)
	}
}

func (l *Logger) formatJSON(entry LogEntry) {
	data, err := json.Marshal(entry)
	if err != nil {
		fmt.Fprintf(os.Stderr, "logger marshal error: %v\n", err)
		return
	}
	fmt.Fprintln(os.Stdout, string(data))
}

func (l *Logger) formatPretty(entry LogEntry) {
	colors := map[LogLevel]string{
		LogLevelDebug: "\033[36m", // cyan
		LogLevelInfo:  "\033[32m", // green
		LogLevelWarn:  "\033[33m", // yellow
		LogLevelError: "\033[31m", // red
	}
	reset := "\033[0m"
	dim := "\033[2m"

	color := colors[entry.Level]
	levelStr := fmt.Sprintf("%-5s", strings.ToUpper(string(entry.Level)))

	output := fmt.Sprintf("%s%s%s %s%s%s %s", dim, entry.Timestamp, reset, color, levelStr, reset, entry.Message)

	if len(entry.Context) > 0 {
		ctx, _ := json.Marshal(entry.Context)
		output += fmt.Sprintf(" %s%s%s", dim, string(ctx), reset)
	}

	fmt.Fprintln(os.Stdout, output)
}

// Debug logs a debug message
func (l *Logger) Debug(message string, context map[string]interface{}) {
	l.log(LogLevelDebug, message, context)
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
