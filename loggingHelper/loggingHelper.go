package loggingHelper

import (
	"context"
	"errors"
	"fmt"
	"github.com/sirupsen/logrus"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"runtime"
	"sync"
	"time"
)

var (
	logHelper *LogHelper
	once      sync.Once
)

// LogHelper is an abstraction for the Logger instance to enable simpler switching between logging libraries.
type LogHelper struct {
	Logger *logrus.Logger
}

// LogrusContextHook is a hook for logrus that adds the file and line number to the log entry.
type LogrusContextHook struct{}

// LogrusOtelHook is a hook for logrus that enables logging to OpenTelemetry.
type LogrusOtelHook struct{}

// Levels returns all log levels for which the LogrusContextHook should be activated (warning level and higher,
// because runtime.Caller is expensive and debug, because it should be disabled in production).
func (hook LogrusContextHook) Levels() []logrus.Level {
	return []logrus.Level{
		logrus.DebugLevel,
		logrus.WarnLevel,
		logrus.ErrorLevel,
		logrus.FatalLevel,
		logrus.PanicLevel,
	}
}

// Fire is called when the LogrusContextHook is activated (when a log entry is made).
func (hook LogrusContextHook) Fire(entry *logrus.Entry) error {
	// Retrieve the call stack
	_, file, line, ok := runtime.Caller(7) // The number of function calls to skip to get to the caller

	// Add the file and line number to the log entry
	if !ok {
		err := errors.New("unable to retrieve the caller information and thus the file and line number")
		GetLogHelper().Debug(entry.Context, err)

		return nil // The hook should not return an error to ensure that other hooks are also executed
	}

	entry.Data["file"] = file
	entry.Data["line"] = line

	return nil
}

// Levels returns all log levels for which the LogrusOtelHook should be activated (warning level and higher).
func (hook LogrusOtelHook) Levels() []logrus.Level {
	return []logrus.Level{
		logrus.WarnLevel,
		logrus.ErrorLevel,
		logrus.FatalLevel,
		logrus.PanicLevel,
	}
}

// Fire is called when the LogrusOtelHook is activated (when a log entry is made).
func (hook LogrusOtelHook) Fire(entry *logrus.Entry) error {

	// Helper function to check the type and set a default value
	getAttributeValue := func(key string, defaultValue string) attribute.KeyValue {
		if value, ok := entry.Data[key]; ok {
			switch v := value.(type) {
			case string:
				return attribute.String(key, v)
			case int:
				return attribute.String(key, fmt.Sprintf("%d", v)) // Convert int to string
			}
		}
		return attribute.String(key, defaultValue)
	}

	// Create attributes
	messageValue := attribute.String("msg", entry.Message)
	levelValue := attribute.String("level", entry.Level.String())
	fileValue := getAttributeValue("file", "unknown")
	lineValue := getAttributeValue("line", "unknown")
	timeValue := attribute.String("time", entry.Time.Format(time.RFC3339))

	addEvent(entry.Context, messageValue, levelValue, fileValue, lineValue, timeValue)

	return nil
}

// initLogHelper initializes the LogHelper instance.
func initLogHelper() {
	// Create a new logrus logger with a JSON formatter
	logrusLogger := logrus.New()
	logrusLogger.SetLevel(logrus.InfoLevel) // Set the default log level to info for production environments
	logrusLogger.SetFormatter(&logrus.JSONFormatter{
		TimestampFormat: time.RFC3339,
	})

	logrusLogger.AddHook(LogrusContextHook{}) // Add the LogrusContextHook to add the file and line number to the log entry
	logrusLogger.AddHook(LogrusOtelHook{})    // Add the LogrusOtelHook to enable logging to OpenTelemetry

	logHelper = &LogHelper{
		Logger: logrusLogger,
	}
}

// GetLogHelper returns the LogHelper instance or creates a new one if it does not exist according to the singleton pattern.
func GetLogHelper() *LogHelper {
	// Create a new LogHelper instance if it does not exist
	once.Do(func() {
		initLogHelper()
	})

	return logHelper
}

// addEvent adds an event to the trace span.
func addEvent(ctx context.Context, args ...attribute.KeyValue) {
	span := trace.SpanFromContext(ctx)
	if span != nil {
		// Add the event to the span
		span.AddEvent("log", trace.WithAttributes(args...))
		// TODO: Use otel log exporter to export logs even if there is no surrounding span
	}
}

// Abstraction for log functions to enable simpler switching between logging libraries.
// Context is required to add the event to the span (if possible).

// Debug logs a message at the debug level.
func (lh *LogHelper) Debug(ctx context.Context, args ...interface{}) {
	lh.Logger.WithContext(ctx).Debug(args...)
}

// Info logs a message at the info level.
func (lh *LogHelper) Info(ctx context.Context, args ...interface{}) {
	lh.Logger.WithContext(ctx).Info(args...)
}

// Warn logs a message at the warning level.
func (lh *LogHelper) Warn(ctx context.Context, args ...interface{}) {
	lh.Logger.WithContext(ctx).Warn(args...)
}

// Error logs a message at the error level.
func (lh *LogHelper) Error(ctx context.Context, args ...interface{}) {
	lh.Logger.WithContext(ctx).Error(args...)
}

// Fatal logs a message at the fatal level.
func (lh *LogHelper) Fatal(ctx context.Context, args ...interface{}) {
	lh.Logger.WithContext(ctx).Fatal(args...)
}

// Level is an enumeration for the log levels to abstract it from the logging library.
type Level uint32

const (
	Debug Level = iota
	Info
	Warn
	Error
	Fatal
)

// getLogrusLevel translates the Level enumeration to the logrus log level.
func (l *Level) getLogrusLevel() logrus.Level {
	switch *l {
	case Debug:
		return logrus.DebugLevel
	case Info:
		return logrus.InfoLevel
	case Warn:
		return logrus.WarnLevel
	case Error:
		return logrus.ErrorLevel
	case Fatal:
		return logrus.FatalLevel
	default:
		return logrus.DebugLevel
	}
}

// SetLogLevel updates the log level of the logger library.
func SetLogLevel(level Level) {
	GetLogHelper().Logger.SetLevel(level.getLogrusLevel())
}
