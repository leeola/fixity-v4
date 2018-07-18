package config

import (
	"fmt"
	"sync"
)

var (
	configures   []func(Config) (Config, error)
	configuresMu sync.Mutex
)

func Configure(f func(Config) (Config, error)) {
	configuresMu.Lock()
	defer configuresMu.Unlock()

	configures = append(configures, f)
}

func NewConfig() (Config, error) {
	var (
		c = Config{
			BlobstoreConfigs: map[string]TypeConfig{},
			IndexConfigs:     map[string]TypeConfig{},
			StoreConfigs:     map[string]TypeConfig{},
		}
		err error
	)

	for _, f := range configures {
		c, err = f(c)
		if err != nil {
			return Config{}, err
		}
	}

	c, err = c.MarshalInterfaces()
	if err != nil {
		return Config{}, fmt.Errorf("marshal typeconfigs: %v", err)
	}

	return c, nil
}
