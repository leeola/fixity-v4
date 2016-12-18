package main

import (
	"flag"

	"github.com/leeola/kala/webui"
)

func main() {
	var configPath string
	flag.StringVar(&configPath, "config", "./config.toml", "path to kala toml config")
	flag.Parse()

	webConfig, err := webui.LoadConfig(configPath)
	if err != nil {
		panic(err)
	}

	w, err := webui.New(webConfig)
	if err != nil {
		panic(err)
	}

	if err := w.ListenAndServe(); err != nil {
		panic(err)
	}
}
