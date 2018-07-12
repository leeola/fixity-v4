package config

import (
	"encoding/json"
	"fmt"
)

type Config struct {
	Store string `json:"store"`

	BlobstoreConfigs map[string]TypeConfig
	IndexConfigs     map[string]TypeConfig
	StoreConfigs     map[string]TypeConfig
}

type TypeConfig struct {
	Type            string `json:"type"`
	Config          json.RawMessage
	ConfigInterface interface{} `json:"-"`
}

func (c Config) BlobstoreConfig(key string, v interface{}) error {
	tc, ok := c.BlobstoreConfigs[key]
	if !ok {
		return fmt.Errorf("blobstore name not found: %q", key)
	}

	return json.Unmarshal(tc.Config, v)
}

func (c Config) IndexConfig(key string, v interface{}) error {
	tc, ok := c.IndexConfigs[key]
	if !ok {
		return fmt.Errorf("index name not found: %q", key)
	}

	return json.Unmarshal(tc.Config, v)
}

func (c Config) StoreConfig(key string, v interface{}) error {
	tc, ok := c.StoreConfigs[key]
	if !ok {
		return fmt.Errorf("store name not found: %q", key)
	}

	return json.Unmarshal(tc.Config, v)
}

func (c Config) MarshalInterfaces() (Config, error) {
	configs := map[string]map[string]TypeConfig{
		"blobstore": c.BlobstoreConfigs,
		"index":     c.IndexConfigs,
		"store":     c.StoreConfigs,
	}

	for configGroupName, configGroup := range configs {
		for k, tc := range configGroup {
			if tc.ConfigInterface != nil {
				b, err := json.Marshal(tc.ConfigInterface)
				if err != nil {
					return Config{}, fmt.Errorf("marshal %s config %q: %v", configGroupName, tc.Type, err)
				}
				tc.Config = b
				tc.ConfigInterface = nil
				configGroup[k] = tc
			}
		}
	}

	return c, nil
}
