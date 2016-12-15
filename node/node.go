package node

import (
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/inconshreveable/log15"
	"github.com/leeola/kala/contenttype"
	"github.com/leeola/kala/database"
	"github.com/leeola/kala/index"
	"github.com/leeola/kala/store"
	"github.com/pressly/chi"
)

type Config struct {
	// The address for those node to listen on
	BindAddr string

	// The store to provide content for this Node.
	Store store.Store `toml:"-"`

	// The database to present this nodes id from.
	Database database.Database `toml:"-"`

	// The indexer to provide index for this Node.
	Index index.Indexer `toml:"-"`

	// The queryer to provide content queries for this Node.
	Query index.Queryer `toml:"-"`

	// optional
	Router *chi.Mux     `toml:"-"`
	Log    log15.Logger `toml:"-"`
}

type Node struct {
	bindAddr string
	log      log15.Logger
	index    index.Indexer
	query    index.Queryer
	store    store.Store
	db       database.Database
	router   *chi.Mux

	contentStorers map[string]contenttype.ContentStorer
}

func New(c Config) (*Node, error) {
	if c.BindAddr == "" {
		return nil, errors.New("missing required Config field: BindAddr")
	}
	if c.Index == nil {
		return nil, errors.New("missing required Config field: Index")
	}
	if c.Query == nil {
		return nil, errors.New("missing required Config field: Query")
	}
	if c.Store == nil {
		return nil, errors.New("missing required Config field: Store")
	}
	if c.Database == nil {
		return nil, errors.New("missing required Config field: Database")
	}

	if c.Log == nil {
		c.Log = log15.New()
	}

	if c.Router == nil {
		c.Router = chi.NewRouter()
	}

	n := &Node{
		bindAddr:       c.BindAddr,
		log:            c.Log,
		index:          c.Index,
		query:          c.Query,
		store:          c.Store,
		db:             c.Database,
		router:         c.Router,
		contentStorers: map[string]contenttype.ContentStorer{},
	}

	if err := n.initDatabase(); err != nil {
		return nil, err
	}

	n.initRouter()

	return n, nil
}

// initDatabase ensures a series of values exist in the db for this node to use.
func (n *Node) initDatabase() error {
	_, err := n.db.GetNodeId()
	if err != nil && err != database.ErrNoRecord {
		return err
	}

	if err == database.ErrNoRecord {
		// TODO(leeola): use hostname+timestamp or uuid for the node.
		nodeId := strconv.FormatInt(time.Now().Unix(), 10)
		if err := n.db.SetNodeId(nodeId); err != nil {
			return err
		}
	}

	return nil
}

func (n *Node) AddUploader(t string, cs contenttype.ContentStorer) {
	n.contentStorers[t] = cs
}

func (n *Node) ListenAndServe() error {
	n.log.Info("Node listening", "bindAddr", n.bindAddr)
	return http.ListenAndServe(n.bindAddr, n.router)
}
