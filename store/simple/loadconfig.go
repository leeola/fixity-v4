package simple

import (
	"os"

	"github.com/BurntSushi/toml"
	"github.com/leeola/errors"
)

type simpleConfig struct {
	Config Config `toml:"simpleStore"`
}

func LoadConfig(configPath string) (Config, error) {
	f, err := os.Open(configPath)
	if err != nil {
		return Config{}, errors.Wrap(err, "failed to open config")
	}
	defer f.Close()

	var conf simpleConfig
	if _, err := toml.DecodeReader(f, &conf); err != nil {
		return Config{}, errors.Wrap(err, "failed to unmarshal config")
	}

	return conf.Config, nil
}
