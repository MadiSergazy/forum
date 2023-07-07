package logger

import (
	"log"
	"os"
)

type Logger struct {
	Error *log.Logger
	Info  *log.Logger
}

func NewLogger() *Logger {
	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)
	return &Logger{Info: infoLog, Error: errorLog}
}
