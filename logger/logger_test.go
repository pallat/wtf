package logger

import (
	"bytes"
	"encoding/json"
	"log/slog"
	"os"
	"testing"
)

func TestNewLogger(t *testing.T) {
	t.Run("ENV not set logger handler should be JSON Handler", func(t *testing.T) {
		os.Unsetenv("ENV")
		buf := bytes.NewBuffer([]byte{})

		defaultLogOutput = buf
		defer func() { defaultLogOutput = os.Stdout }()

		l := New()

		_, ok := l.Handler().(*slog.JSONHandler)
		if !ok {
			t.Errorf("ENV local expect handler type is *slog.JSONHandler but actual is %T", l.Handler())
		}
	})
	t.Run("ENV local logger handler should be Text Handler", func(t *testing.T) {
		os.Setenv("ENV", "local")
		buf := bytes.NewBuffer([]byte{})

		defaultLogOutput = buf
		defer func() { defaultLogOutput = os.Stdout }()

		l := New()

		_, ok := l.Handler().(*slog.TextHandler)
		if !ok {
			t.Errorf("ENV local expect handler type is *slog.TextHandler but actual is %T", l.Handler())
		}
	})
	t.Run("replacer", func(t *testing.T) {
		os.Unsetenv("ENV")
		buf := bytes.NewBuffer([]byte{})

		defaultLogOutput = buf
		defer func() { defaultLogOutput = os.Stdout }()

		replace := func(groups []string, a slog.Attr) (slog.Attr, bool) {
			if a.Key == "arise-testing" {
				return slog.String("infinitas-testing", a.Value.String()), true
			}
			return a, false
		}

		l := New(replace)

		l.Info("message", slog.String("arise-testing", "testing-message"))

		var m map[string]string
		json.NewDecoder(buf).Decode(&m)
		v, ok := m["infinitas-testing"]
		if !ok {
			t.Errorf("replacer replace key should relace %q to %q: actual %v\n", "arise-testing", "infinitas-testing", m)
		}

		if v != "testing-message" {
			t.Errorf("replacer replace key should relace %q to %q: actual %v\n", "arise-testing", "infinitas-testing", m)
		}
	})
}
