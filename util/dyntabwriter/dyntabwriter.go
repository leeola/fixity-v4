package dyntabwriter

import (
	"bytes"
	"fmt"
	"io"
)

const (
	splitChar = '\t'
	joinChar  = ' '
	padChar   = ' '
)

type DynTabWriter struct {
	bufferCount int
	columnB     *bytes.Buffer
	lineB       *bytes.Buffer
	W           io.Writer
	columns     []*Column
}

func New(w io.Writer) *DynTabWriter {
	return &DynTabWriter{
		bufferCount: -1,
		lineB:       bytes.NewBuffer(nil),
		columnB:     bytes.NewBuffer(nil),
		W:           w,
	}
}

func (w *DynTabWriter) Header(s ...string) error {
	w.bufferCount = 2
	v := make([]interface{}, len(s))
	for i, s := range s {
		v[i] = s
	}
	return w.Println(v...)
}

func (w *DynTabWriter) Println(vs ...interface{}) error {
	vsLen := len(vs)
	b := make([][]byte, vsLen)
	for i, v := range vs {
		if v != nil {
			b[i] = []byte(fmt.Sprint(v))
		} else {
			b[i] = []byte{}
		}
	}

	// this func is Println, so append a newline.
	b[vsLen-1] = append(b[vsLen-1], '\n')

	_, err := w.Write(bytes.Join(b, []byte{splitChar}))
	return err
}

func (w *DynTabWriter) writeLine(p []byte) (bool, error) {
	columnTotal := len(w.columns)
	splits := bytes.Split(p, []byte{splitChar})

	for i, split := range splits {
		if i+1 > columnTotal {
			w.columns = append(w.columns, NewColumn(padChar, 2))
		}

		splits[i] = w.columns[i].RightPad(split)
	}

	switch {
	case w.bufferCount > 0:
		w.bufferCount--
		// append the newline, since the column buffer will
		// eventually be sent to the main Write method.
		//
		// TODO(leeola): reorganize this writeLine logic, specifically
		// RightPad.. the flow of code is weird.
		if _, err := w.columnB.Write(append(p, '\n')); err != nil {
			return false, err
		}
		return true, nil

	case w.bufferCount == 0:
		w.bufferCount--
		if _, err := w.columnB.WriteTo(w); err != nil {
			return false, err
		}
	}

	if _, err := w.W.Write(bytes.Join(splits, []byte{joinChar})); err != nil {
		return false, err
	}

	return false, nil
}

func (w *DynTabWriter) Write(p []byte) (int, error) {
	n, err := w.lineB.Write(p)
	if err != nil {
		return n, err
	}

	for {
		b, err := w.lineB.ReadBytes('\n')
		if err != nil && err != io.EOF {
			return n, err
		}

		// if ReadBytes returned EOF, then the buffered data
		// is not a complete line. Write it back, until next time.
		if err == io.EOF {
			if len(b) > 0 {
				if _, err := w.lineB.Write(b); err != nil {
					return 0, err
				}
			}
			break
		}

		buffered, err := w.writeLine(b[:len(b)-1])
		if err != nil {
			return 0, err
		}

		if !buffered {
			// write a newline after each writeLine call.
			if _, err := w.W.Write([]byte{'\n'}); err != nil {
				return 0, err
			}
		}
	}

	return n, nil
}

func (w *DynTabWriter) Flush() error {
	b := append(w.columnB.Bytes(), w.lineB.Bytes()...)
	if len(b) == 0 {
		return nil
	}

	w.bufferCount = -1
	_, err := w.writeLine(b)
	return err
}
