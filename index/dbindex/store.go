package dbindex

import "io"

func (dbi *Dbindex) Exists(s string) (bool, error) {
	return dbi.store.Exists(s)
}

func (dbi *Dbindex) Read(s string) (io.ReadCloser, error) {
	return dbi.store.Read(s)
}

func (dbi *Dbindex) Write(b []byte) (string, error) {
	h, err := dbi.store.Write(b)
	if err == nil {
		if err := dbi.AddEntry(h); err != nil {
			return "", err
		}
	}
	return h, err
}

func (dbi *Dbindex) WriteHash(h string, b []byte) error {
	err := dbi.store.WriteHash(h, b)
	if err == nil {
		if err := dbi.AddEntry(h); err != nil {
			return err
		}
	}
	return err
}
