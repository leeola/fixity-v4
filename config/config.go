package config

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"syscall"

	"github.com/leeola/fixity/config/log"
	homedir "github.com/mitchellh/go-homedir"
)

const DefaultConfigPath = "~/.config/fixity/config.json"

type Config struct {
	Store string `json:"store"`

	Log      bool      `json:"log"`
	LogLevel log.Level `json:"logLevel"`

	BlobstoreConfigs map[string]TypeConfig `json:"blobstoreConfigs"`
	IndexConfigs     map[string]TypeConfig `json:"indexConfigs"`
	StoreConfigs     map[string]TypeConfig `json:"storeConfigs"`
}

type TypeConfig struct {
	Type            string          `json:"type"`
	Config          json.RawMessage `json:"config"`
	ConfigInterface interface{}     `json:"-"`
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

func Open(path string) (Config, error) {
	path, err := homedir.Expand(path)
	if err != nil {
		return Config{}, fmt.Errorf("expand: %v", err)
	}

	f, err := os.OpenFile(path, os.O_RDONLY, 0644)
	if perr, ok := err.(*os.PathError); ok && perr.Err == syscall.ENOENT {
		return Config{}, ErrNotExist
	}
	if err != nil {
		return Config{}, fmt.Errorf("open: %v", err)
	}
	defer f.Close()

	b, err := ioutil.ReadAll(f)
	if err != nil {
		return Config{}, fmt.Errorf("read: %v", err)
	}

	var c Config
	if err := json.Unmarshal(b, &c); err != nil {
		return Config{}, fmt.Errorf("unmarshal: %v", err)
	}

	return c, nil
}

func Save(path string, c Config) error {
	path, err := homedir.Expand(path)
	if err != nil {
		return fmt.Errorf("expand: %v", err)
	}

	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return fmt.Errorf("mkdirall: %v", err)
	}

	f, err := os.OpenFile(path, os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0644)
	if err != nil {
		return fmt.Errorf("open: %v", err)
	}
	defer f.Close()

	b, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return fmt.Errorf("marshal: %v", err)
	}

	if _, err := io.Copy(f, bytes.NewReader(b)); err != nil {
		return fmt.Errorf("copy: %v", err)
	}

	if err := f.Sync(); err != nil {
		return fmt.Errorf("sync: %v", err)
	}

	return nil
}
