package fixity

import (
	"fmt"
	"os"
	"syscall"
	"time"

	"github.com/leeola/fixity/config"
)

type Ref string

type Mutation struct {
	Schema
	ID           string    `json:"id"`
	Namespace    string    `json:"namespace"`
	Signer       string    `json:"signer"`
	Time         time.Time `json:"time"`
	ValuesSchema Ref       `json:"valuesSchema,omitempty"`
	DataSchema   Ref       `json:"dataSchema,omitempty"`
	Signature    string    `json:"signature"`
}

func New() (Store, error) {
	return NewFromConfigpath(config.DefaultConfigPath, "")
}

func NewFromConfigpath(storeName string, path string) (Store, error) {
	c, err := config.Open(path)
	if perr, ok := err.(*os.PathError); ok && perr.Err == syscall.ENOENT {
		// config doesn't exist, generate a default.
		c, err = NewConfig()
		if err != nil {
			return nil, fmt.Errorf("new config: %v", err)
		}

		if err := config.Save(path, c); err != nil {
			return nil, fmt.Errorf("save config: %v", err)
		}
	}
	if err != nil {
		return nil, fmt.Errorf("open config: %v", err)
	}

	return NewFromConfig(storeName, c)
}

func NewFromConfig(storeName string, c config.Config) (Store, error) {
	if storeName == "" {
		storeName = c.Store
	}
	if storeName == "" {
		return nil, fmt.Errorf("missing required argument: storeName")
	}

	tc, ok := c.StoreConfigs[storeName]
	if !ok {
		return nil, fmt.Errorf("store name not found: %q", storeName)
	}

	constructor, ok := storeRegistry[tc.Type]
	if !ok {
		return nil, fmt.Errorf("store type not found: %q", tc.Type)
	}

	s, err := constructor.New(storeName, c)
	if err != nil {
		return nil, fmt.Errorf("store constructor %s: %v", storeName, err)
	}

	return s, nil
}
