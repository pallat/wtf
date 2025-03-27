// Package error provides functions for handling error
package serror

import (
	"errors"
	"fmt"
	"log/slog"
	"path/filepath"
	"runtime"
	"strings"
)

// SError represents error structure with
// line of code that program had been executing when an error occurred
type SError struct {
	err error
	at  string
}

const (
	prefix          = "(("
	sufix           = "))"
	separator       = "+"
	sourceSeparator = ":"
	prefixSize      = len(prefix)
)

// Error returns error message
func (err *SError) Error() string {
	return fmt.Sprintf("((%s+%s))", err.err, err.at)
}

// New returns [SError]
func New(s string) *SError {
	return &SError{
		err: errors.New(s),
		at:  caller(2),
	}
}

// WrapSkip wraps the previous error and includes the line of code the occurred error
// with skip caller given step
//
// Example: if function A caught an error
//
//	func A() {
//	  err := occurred()
//	  err = serror.WrapSkip(err, 0)
//	}
//
// Example: if function A call function B and B caught an error
// but you need to keep line of code in function A
//
//	func A() {
//	  err := B()
//	}
//
//	func B() {
//	  err := occurred()
//	  err = serror.WrapSkip(err, 1)
//	}
func WrapSkip(err error, skip int) *SError {
	skip += 2
	if skip < 2 {
		skip = 2
	}
	return &SError{
		err: err,
		at:  caller(skip),
	}
}

func Wrap(err error) *SError {
	return &SError{
		err: err,
		at:  caller(2),
	}
}

func caller(skip int) string {
	pc, file, no, ok := runtime.Caller(skip)
	if ok {
		b := filepath.Base(file)
		f := filepath.Base(runtime.FuncForPC(pc).Name())
		return fmt.Sprintf("%s:%d:%s", b, no, f)
	}
	return ""
}

// DecodeMessage decodes error message that was generated from serror (New,Wrap or WrapSkip)
// to plain message and []slog.Attr
//
// Example: DecodeMessage("!!original error message:@handler.go:23:handler.Serve")
//
//	 message = "original error message"
//		slog.Attr[0] = slog.String("file", "handler.go")
//		slog.Attr[1] = slog.String("line", "23")
//		slog.Attr[2] = slog.String("func", "handler.Serve")
func DecodeMessage(s string) (msg string, attrs []slog.Attr) {
	if s == "" {
		return "", []slog.Attr{}
	}

	serrorFrom := strings.Index(s, prefix)
	serrorTo := strings.Index(s, sufix)

	if serrorFrom == serrorTo {
		return s, []slog.Attr{}
	}

	serrorMsg := s[serrorFrom+2 : serrorTo]

	elem := strings.Split(serrorMsg, separator)
	if 2 != len(elem) {
		return s, []slog.Attr{}
	}

	sources := strings.Split(elem[1], sourceSeparator)
	if 3 != len(sources) {
		return s, []slog.Attr{}
	}

	s = strings.Replace(s, s[serrorFrom:serrorTo+2], elem[0], 1)

	return s, []slog.Attr{
		slog.String("file", sources[0]),
		slog.String("line", sources[1]),
		slog.String("func", sources[2]),
	}
}
