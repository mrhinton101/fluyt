package logger

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"os"
)

type LogEntry struct {
	Level     slog.Level // Log level
	Err       error      // error object
	Component string     //
	Action    string
	Msg       string
	Target    string
}

// Initialize the Logger and open stream to write to file
func InitLogger(logFileName string) (file *os.File) {

	// Open the file in append mode, create it if it doesn't exist, and allow writing.
	// 0644 represents the file permissions (read/write for owner, read-only for group and others).
	file, err := os.OpenFile(logFileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println("file open fail")
	}
	multiWriter := io.MultiWriter(os.Stdout, file)
	// Use a TextHandler or JSONHandler depending on what you prefer.
	handler := slog.NewJSONHandler(multiWriter, &slog.HandlerOptions{
		Level: slog.LevelError, // Adjust this as needed
	})

	// Set as default logger
	slog.SetDefault(slog.New(handler))

	return file
}

// SLogger provides consistent structured logging using slog with context fields.
func SLogger(entry LogEntry) {
	logArgs := []any{
		"err", entry.Err,
		"component", entry.Component,
		"action", entry.Action,
		"target", entry.Target,
	}

	slog.Log(context.Background(), entry.Level, entry.Msg, logArgs...)
}
