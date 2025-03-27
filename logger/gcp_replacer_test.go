package logger

import (
	"bytes"
	"encoding/json"
	"log/slog"
	"os"
	"testing"
)

func TestGCPKeyReplacer(t *testing.T) {
	t.Run("GCP key replacer", func(t *testing.T) {
		os.Unsetenv("ENV")
		buf := bytes.NewBuffer([]byte{})

		defaultLogOutput = buf
		defer func() { defaultLogOutput = os.Stdout }()

		l := New(GCPKeyReplacer)

		l.Info("message")

		var m map[string]any
		json.NewDecoder(buf).Decode(&m)
		_, ok := m["severity"]
		if !ok {
			t.Errorf("replacer replace key should relace key %q to %q: actual %v\n", "level", "severity", m)
		}

		_, ok = m["message"]
		if !ok {
			t.Errorf("replacer replace key should relace key %q to %q: actual %v\n", "msg", "message", m)
		}

		_, ok = m["timestamp"]
		if !ok {
			t.Errorf("replacer replace key should relace key %q to %q: actual %v\n", "time", "timestamp", m)
		}
	})
	t.Run("not found any key", func(t *testing.T) {
		_, ok := GCPKeyReplacer([]string{}, slog.Attr{})
		if ok {
			t.Errorf("not any matched key expect false")
		}
	})
}
