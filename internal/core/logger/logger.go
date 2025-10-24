package logger

import (
	"log"
	"os"

	tea "github.com/charmbracelet/bubbletea"
)

type Options struct {
	logFile string
	prefix  string
}

type Logger interface {
	Println(v ...any)
	Close() error
}

type fileLogger struct {
	logFile *os.File
}

func NewLogger(opts Options) (Logger, error) {
	l, err := tea.LogToFile(opts.logFile, opts.prefix)
	return &fileLogger{logFile: l}, err
}

func NewDebugLogger() (Logger, error) {
	return NewLogger(Options{
		logFile: "debug.log",
		prefix:  "DEBUG",
	})
}

func NewEmptyLogger() Logger {
	return &emptyLogger{}
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
