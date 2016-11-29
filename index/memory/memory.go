package memory

import (
	"strconv"
	"time"

	"github.com/leeola/errors"
	"github.com/leeola/kala/index"
	"github.com/leeola/kala/store"
)

type Config struct {
	// Indexes currently get their updates by acting as middleware between the
	// real Store and the Node. Due to this, Memory must implement Index *and* Store,
	// and the actual Store must be given to Memory so it can be wrapped.
	Store store.Store
}

type Memory struct {
	version    string
	entryCount int
	entries    map[int]string
	store      store.Store
}

func New(c Config) (*Memory, error) {
	if c.Store == nil {
		return nil, errors.New("missing required Config field: Store")
	}

	return &Memory{
		store: c.Store,
		// TODO(leeola): use hostname or uuid for the index versioning.
		version: strconv.FormatInt(time.Now().Unix(), 10),
		entries: map[int]string{},
	}, nil
}

func (m *Memory) QueryOne(q index.Query) (index.Result, error) {
	if q.IndexVersion != "" && m.version != q.IndexVersion {
		return index.Result{}, index.ErrIndexVersionsDoNotMatch
	}

	if q.IndexEntry != 0 {
		h, ok := m.entries[q.IndexEntry]
		if !ok {
			return index.Result{}, index.ErrNoQueryResults
		}

		return index.Result{
			IndexVersion: m.version,
			Hash:         h,
		}, nil
	}

	return index.Result{}, index.ErrNoQueryResults
}

func (m *Memory) Query(q index.Query) (index.Results, error) {
	if q.IndexVersion != "" && m.version != q.IndexVersion {
		return index.Results{}, index.ErrIndexVersionsDoNotMatch
	}

	// If IndexEntry was specified, there can only be one match.
	//
	// In the future when more query fields are added, the single entry
	// will have to be filtered against the other query fields.
	if q.IndexEntry != 0 {
		h, ok := m.entries[q.IndexEntry]
		if !ok {
			return index.Results{}, index.ErrNoQueryResults
		}

		return index.Results{
			IndexVersion: m.version,
			Hashes:       []string{h},
		}, nil
	}

	return index.Results{}, index.ErrNoQueryResults
}

func (m *Memory) Version() string {
	return m.version
}

func (m *Memory) AddEntry(h string) error {
	m.entryCount += 1
	m.entries[m.entryCount] = h
	return nil
}
