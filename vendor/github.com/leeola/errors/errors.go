package errors

import (
	"errors"
	"fmt"
	"path"
	"runtime"
	"strings"
)

func Cause(err error) error {
	if err == nil {
		return nil
	}

	cErr, ok := err.(causer)
	if !ok {
		return err
	}

	if cause := cErr.Cause(); cause != nil {
		return cause
	}

	return err
}

func Errorf(f string, s ...interface{}) error {
	msg := fmt.Sprintf(f, s...)
	return &errWrap{
		err:      errors.New(msg),
		Msg:      msg,
		SumStack: []string{callerLine() + ": " + msg},
	}
}

// Equals checks if a and b are the same, have the same cause, or variations within.
func Equals(a error, b error) bool {
	if a == b {
		return true
	}

	if a == nil || b == nil {
		return false
	}

	aCause := Cause(a)
	if aCause == b {
		return true
	}

	bCause := Cause(b)
	if bCause == a {
		return true
	}

	return aCause == bCause
}

func New(s string) error {
	return &errWrap{
		Msg:      s,
		SumStack: []string{callerLine() + ": " + s},
	}
}

func Println(err error) {
	if err == nil {
		return
	}

	fmt.Print(Sprintln(err))
}

func Sprintln(err error) string {
	if err == nil {
		return ""
	}

	sErr, ok := err.(*errWrap)
	if !ok {
		return err.Error() + "\n"
	}

	return strings.Join(sErr.SumStack, "\n") + "\n"
}

// Stack wraps the given error with a stack line without a new message.
func Stack(err error) error {
	if err == nil {
		return nil
	}

	sErr, ok := err.(*errWrap)

	stackLine := callerLine()
	errMsg := err.Error()

	if !ok {
		return &errWrap{
			err:      err,
			Msg:      errMsg,
			SumStack: []string{stackLine},
		}
	}

	// if the errWrap is actually the original error, do not modify it. Make
	// a new err as the wrap, and embed the old one.
	//
	// We do this because currently we're not actually embedding all errors, but rather
	// only the original errors. Repeated Wrap() calls just modify the state of the
	// error msg/stack, but we don't want to modify the original error.
	if sErr.IsCause() {
		sErr = &errWrap{
			err:      sErr,
			Msg:      errMsg,
			SumStack: sErr.SumStack[:],
		}
	}

	// the given error is a wrapped error, modify it to the latest error information
	sErr.Msg = errMsg
	sErr.SumStack = append(sErr.SumStack, stackLine)
	return sErr
}

func Wrap(err error, s string) error {
	if err == nil {
		return nil
	}

	return wrap(err, callerLine(), s)
}

func Wrapf(err error, f string, s ...interface{}) error {
	if err == nil {
		return nil
	}

	return wrap(err, callerLine(), fmt.Sprintf(f, s...))
}

func wrap(err error, caller string, s string) error {
	sErr, ok := err.(*errWrap)

	// construct the cascaded error line
	cascadeErrLine := s + ": " + err.Error()

	// If it's not a wrapped error, construct a new one
	if !ok {
		// the previous err is *not* a wrapped error, so include the caller err.Error()
		// message in the stackLine.
		stackLine := caller + ": " + cascadeErrLine

		return &errWrap{
			err:      err,
			Msg:      cascadeErrLine,
			SumStack: []string{stackLine},
		}
	}

	// if the errWrap is actually the original error, do not modify it. Make
	// a new err as the wrap, and embed the old one.
	//
	// We do this because currently we're not actually embedding all errors, but rather
	// only the original errors. Repeated Wrap() calls just modify the state of the
	// error msg/stack, but we don't want to modify the original error.
	if sErr.IsCause() {
		sErr = &errWrap{
			err:      sErr,
			Msg:      sErr.Error(),
			SumStack: sErr.SumStack[:],
		}
	}

	// the given error is a wrapped error, *do not* include the
	// err.Error(), as that will consist of the cascade message.
	stackLine := caller + ": " + s

	// the given error is a wrapped error, modify it to the latest error information
	sErr.Msg = cascadeErrLine
	sErr.SumStack = append(sErr.SumStack, stackLine)
	return sErr
}

// callerLine returns a short path and the line number of the caller
//
// Currently it's returning the two directories above the file for brevity. In
// the future we may want to display the full path relative to the $GOPATH.
func callerLine() string {
	// p comes out as the absolute path, it needs to be trimmed.
	_, p, l, _ := runtime.Caller(2)
	// TODO(leeola): there has to be a better way to write this...
	p, f := path.Dir(p), path.Base(p)
	p, parent := path.Dir(p), path.Base(p)
	p, great := path.Dir(p), path.Base(p)
	p = path.Join(great, parent, f)
	return fmt.Sprintf("%s:%d", p, l)
}

type causer interface {
	Cause() error
}

type errWrap struct {
	// The original error that this error is wrapping. Stored to retrieve the
	// original as needed.
	err error

	Msg      string
	SumStack []string
}

func (e *errWrap) Error() string {
	return e.Msg
}

func (e *errWrap) Errors() []string {
	return e.SumStack[:]
}

func (e *errWrap) IsCause() bool {
	return e == e.Cause()
}

func (e *errWrap) Cause() error {
	if e.err != nil {
		return e.err
	}
	return e
}
