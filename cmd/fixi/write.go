package main

import (
	"errors"
	"fmt"
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

	filePath := ctx.String("file")
	if filePath != "" {
		return errors.New("--file not implemented yet")
	}

	blobStr := ctx.String("blob")
	if blobStr != "" {
		return errors.New("--blob not implemented yet")
	}

	fixi, err := loadFixity(ctx)
	if err != nil {
		return err
	}

	b, err := clijson.CliJson(ctx.Args())
	if err != nil {
		return err
	}

	fields, err := jsonToFields(ctx, b) // might be nil
	if err != nil {
		return err
	}

	var c fixity.Commit
	j := fixity.MultiJson{}
	j.JsonBytesWithFields(ctx.String("json-key"), b, fields)

	hashes, err := fixi.Write(c, j, nil)
	if err != nil {
		return err
	}

	if ctx.Bool("print") {
		for _, h := range hashes {
			fmt.Println(h)
			if err := printHash(fixi, h); err != nil {
				return err
			}
		}
	} else {
		fmt.Println(strings.Join(hashes, "\n"))
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
			Value: v, // might be nil
		})
	}

	for _, f := range ftsFields {
		k, v := splitKeyValue(f)
		fields = append(fields, fixity.Field{
			Field:   k,
			Value:   v, // might be nil
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
