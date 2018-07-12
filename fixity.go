package fixity

import (
	"errors"
	"fmt"
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
	//
	// TODO(leeola): try to load the config here
	//

	if defaultConfig == nil {
		return nil, errors.New("no default config generator specified")
	}

	// config doesn't exist, generate a default.
	c, err := defaultConfig()
	if err != nil {
		return nil, fmt.Errorf("defaultConfigGen: %v", err)
	}

	return NewFromConfig(c.Store, c)
}

func NewFromConfig(name string, c config.Config) (Store, error) {
	if name == "" {
		return nil, fmt.Errorf("missing required config: store")
	}

	tc, ok := c.StoreConfigs[name]
	if !ok {
		return nil, fmt.Errorf("store name not found: %q", name)
	}

	constructor, ok := storeRegistry[tc.Type]
	if !ok {
		return nil, fmt.Errorf("store type not found: %q", tc.Type)
	}

	s, err := constructor.New(name, c)
	if err != nil {
		return nil, fmt.Errorf("store constructor %s: %v", name, err)
	}

	return s, nil
}
