package util

import "github.com/leeola/errors"

// MultiError is a quick implementation to concat errors.
func MultiError(errs ...error) error {
	var err error
	for _, e := range errs {
		if e == nil {
			continue
		}

		if err == nil {
			err = e
		} else {
			err = errors.Errorf("%s; %s", err.Error(), e.Error())
		}
	}
	return err
}
