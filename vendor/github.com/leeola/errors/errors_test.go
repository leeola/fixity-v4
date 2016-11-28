package errors

import (
	"testing"

	"errors"
)

// ensure that Wrap and Wrapf always return nil on a nil error or non-erroring code
// could get really weird for the caller.
func TestNilReturns(t *testing.T) {
	if err := Wrap(nil, "foo"); err != nil {
		t.Fatal("Wrap(nil) should have returned nil. got: ", err.Error())
	}

	if err := Wrapf(nil, "foo"); err != nil {
		t.Fatal("Wrapf(nil) should have returned nil. got: ", err.Error())
	}
}

func TestErrorsCauser(t *testing.T) {
	// make a normal error type
	origErr := errors.New("foo")
	wrappedErr := Wrap(origErr, "bar")
	cause := Cause(wrappedErr)

	if origErr != cause {
		t.Fatalf(
			"Expected original error to be the cause of the wrapped error. "+
				"wanted: %s, got: %s", origErr.Error(), cause.Error(),
		)
	}

	// nest the wraps several layers deep.
	origErr = errors.New("foo")
	wrappedErr = Wrap(origErr, "bar")
	wrappedErr2 := Wrap(wrappedErr, "baz")
	wrappedErr3 := Wrap(wrappedErr2, "bat")

	if origErr != Cause(wrappedErr3) {
		t.Fatalf(
			"Expected nested error to be the cause of the wrapped error. "+
				"wanted: %s, got: %s", origErr.Error(), cause.Error(),
		)
	}

	// ensure that when the original error is of leeola/errors.errWrap type, as is
	// the case with `errors.New()` and `errors.Errorf()` outputs, that we currently
	// return them as the causer.
	origErr = New("foo")
	wrappedErr = Wrap(origErr, "bar")

	if origErr.Error() != "foo" {
		t.Fatalf("Expected original error message not to be modified. "+
			"wanted: %s, got: %s", "foo", origErr.Error())
	}
	if origErr != Cause(wrappedErr) {
		t.Fatalf(
			"Expected original leeola/error error to be the cause of the wrapped error. "+
				"wanted: %s, got: %s", origErr.Error(), cause.Error(),
		)
	}

	// nest the wraps several layers deep
	origErr = New("foo")
	wrappedErr = Wrap(origErr, "bar")
	wrappedErr2 = Wrap(wrappedErr, "baz")
	wrappedErr3 = Wrap(wrappedErr2, "bat")

	if origErr.Error() != "foo" {
		t.Fatalf("Expected original error message not to be modified. "+
			"wanted: %s, got: %s", "foo", origErr.Error())
	}
	if origErr != Cause(wrappedErr3) {
		t.Fatalf(
			"Expected nested leeola/error error to be the cause of the wrapped error. "+
				"wanted: %s, got: %s", origErr.Error(), cause.Error(),
		)
	}
}
