package boltindex

import "io"

func (m *BoltIndex) Exists(s string) (bool, error) {
	return m.store.Exists(s)
}

func (m *BoltIndex) Read(s string) (io.ReadCloser, error) {
	return m.store.Read(s)
}

func (m *BoltIndex) Write(b []byte) (string, error) {
	h, err := m.store.Write(b)
	if err == nil {
		if err := m.AddEntry(h); err != nil {
			return "", err
		}
	}
	return h, err
}

func (m *BoltIndex) WriteHash(h string, b []byte) error {
	err := m.store.WriteHash(h, b)
	if err == nil {
		if err := m.AddEntry(h); err != nil {
			return err
		}
	}
	return err
}
