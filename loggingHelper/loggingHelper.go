package loggingHelper

import (
	"github.com/sirupsen/logrus"
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

// initLogHelper initializes the LogHelper instance.
func initLogHelper() {
	// Create a new logrus logger with a JSON formatter
	logrusLogger := logrus.New()
	logrusLogger.SetLevel(logrus.InfoLevel) // Set the default log level to info for production environments
	logrusLogger.SetFormatter(&logrus.JSONFormatter{
		TimestampFormat: time.RFC3339,
	})

	logrusLogger.AddHook(LogrusContextHook{})      // Add the LogrusContextHook to add the file and line number to the log entry
	logrusLogger.AddHook(LogrusOtelHook{})         // Add the LogrusOtelHook to enable logging to OpenTelemetry
	logrusLogger.AddHook(LogrusOtelShutdownHook{}) // Add the LogrusOtelShutdownHook to ensure that the connection is shut down properly

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
