package blev

import (
	"os"
	"path/filepath"

	"github.com/BurntSushi/toml"
	"github.com/leeola/errors"
	homedir "github.com/mitchellh/go-homedir"
)

type ConfigFile struct {
	RootPath           string
	CreateMissingPaths bool
	DontExpandHome     bool
	BleveIndex         Config
}

func LoadConfig(configPath string) (Config, error) {
	return LoadConfigWithDefault(configPath, ConfigFile{})
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
		if conf.BleveIndex.BleveDir != "" {
			p, err := homedir.Expand(conf.BleveIndex.BleveDir)
			if err != nil {
				return Config{}, errors.Stack(err)
			}
			conf.BleveIndex.BleveDir = p
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
		conf.BleveIndex.BleveDir = filepath.Join(conf.RootPath,
			conf.BleveIndex.BleveDir)
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
