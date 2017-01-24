package crypto

import (
	"bytes"
	"io"
	"io/ioutil"
)

func DecryptReadCloser(c Cryptoer, rc io.ReadCloser) (io.ReadCloser, error) {
	defer rc.Close()
	b, err := ioutil.ReadAll(rc)
	if err != nil {
		return nil, err
	}

	d, err := c.Decrypt(b)
	if err != nil {
		return nil, err
	}

	return ioutil.NopCloser(bytes.NewReader(d)), nil
}
