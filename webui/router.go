package webui

import (
	"github.com/leeola/kala/node"
	"github.com/leeola/kala/webui/handlers"
	"github.com/leeola/kala/webui/webware"
)

func (w *WebUi) initRouter() {
	w.router.Use(node.LoggingMiddleware(w.log))
	w.router.Use(webware.ClientMiddleware(w.client))

	w.router.Get("/", handlers.GetRoot)
}
