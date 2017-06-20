package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strconv"
	"strings"

	"github.com/leeola/fixity"
	"github.com/leeola/fixity/util/fixityutil"
	"github.com/urfave/cli"
)

func WriteCmd(ctx *cli.Context) error {
	if len(ctx.Args()) == 0 && !ctx.Bool("stdin") {
		return cli.ShowCommandHelp(ctx, "write")
	}

	fields, err := fieldsFromCtx(ctx)
	if err != nil {
		return err
	}

	req := fixity.NewWrite(ctx.String("id"), nil)
	req.Fields = fields

	if rollSize := ctx.Int("manual-rollsize"); rollSize != 0 {
		req.RollSize = int64(rollSize)
	}

	if ctx.Bool("cli") {
		req.Blob = ioutil.NopCloser(strings.NewReader(strings.Join(ctx.Args(), " ")))
	} else if ctx.Bool("stdin") {
		req.Blob = ioutil.NopCloser(os.Stdout)
	} else {
		path := ctx.Args().First()

		fi, err := os.Stat(path)
		if err != nil {
			return err
		}

		// TODO(leeola): append unix metadata to fields array

		if req.RollSize == fixity.DefaultRollSize {
			req.SetRollFromFileInfo(fi)
		}

		f, err := os.OpenFile(path, os.O_RDONLY, 0644)
		if err != nil {
			return err
		}
		req.Blob = f
	}

	fixi, err := loadFixity(ctx)
	if err != nil {
		return err
	}

	c, err := fixi.WriteRequest(req)
	if err != nil {
		return err
	}

	b, err := c.Blob()
	if err != nil {
		return err
	}

	out := os.Stdout
	inspect := ctx.Bool("inspect")
	spamBytes := ctx.Bool("spam-bytes")

	for _, h := range b.ChunkHashes {
		fmt.Fprintln(out, h)
		if inspect {
			var c fixity.Chunk
			if err := fixityutil.ReadAndUnmarshal(fixi, h, &c); err != nil {
				return err
			}

			if len(c.ChunkBytes) > 50 && !spamBytes {
				c.ChunkBytes = []byte("...spam bytes hidden...")
			}

			if err := printStruct(out, c); err != nil {
				return err
			}
		}
	}

	fmt.Fprintln(out, b.Hash)
	if inspect {
		if err := printStruct(out, b); err != nil {
			return err
		}
	}

	fmt.Fprintln(out, c.Hash)
	if inspect {
		if err := printStruct(out, c); err != nil {
			return err
		}
	}

	return nil
}

func fieldsFromCtx(ctx *cli.Context) (fixity.Fields, error) {
	indexFields := ctx.StringSlice("index")
	ftsFields := ctx.StringSlice("fts")
	hasIndexFields := len(indexFields) > 0
	hasFtsFields := len(ftsFields) > 0

	if !hasIndexFields && !hasFtsFields {
		return nil, nil
	}

	var fields []fixity.Field
	for _, f := range indexFields {
		k, v := splitKeyValue(f)
		fields = append(fields, fixity.Field{
			Field: k,
			Value: v,
		})
	}

	for _, f := range ftsFields {
		k, v := splitKeyValue(f)
		fields = append(fields, fixity.Field{
			Field:   k,
			Value:   v,
			Options: (fixity.FieldOptions{}).FullTextSearch(),
		})
	}

	return fields, nil
}

func splitKeyValue(s string) (string, interface{}) {
	kv := strings.SplitN(s, "=", 2)
	k := kv[0]
	if len(kv) == 1 {
		return k, nil
	}

	sv := kv[1]
	if v, err := strconv.ParseBool(sv); err == nil {
		return k, v
	}
	if v, err := strconv.Atoi(sv); err == nil {
		return k, v
	}
	return k, sv
}

func printStruct(out io.Writer, v interface{}) error {
	b, err := json.Marshal(v)
	if err != nil {
		return err
	}

	return printJsonBytes(os.Stdout, b)
}
