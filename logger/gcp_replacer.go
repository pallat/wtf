package logger

import (
	"log/slog"
)

var gcpKeys = map[string]string{
	"level": "severity",
	"msg":   "message",
	"time":  "timestamp",
}

// GCPKeyReplacer replaces log keys for GCP standard
func GCPKeyReplacer(groups []string, a slog.Attr) (slog.Attr, bool) {
	for k, v := range gcpKeys {
		if a.Key == k {
			return slog.String(v, a.Value.String()), true
		}
	}
	return a, false
}
