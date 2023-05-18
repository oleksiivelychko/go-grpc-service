package logger

import (
	"fmt"
	"log"
	"os"
)

type Logger struct {
	info  *log.Logger
	error *log.Logger
}

func New() *Logger {
	return &Logger{
		log.New(os.Stdout, "[INFO]\t", log.Ldate|log.Ltime),
		log.New(os.Stderr, "[ERROR]\t", log.Ldate|log.Ltime|log.Lshortfile),
	}
}

func (logger *Logger) Info(message string, args ...any) {
	logger.info.Printf("%s\n", fmt.Sprintf(message, args...))
}

func (logger *Logger) Error(message string, args ...any) {
	logger.error.Printf("%s\n", fmt.Sprintf(message, args...))
}
