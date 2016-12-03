package main

import (
	"os"
	"strings"

	"github.com/BurntSushi/toml"
	"github.com/leeola/errors"
)

func LoadConfig(configPath string) (Config, error) {
	f, err := os.Open(configPath)
	if err != nil {
		return Config{}, errors.Wrap(err, "failed to open config")
	}
	defer f.Close()

	var conf struct {
		Config
		BindAddr string
	}
	if _, err := toml.DecodeReader(f, &conf); err != nil {
		return Config{}, errors.Wrap(err, "failed to unmarshal config")
	}

	if conf.Config.KalaAddr == "" && conf.BindAddr != "" {
		conf.Config.KalaAddr = BindToUrl(conf.BindAddr)
	}

	return conf.Config, nil
}

func BindToUrl(bindAddr string) string {
	switch {
	case strings.HasPrefix(bindAddr, ":"):
		// if the bindaddr is :8000
		return "http://localhost" + bindAddr
	default:
		// if the bind addr is something like localhost:8000
		return "http://" + bindAddr
	}
}
