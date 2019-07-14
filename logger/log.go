package logger

import (
	"io"
	"log"
)

// Logger represents generic logger
type Logger interface {
	Debug(msg string, vals ...interface{})
	Info(msg string, vals ...interface{})
	Warn(msg string, vals ...interface{})
	Error(msg string, vals ...interface{})
}

type myLogger struct {
	logger *log.Logger
}

type Config struct {
	Output io.Writer
}

func NewLogger(c Config) Logger {
	return &myLogger{
		logger: log.New(c.Output, "Movie app:", log.LstdFlags),
	}
}

// TODO Following methods need fixing
func (l *myLogger) Debug(msg string, vals ...interface{}) {
	for _, v := range vals {
		msg += "," + v.(string) + ","
	}
	l.logger.Print(msg)
}

func (l *myLogger) Info(msg string, vals ...interface{}) {
	for _, v := range vals {
		msg += "," + v.(string) + ","
	}
	l.logger.Print(msg)
}

func (l *myLogger) Warn(msg string, vals ...interface{}) {
	for _, v := range vals {
		msg += "," + v.(string) + ","
	}
	l.logger.Print(msg)
}

func (l *myLogger) Error(msg string, vals ...interface{}) {
	for _, v := range vals {
		msg += "," + v.(string) + ","
	}
	l.logger.Print(msg)
}
