package global

import (
	"github.com/sirupsen/logrus"
)

// Global logger module for displaying the logs to the terminal

type LoggerInterface interface {
	info()
	warn()
	error()
}

type Logger struct {
	Message string
}

const (
	StatusError   = "StatusError"
	StatusInfo    = "StatusInfo"
	StatusWarning = "StatusWarning"
)

// Log gets the log message and the status selector
// The log flag is chosen and color coded based on the status received
func (l *Logger) Log(message string, status string) {
	l.Message = message
	switch status {
	case StatusInfo:
		l.info()
		break
	case StatusWarning:
		l.warn()
		break
	case StatusError:
		l.error()
		break
	default:
		l.info()
	}
}

func (l *Logger) info() {
	logger := getLogger()
	logger.Info(l.Message)
}

func (l *Logger) warn() {
	logger := getLogger()
	logger.Warn(l.Message)
}

func (l *Logger) error() {
	logger := getLogger()
	logger.Error(l.Message)
}

func getLogger() *logrus.Logger {
	logger := logrus.New()
	formatter := new(logrus.TextFormatter)
	formatter.FullTimestamp = true
	formatter.TimestampFormat = "2006-01-02 15:04:05"
	formatter.PadLevelText = true
	logger.SetFormatter(formatter)
	return logger
}
