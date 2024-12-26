package FlowWatch

import "github.com/sirupsen/logrus"

// Level is an enumeration for the log levels to abstract it from the logging library.
type Level uint32

const (
	Debug Level = iota
	Info
	Warn
	Error
	Fatal
)

// String returns the string representation of the log level.
func (l Level) String() string {
	switch l {
	case Debug:
		return "Debug"
	case Info:
		return "Info"
	case Warn:
		return "Warn"
	case Error:
		return "Error"
	case Fatal:
		return "Fatal"
	}
	return "Unknown"
}

// getLogrusLevel translates the Level enumeration to the logrus log level.
func (l Level) getLogrusLevel() logrus.Level {
	switch l {
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
