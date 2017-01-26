package node

import (
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/inconshreveable/log15"
	"github.com/leeola/kala/contenttype"
	ct "github.com/leeola/kala/contenttype"
	"github.com/leeola/kala/contenttype/defaults"
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

	// The index to provide indexing, querying and resetting for this Node.
	Index index.Index `toml:"-"`

	// optional
	ContentTypes map[string]ct.ContentType
	Router       *chi.Mux     `toml:"-"`
	Log          log15.Logger `toml:"-"`
}

type Node struct {
	bindAddr     string
	log          log15.Logger
	index        index.Index
	store        store.Store
	db           database.Database
	router       *chi.Mux
	contentTypes map[string]contenttype.ContentType
}

func New(c Config) (*Node, error) {
	if c.Index == nil {
		return nil, errors.New("missing required Config field: Index")
	}
	if c.Store == nil {
		return nil, errors.New("missing required Config field: Store")
	}
	if c.Database == nil {
		return nil, errors.New("missing required Config field: Database")
	}

	if c.BindAddr == "" {
		c.BindAddr = ":7001"
	}

	if c.ContentTypes == nil {
		css, err := defaults.DefaultTypes(c.Store, c.Index)
		if err != nil {
			return nil, err
		}
		c.ContentTypes = css
	}

	if c.Log == nil {
		c.Log = log15.New()
	}

	if c.Router == nil {
		c.Router = chi.NewRouter()
	}

	n := &Node{
		bindAddr:     c.BindAddr,
		log:          c.Log,
		index:        c.Index,
		store:        c.Store,
		db:           c.Database,
		router:       c.Router,
		contentTypes: c.ContentTypes,
	}

	if err := n.initDatabase(); err != nil {
		return nil, err
	}

	n.initRouter()

	// TODO(leeola): Rebuilding the index should be decided by if the store and
	// index are the same versions. Eg, if the store changes, it might be the same
	// index but the index doesn't match the store, so it should be rebuilt.
	//
	// The current solution doesn't handle that. To implement this, the store needs
	// to store a value, and the Node needs to check both Store and Index versions.
	if n.IsNewIndex() {
		if err := n.RebuildIndex(); err != nil {
			return nil, err
		}
	}

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

func (n *Node) ListenAndServe() error {
	n.log.Info("Node listening", "bindAddr", n.bindAddr)
	return http.ListenAndServe(n.bindAddr, n.router)
}

func (n *Node) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	n.router.ServeHTTP(w, r)
}
