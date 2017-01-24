package webui

import (
	"github.com/leeola/errors"
	"github.com/leeola/kala/node/nodeware"
	"github.com/leeola/kala/webui/handlers"
	_ "github.com/leeola/kala/webui/statik"
	"github.com/leeola/kala/webui/webware"
	"github.com/rakyll/statik/fs"
)

func (w *WebUi) initRouter() error {
	w.router.Use(nodeware.LoggingMiddleware("webui", w.log))
	w.router.Use(webware.ClientMiddleware(w.client))
	w.router.Use(webware.TemplatersMiddleware(w.contentTemplaters))

	w.router.Get("/", handlers.GetRoot)
	w.router.Get("/hash/:hash", handlers.GetHash)
	w.router.Get("/hash/:hash/edit", handlers.GetHashEdit)
	w.router.Post("/hash/:hash/edit", handlers.PostHashEdit)
	w.router.Get("/search", handlers.GetSearch)

	statikFS, err := fs.New()
	if err != nil {
		return errors.Stack(err)
	}

	w.router.FileServer("/public", statikFS)

	return nil
}
