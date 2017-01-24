//go:generate statik -src=./public

package webui

import (
	"errors"
	"net/http"

	"github.com/inconshreveable/log15"
	"github.com/leeola/kala/client"
	"github.com/pressly/chi"
)

type Config struct {
	// The address for those node to listen on
	BindAddr string

	// The kala node to use for this UI.
	NodeAddr string

	// optional
	Client *client.Client
	Router *chi.Mux     `toml:"-"`
	Log    log15.Logger `toml:"-"`
}

type WebUi struct {
	bindAddr          string
	client            *client.Client
	log               log15.Logger
	router            *chi.Mux
	contentTemplaters map[string]interface{}
}

func New(c Config) (*WebUi, error) {
	if c.BindAddr == "" {
		return nil, errors.New("missing required Config field: BindAddr")
	}
	if c.NodeAddr == "" && c.Client == nil {
		return nil, errors.New("missing required Config field: NodeAddr")
	}

	if c.Log == nil {
		c.Log = log15.New()
	}

	if c.Router == nil {
		c.Router = chi.NewRouter()
	}

	if c.Client == nil {
		nClient, err := client.New(client.Config{
			KalaAddr: c.NodeAddr,
		})
		if err != nil {
			return nil, err
		}
		c.Client = nClient
	}

	w := &WebUi{
		bindAddr:          c.BindAddr,
		client:            c.Client,
		log:               c.Log,
		router:            c.Router,
		contentTemplaters: map[string]interface{}{},
	}

	if err := w.initRouter(); err != nil {
		return nil, err
	}

	return w, nil
}

func (w *WebUi) ListenAndServe() error {
	w.log.Info("WebUi listening", "bindAddr", w.bindAddr)
	return http.ListenAndServe(w.bindAddr, w.router)
}

func (w *WebUi) AddContentTemplater(t string, cs interface{}) {
	w.contentTemplaters[t] = cs
}
