package simple

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
		SimpleStore        Config `toml:"simpleStore"`
	}
	if _, err := toml.DecodeReader(f, &conf); err != nil {
		return Config{}, errors.Wrap(err, "failed to unmarshal config")
	}

	if !conf.DontExpandHome && conf.SimpleStore.StorePath != "" {
		p, err := homedir.Expand(conf.SimpleStore.StorePath)
		if err != nil {
			return Config{}, errors.Stack(err)
		}
		conf.SimpleStore.StorePath = p
	}

	// Create the db path if it's missing.
	if conf.CreateMissingPaths && conf.SimpleStore.StorePath != "" {
		if _, err := os.Stat(conf.SimpleStore.StorePath); os.IsNotExist(err) {
			if err := os.MkdirAll(conf.SimpleStore.StorePath, 0755); err != nil {
				return Config{}, err
			}
		}
	}

	return conf.SimpleStore, nil
}
