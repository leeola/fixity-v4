package blev

import (
	"os"

	"github.com/BurntSushi/toml"
	"github.com/leeola/errors"
	homedir "github.com/mitchellh/go-homedir"
)

func LoadConfig(configPath string) (Config, error) {
	f, err := os.Open(configPath)
	if err != nil {
		return Config{}, errors.Wrap(err, "failed to open config")
	}
	defer f.Close()

	var conf struct {
		CreateMissingPaths bool
		DontExpandHome     bool
		BleveIndex         Config
	}
	if _, err := toml.DecodeReader(f, &conf); err != nil {
		return Config{}, errors.Wrap(err, "failed to unmarshal config")
	}

	if !conf.DontExpandHome && conf.BleveIndex.BleveDir != "" {
		p, err := homedir.Expand(conf.BleveIndex.BleveDir)
		if err != nil {
			return Config{}, errors.Stack(err)
		}
		conf.BleveIndex.BleveDir = p
	}

	// Create the db path if it's missing.
	if conf.CreateMissingPaths && conf.BleveIndex.BleveDir != "" {
		if _, err := os.Stat(conf.BleveIndex.BleveDir); os.IsNotExist(err) {
			if err := os.MkdirAll(conf.BleveIndex.BleveDir, 0755); err != nil {
				return Config{}, err
			}
		}
	}

	return conf.BleveIndex, nil
}
