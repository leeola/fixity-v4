package main

import (
	"fmt"
	"os"

	"github.com/leeola/errors"
)

func Foo() error {
	_, err := os.Open("fakefile.txt")
	if err != nil {
		return errors.Wrap(err, "os.Open failed")
	}
	return nil
}

func Bar() error {
	err := Foo()
	if err != nil {
		return errors.Wrap(err, "Foo() failed")
	}
	return nil
}

func Baz() error {
	// As an alternate Syntax, Wrap() will return nil if the given error
	// is nil. So you don't need to check `if err != nil {}` before wrapping.
	return errors.Wrap(Bar(), "Bar() failed")
}

func main() {
	err := Baz()

	if err != nil {
		fmt.Println("Baz() returned an error.\n")
		fmt.Println("err.Error():")
		fmt.Println(err.Error())
		fmt.Println("\nerrors.Println(err):")
		errors.Println(err)
		fmt.Println("\nerrors.Cause(err).Error():")
		fmt.Println(errors.Cause(err).Error())
	}
}
