package dyntabwriter

import (
	"bytes"
	"fmt"
	"io"
)

type DynTabWriter struct {
	W       io.Writer
	columns []*Column
}

func New(w io.Writer) *DynTabWriter {
	return &DynTabWriter{
		W: w,
	}
}

func (w *DynTabWriter) Print(vs ...interface{}) error {
	b := make([][]byte, len(vs))
	for i, v := range vs {
		if v != nil {
			b[i] = []byte(fmt.Sprint(v))
		} else {
			b[i] = []byte{}
		}
	}
	_, err := w.Write(bytes.Join(b, []byte{'\t'}))
	return err
}

func (w *DynTabWriter) Write(p []byte) (int, error) {
	splits := bytes.Split(p, []byte{'\t'})

	columnTotal := len(w.columns)

	for i, split := range splits {
		if i+1 > columnTotal {
			w.columns = append(w.columns, NewColumn(' ', 2))
		}

		splits[i] = w.columns[i].RightPad(split)
	}

	return w.W.Write(append(bytes.Join(splits, []byte{' '}), '\n'))
}
