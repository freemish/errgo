package errgo

import (
	"bytes"
	"fmt"
	"runtime"
	"strings"
)

// A StackFrame contains all necessary information about to generate a line
// in a callstack.
type StackFrame struct {
	Caller       uintptr
	File         string
	LineNumber   int
	FunctionName string
	Package      string
}

// NewStackFrame populates a stack frame object from the program counter.
func NewStackFrame(caller uintptr) (frame StackFrame) {
	frame = StackFrame{Caller: caller}
	if frame.Func() == nil {
		return
	}
	frame.Package, frame.FunctionName = packageAndName(frame.Func())
	frame.File, frame.LineNumber = frame.Func().FileLine(caller - 1)
	return
}

// Func returns the function that contained this frame.
func (frame *StackFrame) Func() *runtime.Func {
	if frame.Caller == 0 {
		return nil
	}
	return runtime.FuncForPC(frame.Caller)
}

// String returns the stackframe formatted in the same way as go does
// in runtime/debug.Stack()
func (frame *StackFrame) String() string {
	return fmt.Sprintf("%s: %s: line %d", RelativeFilePath(frame.File), frame.FunctionName, frame.LineNumber)
}

func packageAndName(fn *runtime.Func) (string, string) {
	name := fn.Name()
	pkg := ""

	// The name includes the path name to the package, which is unnecessary
	// since the file name is already included.  Plus, it has center dots.
	// That is, we see
	//  runtime/debug.*T·ptrmethod
	// and want
	//  *T.ptrmethod
	// Since the package path might contains dots (e.g. code.google.com/...),
	// we first remove the path prefix if there is one.
	if lastslash := strings.LastIndex(name, "/"); lastslash >= 0 {
		pkg += name[:lastslash] + "/"
		name = name[lastslash+1:]
	}
	if period := strings.Index(name, "."); period >= 0 {
		pkg += name[:period]
		name = name[period+1:]
	}

	name = strings.Replace(name, "·", ".", -1)
	return pkg, name
}

// RelativeFilePath removes absolute paths - basically cuts
// out addresses at /src/ and before.
func RelativeFilePath(file string) string {
	folders := strings.Split(file, "/")
	count := 0
	for _, folder := range folders {
		if folder == "src" {
			break
		}
		count++
	}
	if count >= len(folders) {
		return file
	}
	buf := bytes.Buffer{}
	for _, folder := range folders[count+1:] {
		buf.WriteString("/")
		buf.WriteString(folder)
	}
	return buf.String()
}
