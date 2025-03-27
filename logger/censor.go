package logger

import "log/slog"

var censors = map[string]string{}

// CensorReplacer replace or mask log value in case of sensitive information
func CensorReplacer(groups []string, a slog.Attr) (slog.Attr, bool) {
	for k, v := range censors {
		if a.Key == k {
			return slog.String(k, v), true
		}
	}
	return a, false
}
