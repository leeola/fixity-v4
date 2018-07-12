package disk

import (
	"errors"

	"github.com/leeola/fixity"
	"github.com/leeola/fixity/config"
)

const configType = "disk"

func init() {
	fixity.RegisterBlobstore(configType, fixity.BlobstoreCreatorFunc(Constructor))
}

func Constructor(name string, c config.Config) (fixity.ReadWriter, error) {
	return nil, errors.New("not implemented")
}
