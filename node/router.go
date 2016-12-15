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
	n.router.Post("/upload/meta", n.PostUploadMetaHandler)
	n.router.Post("/upload/:contentType", n.PostUploadHandler)
	n.router.Get("/download/:hash", n.GetDownloadHandler)
}
