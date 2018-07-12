package disk

import (
	"github.com/leeola/fixity"
	"github.com/leeola/fixity/config"
)

const configType = "disk"

func init() {
	fixity.RegisterBlobstore(configType, fixity.BlobstoreCreatorFunc(Constructor))
}

func Constructor(n string, c config.Config) (fixity.ReadWriter, error) {
	return New(n, c)
}
