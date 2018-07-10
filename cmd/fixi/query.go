package main

import (
	"fmt"
	"os"
	"strings"
	"text/tabwriter"

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

	matches, err := s.Query(q.FromString(qStr))
	if err != nil {
		return fmt.Errorf("query: %v", err)
	}

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 1, ' ', 0)
	fmt.Fprintf(w, "\tREF\tID\t\n")
	for i, m := range matches {
		fmt.Fprintf(w, "%d\t%s\t%s\t\n", i+1, m.Ref, m.ID)
	}
	w.Flush()

	return nil
}
