package logging

import (
	"fmt"
	"log"
	"os"
)

type Logger struct {
	file *os.File
}

func Init(pathToLogfile string) *Logger {
	f, _ := os.OpenFile(pathToLogfile, os.O_RDWR|os.O_CREATE, 0666)
	log.SetOutput(f)
	return &Logger{file: f}
}

func (l *Logger) Close() {
	err := l.file.Close()
	if err != nil {
		return
	}
}

func LogAndSkipError(err error) {
	log.Println(fmt.Sprintf("Skipped file: %s reload started...", err.Error()))
}

func LogMessage(message string) {
	log.Println(message)
}
