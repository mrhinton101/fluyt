package logger

import (
	"context"
	"log/slog"
)

type LogEntry struct {
	Level     slog.Level
	Err       error
	Component string
	Action    string
	Msg       string
	Target    string
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
