package regloader

import (
	"github.com/fatih/structs"
	"github.com/leeola/errors"
	"github.com/leeola/fixity"
	"github.com/leeola/fixity/autoload/registry"
	"github.com/leeola/fixity/indexes/snail"
	cu "github.com/leeola/fixity/util/configunmarshaller"
)

func init() {
	registry.RegisterIndex(Loader)
}

func Loader(cu cu.ConfigUnmarshaller) (fixity.Index, error) {
	var c struct {
		Config snail.Config `toml:"snailIndex"`
	}

	if err := cu.Unmarshal(&c); err != nil {
		return nil, errors.Wrap(err, "failed to unmarshal config")
	}

	// if the config isn't defined, do not load anything. This is allowed.
	if structs.IsZero(c.Config) {
		return nil, nil
	}

	return snail.New(c.Config)
}
