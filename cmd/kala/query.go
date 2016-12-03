package main

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/leeola/kala/client"
	"github.com/leeola/kala/index"
	homedir "github.com/mitchellh/go-homedir"
	"github.com/urfave/cli"
)

func queryCommand(c *cli.Context) error {
	configPath, err := homedir.Expand(c.GlobalString("config"))
	if err != nil {
		fmt.Println(err)
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

	q := index.Query{}
	for _, arg := range c.Args() {
		argSplit := strings.SplitN(arg, "=", 2)
		k := argSplit[0]
		var v string
		if len(argSplit) > 1 {
			v = argSplit[1]
		}

		switch k {
		case "fromEntry":
			i, err := strconv.Atoi(v)
			if err != nil {
				return errors.New("fromEntry must be an integer")
			}
			q.FromEntry = i
		case "limit":
			i, err := strconv.Atoi(v)
			if err != nil {
				return errors.New("limit must be an integer")
			}
			q.Limit = i
		case "indexVersion":
			q.IndexVersion = v
		default:
			Printlnf("warning: unhandled query argument: %s=%s", k, v)
		}
	}

	results, err := client.Query(q)
	if err != nil {
		return err
	}

	for _, h := range results.Hashes {
		Printlnf(h)
	}

	return nil
}
