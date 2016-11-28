package main

import (
	"flag"

	"github.com/leeola/errors"
	"github.com/leeola/kala/node"
	"github.com/leeola/kala/store"
	"github.com/leeola/kala/store/simple"
)

func main() {
	var configPath string
	flag.StringVar(&configPath, "config", "./config.toml", "path to kala toml config")
	flag.Parse()

	// load a store specified in the config.
	store, err := initStoreFromConfig(configPath)
	if err != nil {
		panic(err)
	}

	nodeConfig, err := node.LoadConfig(configPath)
	if err != nil {
		panic(err)
	}

	// fill the nodeConfig with the instances it needs to init.
	nodeConfig.Store = store

	n, err := node.New(nodeConfig)
	if err != nil {
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
