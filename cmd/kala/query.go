package main

import (
	"errors"
	"sort"
	"strconv"
	"strings"

	"github.com/leeola/kala/index"
	"github.com/urfave/cli"
)

func queryCommand(c *cli.Context) error {
	client, err := ClientFromContext(c)
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
			if q.Metadata == nil {
				q.Metadata = index.Metadata{}
			}
			q.Metadata[k] = v
		}
	}

	results, err := client.Query(q)
	if err != nil {
		return err
	}

	sort.Sort(results.Hashes)

	for _, h := range results.Hashes {
		Printlnf("#%d %s", h.Entry, h.Hash)
	}

	return nil
}
