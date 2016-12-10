package peer

import (
	"time"

	"github.com/inconshreveable/log15"
	"github.com/leeola/errors"
	"github.com/leeola/kala/client"
	"github.com/leeola/kala/database"
	"github.com/leeola/kala/index"
	"github.com/leeola/kala/store"
)

const (
	defaultFrequency     = 5 * time.Minute
	defaultPinQueryLimit = 3
)

type Config struct {
	// The Kala client which this Pull struct will talk to
	Client *client.Client

	Store store.Store

	// The database is used mainly to store the last index received for each
	// pinquery that a this peer is configured to pin from other nodes.
	Database database.Database

	Frequency time.Duration

	Pins []index.PinQuery

	// Optional. The number of hashes this peer will query for at a time.
	PinQueryLimit int

	Log log15.Logger
}

// TODO(leeola): Implement a backoff method for both failure *and* success.
type Peer struct {
	// A peer has all the functionality of a client.
	*client.Client

	frequency     time.Duration
	log           log15.Logger
	pins          []index.PinQuery
	store         store.Store
	database      database.Database
	pinQueryLimit int
}

func New(c Config) (*Peer, error) {
	if c.Client == nil {
		return nil, errors.New("missing required Config field: Client")
	}
	if c.Store == nil {
		return nil, errors.New("missing required Config field: Store")
	}
	if c.Database == nil {
		return nil, errors.New("missing required Config field: Database")
	}

	if c.Frequency == 0 {
		c.Frequency = defaultFrequency
	}

	if c.PinQueryLimit == 0 {
		c.PinQueryLimit = defaultPinQueryLimit
	}

	if c.Log == nil {
		c.Log = log15.New()
	}

	return &Peer{
		Client:        c.Client,
		frequency:     c.Frequency,
		log:           c.Log,
		pins:          c.Pins,
		store:         c.Store,
		database:      c.Database,
		pinQueryLimit: c.PinQueryLimit,
	}, nil
}

func (p *Peer) StartPinning() {
	// If there are no pins there's nothing to do.
	if len(p.pins) == 0 {
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
	// If there are no pins there's nothing to do.
	if len(p.pins) == 0 {
		return nil
	}

	for _, pin := range p.pins {
		if err := p.checkPin(pin); err != nil {
			return err
		}
	}

	return nil
}

func (p *Peer) checkPin(pin index.PinQuery) error {
	peerId, err := p.NodeId()
	if err != nil {
		return errors.Wrap(err, "failed to get nodeId from peer")
	}
	lastEntry, err := p.database.GetPeerPinLastEntry(peerId, pin)
	if err != nil && err != database.ErrNoRecord {
		return errors.Wrap(err, "failed to get lastEntry for pin")
	}

	// if there was no record of the last entry, we start from the first entry.
	// Note that it may not actually exists, and that's okay. Zero is a bad
	// number to start from.
	if err == database.ErrNoRecord {
		lastEntry = 1
	}

	res, err := p.Query(index.Query{
		FromEntry: lastEntry,
		Limit:     p.pinQueryLimit,
	})
	if err != nil {
		return err
	}

	lenHashes := len(res.Hashes)
	if lenHashes == 0 {
		return nil
	}

	var highestEntry int
	for i, h := range res.Hashes {
		if err := p.PinHash(h.Hash); err != nil {
			return err
		}

		// TODO(leeola): Once increment highestEntry based off of the eventual Hash
		// struct which also includes the indexEntry of the actual struct!
		// if h.IndexEntry > highestEntry {
		// 	highestEntry = h.IndexEntry
		// }
		highestEntry = lastEntry + i
	}

	// in the future we may want to save per successful pin?
	err = p.database.SetPeerPinLastEntry(peerId, pin, highestEntry+1)
	if err != nil {
		return errors.Wrap(err, "failed to write new peer last entry")
	}

	return nil
}

func (p *Peer) PinHash(h string) error {
	r, err := p.Read(h)
	if err != nil {
		return errors.Wrap(err, "failed to read hash from peer")
	}

	if err := store.WriteHashReader(p.store, h, r); err != nil {
		return errors.Wrap(err, "failed to write peers bytes to local store")
	}

	p.log.Info("pinned blob from peer", "hash", h)

	return nil
}
