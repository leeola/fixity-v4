package main

import (
	"github.com/leeola/kala/client"
	homedir "github.com/mitchellh/go-homedir"
	"github.com/urfave/cli"
)

func idCommand(c *cli.Context) error {
	configPath, err := homedir.Expand(c.GlobalString("config"))
	if err != nil {
		return err
	}

	conf, err := LoadConfig(configPath)
	if err != nil {
		return err
	}

	client, err := client.New(client.Config{
		KalaAddr: conf.KalaAddr,
	})
	if err != nil {
		return err
	}

	id, err := client.NodeId()
	if err != nil {
		return err
	}

	Printlnf(id)

	return nil
}
