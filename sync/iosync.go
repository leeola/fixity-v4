package sync

import (
	"io"
	"io/ioutil"

	"github.com/leeola/fixity"
)

type SyncIo struct {
	fixi   fixity.Fixity
	id     string
	in     io.Reader
	out    io.Writer
	fields []fixity.Field

	hasValue bool
	c        fixity.Content
	err      error
}

func Io(fixi fixity.Fixity, id string, in io.Reader, out io.Writer, f ...fixity.Field) *SyncIo {
	return &SyncIo{
		fixi:   fixi,
		id:     id,
		in:     in,
		out:    out,
		fields: f,
	}
}

func (s *SyncIo) Next() bool {
	return !s.hasValue
}

func (s *SyncIo) Value() (c fixity.Content, err error) {
	if s.hasValue {
		return s.c, s.err
	}

	defer func() {
		s.c = c
		s.err = err
	}()

	req := fixity.NewWrite(s.id, ioutil.NopCloser(s.in))
	req.IgnoreDuplicateBlob = true

	s.hasValue = true
	c, err = s.fixi.WriteRequest(req)
	if err != nil {
		return fixity.Content{}, err
	}

	if err == nil && c.Index > 1 {
		c, err = s.fixi.Read(s.id)
		if err != nil {
			return fixity.Content{}, err
		}
	}

	rc, err := c.Read()
	if err != nil {
		return fixity.Content{}, err
	}
	defer rc.Close()

	if _, err := io.Copy(s.out, rc); err != nil {
		return fixity.Content{}, err
	}

	return c, err
}
