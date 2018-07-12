package fixity

import (
	"errors"

	"github.com/leeola/fixity/config"
)

var defaultConfig func() (config.Config, error)

func SetDefaultConfig(f func() (config.Config, error)) {
	defaultConfig = f
}

func DefaultConfig() (config.Config, error) {
	if defaultConfig == nil {
		return config.Config{}, errors.New("no default configuration defined")
	}

	return defaultConfig()
}
