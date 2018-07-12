package fixity

import (
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

	// config doesn't exist, generate a default.
	c, err := DefaultConfig()
	if err != nil {
		return nil, fmt.Errorf("defaultconfig: %v", err)
	}

	return NewFromConfig(c)
}

func NewFromConfig(c config.Config) (Store, error) {
	if c.Store == "" {
		return nil, fmt.Errorf("missing required config: store")
	}

	tc, ok := c.StoreConfigs[c.Store]
	if !ok {
		return nil, fmt.Errorf("store name not found: %q", c.Store)
	}

	constructor, ok := storeRegistry[tc.Type]
	if !ok {
		return nil, fmt.Errorf("store type not found: %q", tc.Type)
	}

	s, err := constructor.New(c.Store, c)
	if err != nil {
		return nil, fmt.Errorf("store constructor %s: %v", c.Store, err)
	}

	return s, nil
}
