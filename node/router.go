package node

import (
	"github.com/leeola/kala/node/handlers"
	"github.com/leeola/kala/node/nodeware"
)

func (n *Node) initRouter() {
	n.router.Use(nodeware.LoggingMiddleware("node", n.log))
	n.router.Use(nodeware.ContentStorersMiddleware(n.contentStorers))
	n.router.Use(nodeware.DatabaseMiddleware(n.db))
	n.router.Use(nodeware.IndexMiddleware(n.index))
	n.router.Use(nodeware.StoreMiddleware(n.store))
	n.router.Use(nodeware.QueryMiddleware(n.index))

	n.router.Get("/id", handlers.GetNodeId)

	n.router.Get("/index/query", handlers.GetQueryHandler)

	n.router.Post("/blob", handlers.PostBlobHandler)
	n.router.Head("/blob/:hash", handlers.HeadBlobHandler)
	n.router.Get("/blob/:hash", handlers.GetBlobHandler)
	n.router.Put("/blob/:hash", handlers.PutBlobHandler)
	n.router.Get("/blob/:hash/contenttype", handlers.GetBlobContentTypeHandler)

	n.router.Get("/download/:hash", handlers.GetDownloadHandler)
	// n.router.Get("/download/:hash/blob",
	// 	handlers.GetDownloadBlobHandler(n.store, n.index))
	// n.router.Get("/download/:hash/meta/export", handlers.GetMetaExportHandler)

	n.router.Post("/upload", handlers.PostUploadHandler)
	// n.router.Post("/upload/meta", n.PostUploadMetaHandler)
	// multihash and meta currently do the same exact thing, but the
	// api endpoint is being used to allow changes in UX specifically for
	// multihash mutation.
	// n.router.Post("/upload/multihash", n.PostUploadMetaHandler)
}
