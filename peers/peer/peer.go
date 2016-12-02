package peer

import (
	"errors"
	"time"

	"github.com/inconshreveable/log15"
	"github.com/leeola/kala/client"
	"github.com/leeola/kala/index"
	"github.com/leeola/kala/store"
)

type PinQuery struct {
	// no query fields at the moment.
}

type Config struct {
	// The Kala client which this Pull struct will talk to
	Client *client.Client

	Store store.Store

	Frequency time.Duration

	Pins []PinQuery

	Log log15.Logger
}

// TODO(leeola): Implement a backoff method for both failure *and* success.
type Peer struct {
	// A peer has all the functionality of a client.
	*client.Client

	frequency time.Duration
	log       log15.Logger
	pins      []PinQuery
	store     store.Store

	version   string
	lastEntry int
}

func New(c Config) (*Peer, error) {
	if c.Client == nil {
		return nil, errors.New("missing required Config field: Client")
	}
	if c.Store == nil {
		return nil, errors.New("missing required Config field: Store")
	}

	if c.Frequency == 0 {
		c.Frequency = 5 * time.Minute
	}

	if c.Log == nil {
		c.Log = log15.New()
	}

	return &Peer{
		Client:    c.Client,
		frequency: c.Frequency,
		log:       c.Log,
		pins:      c.Pins,
		// TODO(leeola): remove in favor of a bolt backed entry.
		lastEntry: 1,
	}, nil
}

func (p *Peer) StartPinning() {
	// If there are no pins there's nothing to do.
	if p.pins == nil {
		return
	}

	go func() {
		// Switch this to NewTicker if we have a need to shut down the pulling
		// in the future.
		t := time.Tick(p.frequency)
		for {
			select {
			case <-t:
				if err := p.CheckPins(); err != nil {
					p.log.Error("failed get updated pins", "err", err)
				}
			}
		}
	}()
}

func (p *Peer) CheckPins() error {
	res, err := p.Query(index.Query{
		FromEntry: p.lastEntry,
		Limit:     3,
	})
	if err != nil {
		return err
	}

	p.log.Debug("got pin results", "peer")

	lenHashes := len(res.Hashes)
	if lenHashes == 0 {
		return nil
	}

	p.lastEntry += lenHashes

	for _, h := range res.Hashes {
		if err := p.PinHash(h); err != nil {
			return err
		}
	}

	return nil
}

func (p *Peer) PinHash(h string) error {
	p.log.Warn("PinHash not implemented, but letting that slide..")
	return nil
}

func (pq PinQuery) IsZero() bool {
	switch {
	default:
		return true
	}
}
