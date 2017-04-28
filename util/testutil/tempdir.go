package testutil

import (
	"io/ioutil"
	"os"
)

func MustTempDir(prefix string) string {
	dir, err := ioutil.TempDir(os.TempDir(), prefix)
	if err != nil {
		panic(err)
	}
	return dir
}
