package logger

import (
	"encoding/json"
	"fmt"
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

// Config is used to providing logging config
type Config struct {
	Output io.Writer
	// Could have format as well as a config like json
}

// NewLogger return a Logger
func NewLogger(c Config) Logger {
	return &myLogger{
		logger: log.New(c.Output, "", log.LstdFlags),
	}
}

func (l *myLogger) Debug(msg string, vals ...interface{}) {
	values := make(map[string]string, 0)
	values["msg"] = msg
	values["loglevel"] = "DEBUG"
	values["values"] = fmt.Sprintf("%+v", vals)
	b, _ := json.Marshal(values)
	l.logger.Print(string(b))

}

func (l *myLogger) Info(msg string, vals ...interface{}) {
	values := make(map[string]string, 0)
	values["msg"] = msg
	values["loglevel"] = "INFO"
	values["values"] = fmt.Sprintf("%+v", vals)
	b, _ := json.Marshal(values)
	l.logger.Print(string(b))

}

func (l *myLogger) Warn(msg string, vals ...interface{}) {
	values := make(map[string]string, 0)
	values["msg"] = msg
	values["loglevel"] = "WARN"
	values["values"] = fmt.Sprintf("%+v", vals)
	b, _ := json.Marshal(values)
	l.logger.Print(string(b))

}

func (l *myLogger) Error(msg string, vals ...interface{}) {
	values := make(map[string]string, 0)
	values["msg"] = msg
	values["loglevel"] = "ERROR"
	values["values"] = fmt.Sprintf("%+v", vals)
	b, _ := json.Marshal(values)
	l.logger.Print(string(b))

}
