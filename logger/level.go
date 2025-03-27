package logger

import (
	"log/slog"
	"os"
	"strings"
)

var logLevel = map[string]slog.Level{
	"DEBUG": slog.LevelDebug,
	"INFO":  slog.LevelInfo,
	"WARN":  slog.LevelWarn,
	"ERROR": slog.LevelError,
}

var defaultLogLevel = "ERROR"
var LogLevel slog.Level

func init() {
	logLevelConf := strings.ToUpper(os.Getenv("LOG_LEVEL"))
	if logLevelConf == "" {
		logLevelConf = defaultLogLevel
	}

	LogLevel = logLevel[logLevelConf]
}
