package disk

import (
	"github.com/leeola/fixity"
	"github.com/leeola/fixity/config"
)

const configType = "disk"

func init() {
	fixity.RegisterBlobstore(configType, fixity.BlobstoreConstructorFunc(Constructor))
}

func Constructor(n string, c config.Config) (fixity.Blobstore, error) {
	return New(n, c)
}
