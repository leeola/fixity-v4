package autoload

import (
	"errors"

	"github.com/leeola/fixity"
	"github.com/leeola/fixity/autoload/registry"
	cu "github.com/leeola/fixity/util/configunmarshaller"

	_ "github.com/leeola/fixity/impl/local/regloader"
	_ "github.com/leeola/fixity/indexes/snail/regloader"
	_ "github.com/leeola/fixity/stores/disk/regloader"
)

func LoadFixity(configPath string) (fixity.Fixity, error) {
	if configPath == "" {
		return nil, errors.New("config path is required")
	}

	cu := cu.New([]string{configPath})
	return registry.LoadFixity(cu)
}
