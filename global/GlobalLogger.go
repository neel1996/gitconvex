package global

import (
	"github.com/TwinProduction/go-color"
	"log"
)

// Global logger module for displaying the logs to the terminal

type LoggerInterface interface {
	LogInfo()
	LogWarning()
	LogError()
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

func (logger *Logger) Log(message string, status string) {
	logger.Message = message
	switch status {
	case StatusInfo:
		logger.LogInfo()
		break
	case StatusWarning:
		logger.LogWarning()
		break
	case StatusError:
		logger.LogError()
		break
	default:
		logger.LogInfo()
	}
}

func (logger *Logger) LogInfo() {
	log.Printf("%v[INFO]: %v%v\n", color.Cyan, color.Reset, logger.Message)
}

func (logger *Logger) LogWarning() {
	log.Printf("%v[WARNING]: %v%v\n", color.Yellow, color.Reset, logger.Message)
}

func (logger *Logger) LogError() {
	log.Printf("%v[ERROR]: %v%v\n", color.Red, color.Reset, logger.Message)
}
