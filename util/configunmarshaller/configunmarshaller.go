package configunmarshaller

import (
	"io"
	"io/ioutil"
	"os"
	"strings"

	"github.com/BurntSushi/toml"
	"github.com/leeola/errors"
)

type ConfigUnmarshaller interface {
	Unmarshal(interface{}) error
}

type ConfigUnmarshallerFunc func(interface{}) error

func (f ConfigUnmarshallerFunc) Unmarshal(v interface{}) error {
	return f(v)
}

// New returns an unmarshaller for the given configPaths.
//
// It reuses the bytes of the given configPaths, thereby reducing filesystem
// reads for repeated unmarshalls.
//
// Avoid storing the returned ConfigUnmarshaller as to avoid storing the
// configBytes in memory permanently.
func New(configPaths []string) ConfigUnmarshaller {
	var configData string

	return ConfigUnmarshallerFunc(func(v interface{}) error {
		if configData == "" {
			b, err := ConfigPathsToBytes(configPaths)
			if err != nil {
				return err
			}
			configData = string(b)
		}

		if _, err := toml.Decode(configData, v); err != nil {
			return errors.Stack(err)
		}

		return nil
	})
}

func UnmarshalConfigs(configPaths []string, v interface{}) error {
	b, err := ConfigPathsToBytes(configPaths)
	if err != nil {
		return errors.Stack(err)
	}

	if _, err := toml.Decode(string(b), v); err != nil {
		return errors.Stack(err)
	}

	return nil
}

func ConfigPathsToBytes(configPaths []string) ([]byte, error) {
	// *2 becuase we're adding a newline after every config to ensure no concat
	// issues.
	configReaders := make([]io.Reader, len(configPaths)*2)
	for i, p := range configPaths {
		f, err := os.Open(p)
		if err != nil {
			return nil, errors.Wrapf(err, "failed to open config: %s", p)
		}
		defer f.Close()

		configReaders[i] = f
		configReaders[i+1] = strings.NewReader("\n")
	}

	return ioutil.ReadAll(io.MultiReader(configReaders...))
}
