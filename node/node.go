package node

import (
	"errors"
	"net/http"

	"github.com/inconshreveable/log15"
	"github.com/leeola/kala/store"
	"github.com/pressly/chi"
)

type Config struct {
	// The address for those node to listen on
	BindAddr string

	// required
	Store store.Store `toml:"-"`

	// optional
	Router *chi.Mux
	Log    log15.Logger `toml:"-"`
}

type Node struct {
	bindAddr string
	log      log15.Logger
	store    store.Store
	router   *chi.Mux
}

func New(c Config) (*Node, error) {
	if c.BindAddr == "" {
		return nil, errors.New("missing required Config field: BindAddr")
	}
	if c.Store == nil {
		return nil, errors.New("missing required Config field: Store")
	}

	if c.Log == nil {
		c.Log = log15.New()
	}

	if c.Router == nil {
		c.Router = chi.NewRouter()
	}

	n := &Node{
		bindAddr: c.BindAddr,
		log:      c.Log,
		store:    c.Store,
		router:   c.Router,
	}

	n.initRouter()

	return n, nil
}

func (n *Node) ListenAndServe() error {
	n.log.Info("Node listening", "bindAddr", n.bindAddr)
	return http.ListenAndServe(n.bindAddr, n.router)
}
