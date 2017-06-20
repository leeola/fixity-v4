package fixityutil

import (
	"encoding/json"
	"io/ioutil"

	"github.com/leeola/fixity"
)

func ReadAndUnmarshal(f fixity.Fixity, h string, v interface{}) error {
	rc, err := f.Blob(h)
	if err != nil {
		return err
	}
	defer rc.Close()

	b, err := ioutil.ReadAll(rc)
	if err != nil {
		return err
	}

	if err := json.Unmarshal(b, v); err != nil {
		return err
	}

	return nil
}
