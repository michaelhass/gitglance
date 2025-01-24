package logger

import (
	"log"
	"os"

	tea "github.com/charmbracelet/bubbletea"
)

type Logger interface {
	Println(v ...any)
	Close() error
}

type fileLogger struct {
	logFile *os.File
}

func NewLogger(isDebug bool) (Logger, error) {
	if isDebug {
		l, err := tea.LogToFile("debug.log", "DEBUG")
		return &fileLogger{logFile: l}, err
	}
	return &emptyLogger{}, nil
}

func (l *fileLogger) Close() error {
	return l.logFile.Close()
}

func (l *fileLogger) Println(v ...any) {
	log.Println(v...)
}

type emptyLogger struct{}

func (l *emptyLogger) Close() error {
	return nil
}

func (l *emptyLogger) Println(v ...any) {
}
