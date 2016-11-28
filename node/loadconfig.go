package node

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

	var conf Config
	if _, err := toml.DecodeReader(f, &conf); err != nil {
		return Config{}, errors.Wrap(err, "failed to unmarshal config")
	}

	return conf, nil
}
