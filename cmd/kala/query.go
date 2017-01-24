package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
	"text/tabwriter"

	"github.com/leeola/kala/client"
	"github.com/leeola/kala/index"
	"github.com/leeola/kala/util/strutil"
	"github.com/urfave/cli"
)

func queryCommand(ctx *cli.Context) error {
	c, err := ClientFromContext(ctx)
	if err != nil {
		return err
	}

	q := index.Query{Metadata: index.Metadata{}}
	for _, arg := range ctx.Args() {
		k, v := strutil.SplitQueryField(arg)

		switch k {
		case "anchor", "previousHash", "multiPart", "multiHash":
			// Wrap the value in quotes for an easy UX. This will likely be changed
			// in the future, by opting out, so the user can explicitly declare quotes
			// as needed.
			v = "\"" + v + "\""
			q.Metadata[k] = v
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
		case "":
			if _, ok := q.Metadata["name"]; !ok {
				// if the key was not specified, default the value to name
				// only if name was not already set.
				//
				// Since this is implicit logic, we should not overwrite the users explicit
				// filename value.
				q.Metadata["name"] = v
			} else {
				fmt.Printf("warning: value with no key: %q\n", v)
			}
		default:
			// if the key was specified, set the requested metadata key with the val
			q.Metadata[k] = v
		}
	}

	// set the default index sort. Note that we're doing this
	// here because if we used StringSliceFlag.Value to set the default
	// then the user can't override the default. "indexEntry" would always be
	// the first ascending sort.
	//
	// So by setting it manually as the default, the user can override it.
	ascendingFlags := ctx.StringSlice("ascending")
	if len(ascendingFlags) == 0 && len(ctx.StringSlice("descending")) == 0 {
		ascendingFlags = []string{"indexEntry"}
	}

	sorts := []index.SortBy{}
	for _, s := range ascendingFlags {
		sorts = append(sorts, index.SortBy{Field: s})
	}
	for _, s := range ctx.StringSlice("descending") {
		sorts = append(sorts, index.SortBy{
			Field:      s,
			Descending: true,
		})
	}

	results, err := c.Query(q, sorts...)
	if err != nil {
		return err
	}

	previewFields := ctx.StringSlice("preview")
	if len(previewFields) == 0 {
		previewFields = []string{"name"}
	}

	previewColumns := make([]string, len(previewFields)+2)
	previewColumns[0] = "ENTRY"
	previewColumns[1] = "CONTENT ADDRESS"
	for i, p := range previewFields {
		i += 2
		previewColumns[i] = strings.ToUpper(p)
	}

	w := tabwriter.NewWriter(os.Stdout, 2, 0, 2, ' ', 0)
	fmt.Fprintln(w, strings.Join(previewColumns, "\t"))
	for _, h := range results.Hashes {
		previewValues := make([]interface{}, len(previewFields)+2)
		previewValues[0] = h.Entry
		previewValues[1] = h.Hash

		if len(previewFields) > 0 {
			blob, err := hashToMap(c, h.Hash)
			if err != nil {
				return err
			}

			for i, k := range previewFields {
				i += 2
				v, ok := blob[k]
				if ok {
					previewValues[i] = v
				} else {
					previewValues[i] = ""
				}
			}
		}

		var fStr string
		for i, _ := range previewValues {
			if i != 0 {
				fStr += "\t%s"
			} else {
				fStr += "%d"
			}
		}
		fStr += "\n"
		fmt.Fprintf(w, fStr, previewValues...)

	}
	return w.Flush()
}

func hashToMap(c *client.Client, h string) (map[string]interface{}, error) {
	rc, err := c.GetBlob(h)
	if err != nil {
		return nil, err
	}
	defer rc.Close()

	b, err := ioutil.ReadAll(rc)
	if err != nil {
		return nil, err
	}

	blob := map[string]interface{}{}
	if err := json.Unmarshal(b, &blob); err != nil {
		return nil, err
	}

	return blob, nil
}
