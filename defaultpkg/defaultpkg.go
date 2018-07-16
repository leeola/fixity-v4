package defaultpkg

import (
	"github.com/leeola/fixity/blobstore/disk"
	"github.com/leeola/fixity/config"
	"github.com/leeola/fixity/config/log"
	"github.com/leeola/fixity/index/bleve"
	"github.com/leeola/fixity/store/nosign"
)

func init() {
	config.Configure(DefaultConfigure)
}

func DefaultConfigure(c config.Config) (config.Config, error) {
	c.Store = "default"
	c.Log = true
	c.LogLevel = log.Info
	c.BlobstoreConfigs["default"] = config.TypeConfig{
		Type: "disk",
		ConfigInterface: disk.Config{
			Path: "_store",
			Flat: true,
		},
	}

	c.IndexConfigs["default"] = config.TypeConfig{
		Type: "bleve",
		ConfigInterface: bleve.Config{
			Path: "_store",
		},
	}

	c.StoreConfigs["default"] = config.TypeConfig{
		Type: "nosign",
		ConfigInterface: nosign.Config{
			BlobstoreName: "default",
			IndexName:     "default",
		},
	}

	return c, nil
}
