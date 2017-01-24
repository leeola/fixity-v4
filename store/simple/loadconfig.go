package simple

import (
	"os"
	"path/filepath"

	"github.com/BurntSushi/toml"
	"github.com/leeola/errors"
	homedir "github.com/mitchellh/go-homedir"
)

var DefaultConfigFile = ConfigFile{}

type ConfigFile struct {
	RootPath           string
	CreateMissingPaths bool
	DontExpandHome     bool
	SimpleStore        Config `toml:"simpleStore"`
}

func LoadConfig(configPath string) (Config, error) {
	return LoadConfigWithDefault(configPath, DefaultConfigFile)
}

func LoadConfigWithDefault(configPath string, conf ConfigFile) (Config, error) {
	f, err := os.Open(configPath)
	if err != nil {
		return Config{}, errors.Wrap(err, "failed to open config")
	}
	defer f.Close()

	if _, err := toml.DecodeReader(f, &conf); err != nil {
		return Config{}, errors.Wrap(err, "failed to unmarshal config")
	}

	if !conf.DontExpandHome {
		if conf.SimpleStore.StorePath != "" {
			p, err := homedir.Expand(conf.SimpleStore.StorePath)
			if err != nil {
				return Config{}, errors.Stack(err)
			}
			conf.SimpleStore.StorePath = p
		}
		if conf.RootPath != "" {
			p, err := homedir.Expand(conf.RootPath)
			if err != nil {
				return Config{}, errors.Stack(err)
			}
			conf.RootPath = p
		}
	}

	if conf.RootPath != "" {
		conf.SimpleStore.StorePath = filepath.Join(conf.RootPath,
			conf.SimpleStore.StorePath)
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
