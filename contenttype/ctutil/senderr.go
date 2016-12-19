package ctutil

import ct "github.com/leeola/kala/contenttype"

// SendErr is a helper function allowing a caller to return an error to be sent.
//
// If no error is returned (ie, nil return value), then nothing is done.
func SendErrAndClose(ch chan ct.Result, err error) {
	if err != nil {
		ch <- ct.Result{Error: err}
		close(ch)
	}
}
