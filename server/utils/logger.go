package utils

import (
	"fmt"
	"log/slog"
	"os"
)

var (
	logger *slog.Logger
)

func NewLogger(filename string, level string) *slog.Logger {
	logfile, err := os.OpenFile(filename, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0666)
	if err != nil {
		panic("Couldn't open file: %s" + filename)
	}
	logger = slog.New(slog.NewJSONHandler(logfile, &slog.HandlerOptions{
		AddSource: true,
		Level:     getLogLevel(level),
	}))
	return logger
}

func GetLogger() *slog.Logger {
	if logger == nil {
		fmt.Println("logger not intialized. call new logger")
		os.Exit(0)
	}
	return logger
}

func getLogLevel(level string) slog.Leveler {
	switch level {
	case "INFO":
		return slog.LevelInfo
	case "DEBUG":
		return slog.LevelDebug
	case "ERROR":
		return slog.LevelError
	default:
		return slog.LevelInfo
	}
}
