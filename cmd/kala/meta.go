package main

import (
	"strings"

	"github.com/leeola/kala/store"
	"github.com/urfave/cli"
)

func metaCommand(c *cli.Context) error {
	client, err := ClientFromContext(c)
	if err != nil {
		return err
	}

	metaChanges := argsToMetaChanges(c.Args())

	hashes, err := client.UploadMeta(metaChanges)
	if err != nil {
		return err
	}

	for _, h := range hashes {
		Printlnf(h)
	}

	return nil
}

func argsToMetaChanges(args []string) store.MetaChanges {
	mc := store.MetaChanges{}
	argLen := len(args)
	for i := 0; i < argLen; i++ {
		key, value := keyValue(args[i])

		if value == "" {
			if i+1 < argLen {
				nextKey, nextValue := keyValue(args[i+1])

				if nextValue != "" {
					// if the next arg is a normal keyvalue arg (ie has both key and value),
					// set the current value to true, to treat it like a bool arg.
					value = "true"
				} else {
					// if the next value is empty, then it only has a key, so
					// treat the next key as the value to the current arg, and increment
					// the index.
					value = nextKey
					i++
				}
			} else {
				// There's no more args to look for a value from,
				// yet the value is empty, so treat the value as true
				value = "true"
			}
		}

		mc[key] = value
	}
	return mc
}

func keyValue(s string) (string, string) {
	sSplit := strings.SplitN(s, "=", 2)

	if len(sSplit) > 1 {
		return sSplit[0], sSplit[1]
	}

	return sSplit[0], ""
}
