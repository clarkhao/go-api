package logger

import (
	"log"
	"os"
)

type LogLevel int

const (
	DebugLevel LogLevel = iota
	InfoLevel
	ErrorLevel
)

const (
	colorBlack   = "\033[30m"
	colorRed     = "\033[31m"
	colorGreen   = "\033[32m"
	colorYellow  = "\033[33m"
	colorBlue    = "\033[34m"
	colorMagenta = "\033[35m"
	colorCyan    = "\033[36m"
	colorWhite   = "\033[37m"
	colorReset   = "\033[0m"
)

type Logger struct {
	logger   *log.Logger
	logLevel LogLevel
}

func NewLogger(logLevel LogLevel) *Logger {
	return &Logger{
		logger:   log.New(os.Stdout, "", log.LstdFlags),
		logLevel: logLevel,
	}
}

func (l *Logger) Debug(message any) {
	if l.logLevel <= DebugLevel {
		l.logger.Printf("%v[DEBUG]: %v%v", colorBlue, message, colorReset)
	}
}

func (l *Logger) Info(message any) {
	if l.logLevel <= InfoLevel {
		l.logger.Printf("%v[INFO]: %v%v", colorGreen, message, colorReset)
	}
}

func (l *Logger) Error(message any) {
	if l.logLevel <= ErrorLevel {
		l.logger.Printf("%v[ERROR]: %v%v", colorRed, message, colorReset)
	}
}

var Log = NewLogger(DebugLevel)
