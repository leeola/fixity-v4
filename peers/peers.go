package peers

import (
	"errors"
	"io"

	"github.com/inconshreveable/log15"
	"github.com/leeola/kala/client"
	"github.com/leeola/kala/store"
)

type Config struct {
	PeerAddrs []string
	Store     store.Store
	Log       log15.Logger
}

// NOTE: Peers has a lot of room for optimization in how it reads/writes/etc
// to the peers. Eg: Read can easily be optimized to do exists checks on all,
// and only read from the match(es). Parallel execution of Exists checks and writes
// are another option.
//
// For now though, avoiding premature optimization in Peers.
type Peers struct {
	peerClients []*client.Client
	store       store.Store
	log         log15.Logger
}

func New(c Config) (*Peers, error) {
	if c.Store == nil {
		return nil, errors.New("missing Config field: Store")
	}

	if c.Log == nil {
		c.Log = log15.New()
	}

	clients := make([]*client.Client, len(c.PeerAddrs))
	for i, addr := range c.PeerAddrs {
		c, err := client.New(client.Config{
			KalaAddr: addr,
		})
		if err != nil {
			return nil, err
		}
		clients[i] = c
	}

	// TODO(leeola): sort peerClients based on a pre-configured means. Eg, network
	// locality, etc.

	return &Peers{
		peerClients: clients,
		store:       c.Store,
		log:         c.Log,
	}, nil
}

func (p *Peers) Exists(h string) (bool, error) {
	return false, errors.New("not implemented")
}

func (p *Peers) Read(h string) (io.ReadCloser, error) {
	rc, err := p.store.Read(h)
	if err != nil && err != store.HashNotFoundErr {
		return nil, err
	}

	if err == nil {
		return rc, nil
	}

	// if execution gets here, the local store could not find the hash. Attempt
	// to get it from the clients.
	//
	// This is a point of potential optimization. See Peers docstring for details.
	for _, c := range p.peerClients {
		peerRc, err := c.Read(h)

		// Note that we're dropping the error because we expect multiple peers to not
		// exist or not be reachable, and thus error at any point in time.
		//
		// TODO(leeola): figure out how to separate expected errors and unexpected
		// errors.
		if err != nil {
			p.log.Debug("peer Read error", "err", err)
		}

		if peerRc != nil {
			return peerRc, nil
		}
	}

	return nil, store.HashNotFoundErr
}

func (p *Peers) Write(b []byte) (string, error) {
	return p.store.Write(b)
}

func (p *Peers) WriteHash(h string, b []byte) error {
	return p.store.WriteHash(h, b)
}
