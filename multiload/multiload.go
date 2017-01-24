// MultiLoad loads various instances based on the given configuration.
//
// This is primarily used by cmd/kalanode and integration tests for easily spawning
// nodes based on various configs.

package multiload

import (
	"github.com/leeola/errors"
	"github.com/leeola/kala/database"
	"github.com/leeola/kala/database/bolt"
	"github.com/leeola/kala/index"
	"github.com/leeola/kala/index/blev"
	"github.com/leeola/kala/node"
	"github.com/leeola/kala/store"
	"github.com/leeola/kala/store/simple"
)

type RootConfig struct {
	RootPath string
	// Unmarshallers map[string]contenttype.MetadataUnmarshaller
}

func LoadNode(p string, d database.Database, i index.Index, s store.Store) (*node.Node, error) {
	return LoadNodeWithDefault(p, RootConfig{}, d, i, s)
}

func LoadNodeWithDefault(p string, c RootConfig, d database.Database, i index.Index, s store.Store) (*node.Node, error) {
	if s == nil {
		st, err := LoadStoreWithDefault(p, c)
		if err != nil {
			return nil, err
		}
		s = st
	}

	if d == nil {
		db, err := LoadDatabaseWithDefault(p, c)
		if err != nil {
			return nil, err
		}
		d = db
	}

	if i == nil {
		in, err := LoadIndexWithDefault(p, c, s)
		if err != nil {
			return nil, err
		}
		i = in
	}

	nodeConfig, err := node.LoadConfig(p)
	if err != nil {
		return nil, err
	}

	// fill the nodeConfig with the instances it needs to init.
	nodeConfig.Store = s
	nodeConfig.Index = i
	nodeConfig.Database = d
	// nodeConfig.Unmarshallers = c.Unmarshallers

	return node.New(nodeConfig)
}

func LoadDatabaseWithDefault(configPath string, c RootConfig) (database.Database, error) {
	// currently only bolt
	dbConfig, err := bolt.LoadConfigWithDefault(configPath, bolt.ConfigFile{
		RootPath: c.RootPath,
	})
	if err != nil {
		return nil, err
	}

	db, err := bolt.New(dbConfig)
	if err != nil {
		return nil, err
	}

	return db, nil
}

func LoadIndexWithDefault(configPath string, c RootConfig, s store.Store) (index.Index, error) {
	// currently only blev
	indexConfig, err := blev.LoadConfigWithDefault(configPath, blev.ConfigFile{
		RootPath: c.RootPath,
	})
	if err != nil {
		return nil, err
	}

	return blev.New(indexConfig)
}

func LoadStoreWithDefault(configPath string, c RootConfig) (store.Store, error) {
	// currently only simple
	simpleConfig, err := simple.LoadConfigWithDefault(configPath, simple.ConfigFile{
		RootPath: c.RootPath,
	})
	if err != nil {
		return nil, errors.Wrap(err, "failed to load simple store config")
	}

	return simple.New(simpleConfig)
}
