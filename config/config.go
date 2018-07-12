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
	Type   string `json:"type"`
	Config json.RawMessage
}

func (c Config) StoreConfig(key string, v interface{}) error {
	tc, ok := c.StoreConfigs[key]
	if !ok {
		return fmt.Errorf("store name not found: %q", key)
	}

	return json.Unmarshal(tc.Config, v)
}
