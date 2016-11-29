package peers

import (
	"os"

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
		BindAddr string
		Peers    []struct {
			Addr string
		} `toml:"peers"`
	}

	if _, err := toml.DecodeReader(f, &conf); err != nil {
		return Config{}, errors.Wrap(err, "failed to unmarshal config")
	}

	var peerAddrs []string
	if len(conf.Peers) != 0 {
		peerAddrs = make([]string, len(conf.Peers))
		for i, confPeer := range conf.Peers {
			peerAddrs[i] = confPeer.Addr
		}
	}

	return Config{
		PeerAddrs: peerAddrs,
	}, nil
}

func (c Config) IsZero() bool {
	switch {
	case c.PeerAddrs != nil:
		return false
	case c.Store != nil:
		return false
	default:
		return true
	}
}
