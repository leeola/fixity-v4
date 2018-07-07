package main

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/urfave/cli"
)

func WriteCmd(clictx *cli.Context) error {
	useStdin := clictx.Bool("stdin")
	if !useStdin {
		return errors.New("only stdin currently supported")
	}

	s, err := storeFromCli(clictx)
	if err != nil {
		// no wrap above helper errs
		return err
	}

	id := "foo"

	hashes, err := s.Write(context.Background(), id, nil, strings.NewReader("foo"))
	if err != nil {
		return fmt.Errorf("write: %v", err)
	}

	for _, h := range hashes {
		fmt.Println(h)
	}

	return nil
}
