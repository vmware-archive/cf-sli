package logger

import "log"

type Logger interface  {
	Printf(format string, v ...interface{})
}

type logger struct {
}

func NewLogger() logger {
	return logger {
	}
}

func (l logger) Printf(format string, v ...interface{}) {
	log.Printf(format, v)
}
