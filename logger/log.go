package logger

import (
	"log"
	"os"
)

func ErrLogger() *log.Logger {
	return log.New(os.Stderr, "[ERR]", log.Default().Flags())
}

func InfoLogger() *log.Logger {
	return log.New(os.Stderr, "[INFO]", log.Default().Flags())
}
