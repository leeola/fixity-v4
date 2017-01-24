package jsonutil

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"os"

	"github.com/leeola/errors"
)

func MarshalToWriter(w io.Writer, v interface{}) (int, error) {
	b, err := json.Marshal(v)
	if err != nil {
		return 0, err
	}

	n, err := w.Write(b)
	if err != nil {
		return 0, err
	}

	return n, nil
}

func MarshalToPath(p string, v interface{}) error {
	f, err := os.OpenFile(p, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = MarshalToWriter(f, v)
	return err
}

func UnmarshalReader(r io.Reader, v interface{}) error {
	b, err := ioutil.ReadAll(r)
	if err != nil {
		return errors.Stack(err)
	}

	if err := json.Unmarshal(b, v); err != nil {
		return errors.Stack(err)
	}

	return nil
}
