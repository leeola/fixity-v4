package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"

	"github.com/leeola/fixity"
	"github.com/leeola/fixity/util/clijson"
	"github.com/urfave/cli"
)

func WriteCmd(ctx *cli.Context) error {
	if len(ctx.Args()) == 0 {
		return cli.ShowCommandHelp(ctx, "write")
	}

	out := os.Stdout

	filePath := ctx.String("file")
	if filePath != "" {
		return errors.New("--file not implemented yet")
	}

	fixi, err := loadFixity(ctx)
	if err != nil {
		return err
	}

	jsonB, err := clijson.CliJson(ctx.Args())
	if err != nil {
		return err
	}

	fields, err := jsonToFields(ctx, jsonB)
	if err != nil {
		return err
	}

	c, err := fixi.Write(ctx.String("id"), bytes.NewReader(jsonB), fields...)
	if err != nil {
		return err
	}

	b, err := c.Blob()
	if err != nil {
		return err
	}

	inspect := ctx.Bool("inspect")
	fmt.Fprintln(out, c.Hash)

	if inspect {
		if err := printStruct(out, c); err != nil {
			return err
		}
	}

	fmt.Fprintln(out, b.Hash)
	if inspect {
		if err := printStruct(out, c); err != nil {
			return err
		}
	}

	// TODO(leeola): cap the total chunks printed to something small..
	// like 5. The ux should probably also be limited by an --unsafe flag.
	// Eg, printing even a single chunkhash might be massive if the rollsize
	// was set to 5MB or something.
	for _, h := range b.ChunkHashes {
		fmt.Fprintln(out, h)
		if inspect {
			if err := printHash(fixi, h); err != nil {
				return err
			}
		}
	}

	return nil
}

func jsonToFields(ctx *cli.Context, b []byte) (fixity.Fields, error) {
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
