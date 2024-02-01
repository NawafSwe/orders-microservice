package logger

import (
	"github.com/sirupsen/logrus"
	"io"
	"log"
	"os"
)

type Logger interface {
	Debug(data map[string]any, args any)
	Info(data map[string]any, args any)
	Warn(data map[string]any, args any)
	Error(data map[string]any, args any)
}

// getConfiguredLogger
// function to config the logger
func getConfiguredLogger() *logrus.Logger {
	l := logrus.New()
	l.SetFormatter(&logrus.JSONFormatter{})
	f, err := os.OpenFile("app.log", os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		log.Fatalf("failed to create a log file, err: %v", err)
	}
	// writing to log file and to datadog which is automatically will take what is written to os.Stdout
	l.SetOutput(io.MultiWriter(os.Stdout, f))
	l.SetLevel(logrus.DebugLevel)
	return l
}
func NewLogger() Logger {
	return LoggingImpl{logger: getConfiguredLogger()}
}

type LoggingImpl struct {
	logger *logrus.Logger
}

func (l LoggingImpl) Info(data map[string]any, args any) {
	l.logger.WithFields(data).Info(args)
}
func (l LoggingImpl) Warn(data map[string]any, args any) {
	l.logger.WithFields(data).Warn(args)
}
func (l LoggingImpl) Debug(data map[string]any, args any) {
	l.logger.WithFields(data).Debug(args)
}
func (l LoggingImpl) Error(data map[string]any, args any) {
	l.logger.WithFields(data).Error(args)
}
