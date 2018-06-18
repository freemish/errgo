package errgo

import (
	"bytes"
	"fmt"
	"runtime"
)

// The maximum number of stackframes on any error.
var MaxStackDepth = 50

// Error is an error with an attached stacktrace. It can be used
// wherever the builtin error interface is expected.
type StackableError struct {
	Err      error
	Prefixes []string
	stack    []uintptr
	frames   []StackFrame
}

// Error returns the prefixed error message.
func (err *StackableError) Error() string {
	msg := err.Err.Error()

	for _, prefix := range err.Prefixes {
		msg = fmt.Sprintf("%s: %s", prefix, msg)
	}

	return msg
}

// Callers allows access to program counters.
func (err *StackableError) Callers() []uintptr {
	return err.stack
}

// Wrap makes a StackableError from an interface;
// returns itself if the interface is already a *StackableError.
func Wrap(e interface{}) *StackableError {
	var err error

	switch e := e.(type) {
	case *StackableError:
		return e // this adds a caller to the stack!
	case error:
		err = e
	default:
		err = fmt.Errorf("%v", e)
	}

	return newStackableError(err, 1)
}

// WrapPrefix makes a StackableError from the given value. If that value is already an
// error then it will be used directly, if not, it will be passed to
// fmt.Errorf("%v"). The prefix parameter is used to add a prefix to the
// error message when calling Error().
func WrapPrefix(e interface{}, prefix string) *StackableError {
	err := Wrap(e)
	err.Prefixes = append(err.Prefixes, prefix)
	return err
}

func newStackableError(e error, skip int) *StackableError {
	var prefixes []string
	stack := make([]uintptr, MaxStackDepth)
	length := runtime.Callers(2+skip, stack)
	return &StackableError{
		Err:      e,
		stack:    stack[:length],
		Prefixes: prefixes,
	}
}

// Is detects whether the error is equal to a given error. Errors
// are considered equal by this function if they are the same object,
// or if they both contain the same error inside an errors.Error.
func Is(e error, original error) bool {
	if e == original {
		return true
	}

	if e, ok := e.(*StackableError); ok {
		return Is(e.Err, original)
	}

	if original, ok := original.(*StackableError); ok {
		return Is(e, original.Err)
	}

	return false
}

// StackFrames returns an array of frames containing information about the
// stack.
func (err *StackableError) StackFrames() []StackFrame {
	if err.frames == nil {
		err.frames = make([]StackFrame, len(err.stack))
		for i, pc := range err.stack {
			err.frames[i] = NewStackFrame(pc)
		}
	}

	return err.frames
}

// Stack returns the callstack formatted the same way that go does
// in runtime/debug.Stack()
func (err *StackableError) Stack() string {
	buf := bytes.Buffer{}

	for _, frame := range err.StackFrames() {
		buf.WriteString(frame.String())
		buf.WriteString("\n")
	}

	return buf.String()
}

// StackTrace prints a stacktrace like:
// ERROR: (prefixed message)
// (stack returned by Stack())
func (err *StackableError) StackTrace() string {
	return "ERROR: " + err.Error() + "\n" + err.Stack()
}
