package main

import (
	"errors"
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

	// set the default index sort. Note that we're doing this
	// here because if we used StringSliceFlag.Value to set the default
	// then the user can't override the default. "indexEntry" would always be
	// the first ascending sort.
	//
	// So by setting it manually as the default, the user can override it.
	ascendingFlags := c.StringSlice("ascending")
	if len(ascendingFlags) == 0 && len(c.StringSlice("descending")) == 0 {
		ascendingFlags = []string{"indexEntry"}
	}

	sorts := []index.SortBy{}
	for _, s := range ascendingFlags {
		sorts = append(sorts, index.SortBy{Field: s})
	}
	for _, s := range c.StringSlice("descending") {
		sorts = append(sorts, index.SortBy{
			Field:      s,
			Descending: true,
		})
	}

	results, err := client.Query(q, sorts)
	if err != nil {
		return err
	}

	for _, h := range results.Hashes {
		Printlnf("#%d %s", h.Entry, h.Hash)
	}

	return nil
}
