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

	if !conf.DontExpandHome && conf.SimpleStore.Path != "" {
		p, err := homedir.Expand(conf.SimpleStore.Path)
		if err != nil {
			return Config{}, errors.Stack(err)
		}
		conf.SimpleStore.Path = p
	}

	// Create the db path if it's missing.
	if conf.CreateMissingPaths && conf.SimpleStore.Path != "" {
		if _, err := os.Stat(conf.SimpleStore.Path); os.IsNotExist(err) {
			if err := os.MkdirAll(conf.SimpleStore.Path, 0755); err != nil {
				return Config{}, err
			}
		}
	}

	return conf.SimpleStore, nil
}
