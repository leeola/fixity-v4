package autoload

import (
	"github.com/leeola/fixity"
	"github.com/leeola/fixity/autoload/registry"
	cu "github.com/leeola/fixity/util/configumarshaller"
)

func LoadFixity(configPath string) (fixity.Fixity, error) {
	cu := cu.New([]string{configPath})
	return registry.LoadFixity(cu)
}
