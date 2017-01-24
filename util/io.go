package util

import (
	"io"
	"io/ioutil"
	"os"
)

func WriteReader(r io.Reader, p string, perm os.FileMode) error {
	b, err := ioutil.ReadAll(r)
	if err != nil {
		return err
	}

	return ioutil.WriteFile(p, b, perm)
}
