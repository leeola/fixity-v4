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
