package index

import "github.com/leeola/kala/q"

type Index interface {
	Search(q.Query) ([][]byte, error)
}
