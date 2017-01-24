package webui

import (
	"os"

	"github.com/BurntSushi/toml"
	"github.com/leeola/errors"
	"github.com/leeola/kala/client"
)

func LoadConfig(configPath string) (Config, error) {
	f, err := os.Open(configPath)
	if err != nil {
		return Config{}, errors.Wrap(err, "failed to open config")
	}
	defer f.Close()

	var conf struct {
		BindAddr string
		Web      Config
	}
	if _, err := toml.DecodeReader(f, &conf); err != nil {
		return Config{}, errors.Wrap(err, "failed to unmarshal config")
	}

	if conf.Web.NodeAddr == "" && conf.BindAddr != "" {
		conf.Web.NodeAddr = client.BindToHttp(conf.BindAddr)
	}

	return conf.Web, nil
}
