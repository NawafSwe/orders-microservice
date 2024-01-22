package logger

import "log"

type Logger interface {
	Log(data any)
}

func NewLogger() Logger {
	return LoggingImpl{}
}

type LoggingImpl struct {
}

func (l LoggingImpl) Log(data any) {
	log.Printf("%v\n", data)
}
