package node

func (n *Node) initRouter() {
	n.router.Use(LoggingMiddleware(n.log))

	n.router.Get("/content/:hash", n.GetContentHandler)
	n.router.Put("/content/:hash", n.PutContentHandler)
	n.router.Post("/content", n.PostContentHandler)
}
