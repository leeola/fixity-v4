package peers

import (
	"io"
	"time"

	"github.com/inconshreveable/log15"
	"github.com/leeola/errors"
	"github.com/leeola/kala/client"
	"github.com/leeola/kala/database"
	"github.com/leeola/kala/index"
	"github.com/leeola/kala/peers/peer"
	"github.com/leeola/kala/store"
)

type Config struct {
	Peers []PeerConfig
	Store store.Store

	// The database is used mainly to store the last index received for each
	// pinquery that a this node is configured to pin from other nodes.
	Database database.Database

	Log log15.Logger
}

// Note that this exists to construct an individual Peer, in combination with the
// Peers Config. This should not be confused with Peer.Config.
type PeerConfig struct {
	Addr      string
	Frequency time.Duration
	Pins      []index.PinQuery
}

// NOTE: Peers has a lot of room for optimization in how it reads/writes/etc
// to the peers. Eg: Read can easily be optimized to do exists checks on all,
// and only read from the match(es). Parallel execution of Exists checks and writes
// are another option.
//
// For now though, avoiding premature optimization in Peers.
type Peers struct {
	peers []*peer.Peer
	store store.Store
	log   log15.Logger
}

func New(c Config) (*Peers, error) {
	if c.Store == nil {
		return nil, errors.New("missing Config field: Store")
	}
	if c.Database == nil {
		return nil, errors.New("missing Config field: Database")
	}

	if c.Log == nil {
		c.Log = log15.New()
	}

	peers := make([]*peer.Peer, len(c.Peers))
	for i, pc := range c.Peers {
		kc, err := client.New(client.Config{
			KalaAddr: pc.Addr,
		})
		if err != nil {
			return nil, errors.Wrap(err, "failed to construct kala client")
		}

		p, err := peer.New(peer.Config{
			Client:    kc,
			Pins:      pc.Pins,
			Store:     c.Store,
			Database:  c.Database,
			Log:       c.Log,
			Frequency: pc.Frequency,
		})
		if err != nil {
			return nil, errors.Wrap(err, "failed to construct peer")
		}

		peers[i] = p
	}

	// TODO(leeola): sort peerClients based on a pre-configured means. Eg, network
	// locality, etc.

	return &Peers{
		peers: peers,
		store: c.Store,
		log:   c.Log,
	}, nil
}

func (p *Peers) Exists(h string) (bool, error) {
	return p.store.Exists(h)
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
	for _, c := range p.peers {
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

func (ps *Peers) StartPinning() {
	ps.log.Info("Peers are starting to pin")

	for _, p := range ps.peers {
		p.StartPinning()
	}
}
