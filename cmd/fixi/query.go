package main

import (
	"fmt"
	"strings"

	"github.com/leeola/fixity/q"
	"github.com/urfave/cli"
)

func QueryCmd(clictx *cli.Context) error {
	s, err := storeFromCli(clictx)
	if err != nil {
		// no wrap above helper errs
		return err
	}

	qStr := strings.Join(clictx.Args(), " ")

	refs, err := s.Query(q.FromString(qStr))
	if err != nil {
		return fmt.Errorf("query: %v", err)
	}

	for _, ref := range refs {
		fmt.Println("ref:", ref)
	}

	return nil
}
