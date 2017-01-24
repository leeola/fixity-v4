package main

import (
	"flag"

	"github.com/leeola/kala/client"
	"github.com/leeola/kala/contenttype/inventory"
	"github.com/leeola/kala/webui"
)

func main() {
	var configPath string
	flag.StringVar(&configPath, "config", "./config.toml", "path to kala toml config")
	flag.Parse()

	webConfig, err := webui.LoadConfig(configPath)
	if err != nil {
		panic(err)
	}

	nClient, err := client.New(client.Config{
		KalaAddr: webConfig.NodeAddr,
	})
	if err != nil {
		panic(err)
	}
	webConfig.Client = nClient

	w, err := webui.New(webConfig)
	if err != nil {
		panic(err)
	}

	if err := addDefaultTemplaters(w, nClient); err != nil {
		panic(err)
	}

	if err := w.ListenAndServe(); err != nil {
		panic(err)
	}
}

func addDefaultTemplaters(w *webui.WebUi, c *client.Client) error {
	var t interface{}
	// Temporarily disabled while the interface changes.
	// t = contenttype.TemplaterFunc(data.Templater)
	// w.AddContentTemplater("data", t)

	// t = contenttype.TemplaterFunc(folder.Templater)
	// w.AddContentTemplater("folder", t)

	// t = contenttype.TemplaterFunc(file.Templater)
	// w.AddContentTemplater("file", t)

	// t = contenttype.TemplaterFunc(image.Templater)
	// w.AddContentTemplater("image", t)

	// t = contenttype.TemplaterFunc(video.Templater)
	// w.AddContentTemplater("video", t)

	t = inventory.NewTemplater(c)
	w.AddContentTemplater(inventory.TypeKey, t)

	return nil
}
