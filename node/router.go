package node

func (n *Node) initRouter() {
	n.router.Use(LoggingMiddleware(n.log))

	n.router.Get("/id", n.GetNodeId)
	n.router.Head("/blob/:hash", n.HeadBlobHandler)
	n.router.Get("/blob/:hash", n.GetBlobHandler)
	n.router.Put("/blob/:hash", n.PutBlobHandler)
	n.router.Post("/blob", n.PostBlobHandler)
	n.router.Get("/index/query", n.GetQueryHandler)
	n.router.Get("/index/content", n.GetIndexContentHandler)
	n.router.Get("/download/:hash", n.GetDownloadHandler)
	n.router.Post("/upload", n.PostUploadHandler)
	n.router.Post("/upload/meta", n.PostUploadMetaHandler)
	// multihash and meta currently do the same exact thing, but the
	// api endpoint is being used to allow changes in UX specifically for
	// multihash mutation.
	n.router.Post("/upload/multihash", n.PostUploadMetaHandler)
}
