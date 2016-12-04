package main

import (
	"os"
	"strings"

	"github.com/BurntSushi/toml"
	"github.com/leeola/errors"
	"github.com/leeola/kala/client"
	homedir "github.com/mitchellh/go-homedir"
	"github.com/urfave/cli"
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

func ClientFromContext(c *cli.Context) (*client.Client, error) {
	config, _ := ConfigFromFile(c)
	if kalaAddr := c.String("KalaAddr"); kalaAddr != "" {
		config.KalaAddr = kalaAddr
	}

	if config.KalaAddr == "" {
		return nil, errors.New("missing KalaAddr")
	}

	client, err := client.New(client.Config{
		KalaAddr: config.KalaAddr,
	})
	if err != nil {
		return nil, err
	}

	return client, nil
}

func ConfigFromFile(c *cli.Context) (Config, error) {
	configPath, err := homedir.Expand(c.GlobalString("config"))
	if err != nil {
		return Config{}, err
	}

	conf, err := LoadConfig(configPath)
	if err != nil {
		return Config{}, err
	}

	return conf, nil
}
