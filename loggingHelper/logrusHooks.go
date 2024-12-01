package loggingHelper

import (
	"context"
	"fmt"
	"github.com/LucaSchmitz2003/FlowWatch/otelHelper"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"runtime"
	"time"
)

// LogrusContextHook is a hook for logrus that adds the file and line number to the log entry.
type LogrusContextHook struct{}

// LogrusOtelHook is a hook for logrus that enables logging to OpenTelemetry.
type LogrusOtelHook struct{}

// LogrusOtelShutdownHook is a hook for logrus that ensures that the connection to OpenTelemetry is shut down properly.
type LogrusOtelShutdownHook struct{}

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

// addEvent adds an event to the trace span.
func addEvent(ctx context.Context, args ...attribute.KeyValue) {
	span := trace.SpanFromContext(ctx)
	if span != nil {
		// Add the event to the span
		span.AddEvent("log", trace.WithAttributes(args...))
		// TODO: Use otel log exporter to export logs even if there is no surrounding span
	}
}

// Levels returns all log levels for which the LogrusOtelShutdownHook should be activated
// (fatal level and higher since it terminates the program).
func (hook LogrusOtelShutdownHook) Levels() []logrus.Level {
	return []logrus.Level{
		logrus.FatalLevel,
		logrus.PanicLevel,
	}
}

// Fire is called when the LogrusOtelShutdownHook is activated (when a fatal log entry is made).
func (hook LogrusOtelShutdownHook) Fire(entry *logrus.Entry) error {
	otelHelper.Shutdown() // Shutdown the OpenTelemetry connection
	return nil
}
