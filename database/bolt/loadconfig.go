package bolt

import (
	"os"

	"github.com/BurntSushi/toml"
	"github.com/leeola/errors"
)

func LoadConfig(configPath string) (Config, error) {
	f, err := os.Open(configPath)
	if err != nil {
		return Config{}, errors.Wrap(err, "failed to open config")
	}
	defer f.Close()

	var conf struct {
		CreateMissingPaths bool
		BoltDatabase       Config
	}
	if _, err := toml.DecodeReader(f, &conf); err != nil {
		return Config{}, errors.Wrap(err, "failed to unmarshal config")
	}

	// Create the db path if it's missing.
	if conf.CreateMissingPaths && conf.BoltDatabase.BoltPath != "" {
		if _, err := os.Stat(conf.BoltDatabase.BoltPath); os.IsNotExist(err) {
			if err := os.MkdirAll(conf.BoltDatabase.BoltPath, 0755); err != nil {
				return Config{}, err
			}
		}
	}

	return conf.BoltDatabase, nil
}
