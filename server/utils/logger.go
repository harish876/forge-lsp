package utils

import (
	"log"
	"os"
)

func GetLogger(filename string) *log.Logger {
	logfile, err := os.OpenFile(filename, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0666)
	if err != nil {
		panic("Couldn't open file: %s" + filename)
	}

	return log.New(logfile, "[forge-lsp] ", log.Ldate|log.Ltime|log.Lshortfile)
}
