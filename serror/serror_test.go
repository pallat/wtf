package serror

import (
	"errors"
	"fmt"
	"path/filepath"
	"runtime"
	"testing"
)

func TestWrap(t *testing.T) {
	given := "testing"
	var pc uintptr
	var file string
	var line int

	var first = func() error {
		pc, file, line, _ = runtime.Caller(0)
		return Wrap(errors.New(given))
	}

	err := first()

	expect := fmt.Sprintf("((testing+%s:%d:%s))", filepath.Base(file), line+1, filepath.Base(runtime.FuncForPC(pc).Name()))

	if err.Error() != expect {
		t.Errorf("expect:\n%q\n\nactual:\n%q\n", expect, err.Error())
	}
}

func TestWrapSkipZeroStep_ErrorOccurredFromCaller(t *testing.T) {
	given := "testing"
	var pc uintptr
	var file string
	var line int

	var first = func() error {
		pc, file, line, _ = runtime.Caller(0)
		return WrapSkip(errors.New(given), 0)
	}

	err := first()

	expect := fmt.Sprintf("((testing+%s:%d:%s))", filepath.Base(file), line+1, filepath.Base(runtime.FuncForPC(pc).Name()))

	if err.Error() != expect {
		t.Errorf("expect:\n%q\n\nactual:\n%q\n", expect, err.Error())
	}
}

func TestWrapSkipOneStepBack_ErrorOccurredFromCallerofCaller(t *testing.T) {
	given := "testing"
	var pc uintptr
	var file string
	var line int

	var second = func() error {
		return WrapSkip(errors.New(given), 1)
	}

	var first = func() error {
		pc, file, line, _ = runtime.Caller(0)
		return second()
	}

	err := first()

	expect := fmt.Sprintf("((testing+%s:%d:%s))", filepath.Base(file), line+1, filepath.Base(runtime.FuncForPC(pc).Name()))

	if err.Error() != expect {
		t.Errorf("expect:\n%q\n\nactual:\n%q\n", expect, err.Error())
	}
}

func TestCallerNotRecoverSkip(t *testing.T) {
	const notRecover = 4
	s := caller(notRecover)
	if "" != s {
		t.Error("the origin not recover skip is 4")
	}
}

func TestDecodeMessageEmptyString(t *testing.T) {
	msg, atts := DecodeMessage("")

	if "" != msg {
		t.Errorf("given empty string to Decode expect empty string msg but actual %q\n", msg)
	}

	if 0 != len(atts) {
		t.Errorf("given empty string to DecodeMessage expect 0 lenght Attrs but actual %d\n", len(atts))
	}
}

func TestSErrorDecode(t *testing.T) {
	given := "test message"
	err := New(given)
	msg, attrs := DecodeMessage(err.Error())

	if msg != "test message" {
		t.Errorf("%q message is expected but actual %q", given, msg)
	}

	if len(attrs) != 3 {
		t.Errorf("attrs should have 3 elements but actual %d", len(attrs))
	}
}

func TestSErrorWithMessageDecode(t *testing.T) {
	given := "test message"
	err := New(given)

	message := fmt.Sprintf("first message: %s", err)

	msg, attrs := DecodeMessage(message)

	if msg != "first message: test message" {
		t.Errorf("message: %q\n", message)
		t.Errorf("%q message is expected but actual %q\n", "first message: "+given, msg)
	}

	if len(attrs) != 3 {
		t.Errorf("attrs should have 3 elements but actual %d", len(attrs))
	}
}

func TestDecodeMessagePlainError(t *testing.T) {
	given := "test message"
	err := errors.New(given)
	msg, attrs := DecodeMessage(err.Error())

	if msg != "test message" {
		t.Errorf("%q message is expected but actual %q", given, msg)
	}

	if len(attrs) != 0 {
		t.Errorf("attrs should have 3 elements but actual %d", len(attrs))
	}
}

func TestDecodeMessageWrongPattern(t *testing.T) {
	givenMsg := prefix + "message" + sufix
	msg, atts := DecodeMessage(givenMsg)

	if givenMsg != msg {
		t.Errorf("expect %q string msg but actual %q\n", givenMsg, msg)
	}

	if 0 != len(atts) {
		t.Errorf("given wrong pattern to DecodeMessage expect 0 lenght Attrs but actual %d\n", len(atts))
	}
}

func TestDecodeMessageWrongSourcePattern(t *testing.T) {
	givenMsg := prefix + "message" + separator + "filename.go" + sufix
	msg, atts := DecodeMessage(givenMsg)

	if givenMsg != msg {
		t.Errorf("expect %q string msg but actual %q\n", givenMsg, msg)
	}

	if 0 != len(atts) {
		t.Errorf("given wrong pattern to DecodeMessage expect 0 lenght Attrs but actual %d\n", len(atts))
	}
}
