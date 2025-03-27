// Package logger provides functions for creating slog.Logger
package logger

import (
	"io"
	"log/slog"
	"os"
	"strings"
)

// ReplacerFunc defines the function of replacing the key or value of a log message.
type ReplacerFunc func(groups []string, a slog.Attr) (slog.Attr, bool)

var defaultLogOutput io.Writer = os.Stdout

// New returns *slog.Logger
func New(r ...ReplacerFunc) *slog.Logger {
	logLevel := &slog.LevelVar{}
	logLevel.Set(LogLevel)

	opts := &slog.HandlerOptions{
		Level: logLevel,
		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
			for _, replace := range r {
				if v, ok := replace(groups, a); ok {
					return v
				}
			}
			return a
		},
	}

	var handler slog.Handler = slog.NewJSONHandler(defaultLogOutput, opts)
	if strings.ToLower(os.Getenv("ENV")) == "local" {
		handler = slog.NewTextHandler(defaultLogOutput, opts)
	}

	logger := slog.New(handler)

	slog.SetDefault(logger)

	return logger
}
