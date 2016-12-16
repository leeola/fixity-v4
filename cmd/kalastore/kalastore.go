package main

import (
	"flag"

	"github.com/leeola/errors"
	"github.com/leeola/kala/contenttype"
	"github.com/leeola/kala/contenttype/data"
	"github.com/leeola/kala/contenttype/file"
	"github.com/leeola/kala/database/bolt"
	"github.com/leeola/kala/index"
	"github.com/leeola/kala/index/blev"
	"github.com/leeola/kala/node"
	"github.com/leeola/kala/peers"
	"github.com/leeola/kala/store"
	"github.com/leeola/kala/store/simple"
)

func main() {
	var configPath string
	flag.StringVar(&configPath, "config", "./config.toml", "path to kala toml config")
	flag.Parse()

	// init db (currently only bolt)
	dbConfig, err := bolt.LoadConfig(configPath)
	if err != nil {
		panic(err)
	}
	db, err := bolt.New(dbConfig)
	if err != nil {
		panic(err)
	}

	indexConfig, err := blev.LoadConfig(configPath)
	if err != nil {
		panic(err)
	}

	index, err := blev.New(indexConfig)
	if err != nil {
		panic(err)
	}

	// load a store specified in the config.
	store, err := initStoreFromConfig(configPath)
	if err != nil {
		panic(err)
	}

	// wrap the store with our peers, if configured.
	peersConfig, err := peers.LoadConfig(configPath)
	if err != nil {
		panic(err)
	}
	if !peersConfig.IsZero() {
		peersConfig.Store = store
		peersConfig.Database = db
		p, err := peers.New(peersConfig)
		if err != nil {
			panic(err)
		}

		// start the pinning
		p.StartPinning()

		store = p
	}

	nodeConfig, err := node.LoadConfig(configPath)
	if err != nil {
		panic(err)
	}

	// fill the nodeConfig with the instances it needs to init.
	nodeConfig.Store = store
	nodeConfig.Index = index
	nodeConfig.Query = index
	nodeConfig.Database = db

	n, err := node.New(nodeConfig)
	if err != nil {
		panic(err)
	}

	if err := addDefaultUploads(n, store, index); err != nil {
		panic(err)
	}

	if err := n.ListenAndServe(); err != nil {
		panic(err)
	}
}

func initStoreFromConfig(configPath string) (store.Store, error) {
	// first try the SimpleStore
	simpleConfig, err := simple.LoadConfig(configPath)
	if err != nil {
		return nil, errors.Wrap(err, "failed to load simple store config")
	}

	// if there is a config for the Simple store, use it.
	if !simpleConfig.IsZero() {
		simpleStore, err := simple.New(simpleConfig)
		// errors.Wrap() returns nil if err is nil, this is safe.
		return simpleStore, errors.Wrap(err, "failed to init simple store")
	}

	// no more store implementations to load from config.
	return nil, nil
}

func addDefaultUploads(n *node.Node, s store.Store, i index.Indexer) error {
	var cs contenttype.ContentStorer
	cs, err := data.New(data.Config{Store: s, Index: i})
	if err != nil {
		return err
	}
	n.AddContentStorer("data", cs)

	cs, err = file.New(file.Config{Store: s, Index: i})
	if err != nil {
		return err
	}
	n.AddContentStorer("file", cs)

	return nil
}
