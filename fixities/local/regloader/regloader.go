package regloader

import (
	"github.com/fatih/structs"
	"github.com/leeola/errors"
	"github.com/leeola/fixity"
	"github.com/leeola/fixity/autoload/registry"
	"github.com/leeola/fixity/impl/local"
	cu "github.com/leeola/fixity/util/configunmarshaller"
)

func init() {
	registry.RegisterFixity(Loader)
}

func Loader(cu cu.ConfigUnmarshaller) (fixity.Fixity, error) {
	// We're not even really using the rootPath, index and store uses those,
	// we mainly use the rootPath as a way to verify that the configuration
	// specifies a local Fixity implementation.
	var c struct {
		DontExpandHome bool         `toml:"dontExpandHome"`
		Config         local.Config `toml:"localFixity"`
	}

	if err := cu.Unmarshal(&c); err != nil {
		return nil, errors.Wrap(err, "failed to unmarshal config")
	}

	// if the config isn't defined, do not load anything. This is allowed.
	if structs.IsZero(c.Config) {
		return nil, nil
	}

	i, err := registry.LoadIndex(cu)
	if err != nil {
		return nil, err
	}

	s, err := registry.LoadStore(cu)
	if err != nil {
		return nil, err
	}

	config := c.Config
	config.Index = i
	config.Store = s

	return local.New(config)
}
