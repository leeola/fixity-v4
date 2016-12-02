package peers

import (
	"os"
	"time"

	"github.com/BurntSushi/toml"
	"github.com/leeola/errors"
	"github.com/leeola/kala/peers/peer"
)

func LoadConfig(configPath string) (Config, error) {
	f, err := os.Open(configPath)
	if err != nil {
		return Config{}, errors.Wrap(err, "failed to open config")
	}
	defer f.Close()

	var conf struct {
		Peers []struct {
			// Embedded PeerConfig
			PeerConfig

			// Add in a friendly default option to init a empty pin query, which
			// will pin everything.
			PinAll bool
		} `toml:"peers"`
	}

	if _, err := toml.DecodeReader(f, &conf); err != nil {
		return Config{}, errors.Wrap(err, "failed to unmarshal config")
	}

	var peerConfigs []PeerConfig
	if len(conf.Peers) != 0 {
		peerConfigs = make([]PeerConfig, len(conf.Peers))
		for i, peerStruct := range conf.Peers {
			peerConfig := peerStruct.PeerConfig
			if peerStruct.PinAll && !hasNoFilterPinQuery(peerConfig.Pins) {
				// Add an empty pin query, which by nature includes all hashes.
				// See Peer and PinQuery for further details.
				peerConfig.Pins = append(peerConfig.Pins, peer.PinQuery{})
			}

			// multiply the frequency by Seconds so that in the config it is based
			// off of seconds.
			peerConfig.Frequency = peerConfig.Frequency * time.Second

			peerConfigs[i] = peerConfig
		}
	}

	return Config{
		Peers: peerConfigs,
	}, nil
}

func (c Config) IsZero() bool {
	switch {
	case c.Peers != nil:
		return false
	case c.Store != nil:
		return false
	case c.Log != nil:
		return false
	default:
		return true
	}
}

func hasNoFilterPinQuery(pinQueries []peer.PinQuery) bool {
	if pinQueries == nil {
		return false
	}

	for _, pq := range pinQueries {
		if pq.IsZero() {
			return true
		}
	}

	return false
}
