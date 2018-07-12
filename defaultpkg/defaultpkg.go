package defaultpkg

import (
	"github.com/leeola/fixity"
	"github.com/leeola/fixity/blobstore/disk"
	"github.com/leeola/fixity/config"
	"github.com/leeola/fixity/index/bleve"
	"github.com/leeola/fixity/store/nosign"
)

func init() {
	fixity.SetDefaultConfig(DefaultGenerator)
}

func DefaultGenerator() (config.Config, error) {
	return config.Config{
		Store: "default",
		BlobstoreConfigs: map[string]config.TypeConfig{
			"default": {
				Type: "disk",
				ConfigInterface: disk.Config{
					Path: "_store",
				},
			},
		},
		IndexConfigs: map[string]config.TypeConfig{
			"default": {
				Type: "bleve",
				ConfigInterface: bleve.Config{
					Path: "_store",
				},
			},
		},
		StoreConfigs: map[string]config.TypeConfig{
			"default": {
				Type: "nosign",
				ConfigInterface: nosign.Config{
					BlobstoreKey: "default",
					IndexKey:     "default",
				},
			},
		},
	}.MarshalInterfaces()
}
