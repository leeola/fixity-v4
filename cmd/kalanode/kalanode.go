package main

import (
	"flag"

	"github.com/leeola/kala/multiload"
)

func main() {
	var configPath string
	flag.StringVar(&configPath, "config", "./config.toml", "path to kala toml config")
	flag.Parse()

	s, err := multiload.LoadStoreWithDefault(configPath, multiload.RootConfig{})
	if err != nil {
		panic(err)
	}
	i, err := multiload.LoadIndexWithDefault(configPath, multiload.RootConfig{}, s)
	if err != nil {
		panic(err)
	}
	n, err := multiload.LoadNode(configPath, nil, i, s)
	if err != nil {
		panic(err)
	}

	// TODO(leeola): Move the peer init to the multiload
	//
	// wrap the store with our peers, if configured.
	// peersConfig, err := peers.LoadConfig(configPath)
	// if err != nil {
	// 	panic(err)
	// }
	// if !peersConfig.IsZero() {
	// 	peersConfig.Store = store
	// 	peersConfig.Database = db
	// 	p, err := peers.New(peersConfig)
	// 	if err != nil {
	// 		panic(err)
	// 	}

	// 	// start the pinning
	// 	p.StartPinning()

	// 	store = p
	// }

	if err := n.ListenAndServe(); err != nil {
		panic(err)
	}
}
