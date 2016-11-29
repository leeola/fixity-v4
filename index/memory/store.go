package memory

import "io"

func (m *Memory) Exists(s string) (bool, error) {
	return m.store.Exists(s)
}

func (m *Memory) Read(s string) (io.ReadCloser, error) {
	return m.store.Read(s)
}

func (m *Memory) Write(b []byte) (string, error) {
	h, err := m.store.Write(b)
	if err == nil {
		if err := m.AddEntry(h); err != nil {
			return "", err
		}
	}
	return h, err
}

func (m *Memory) WriteHash(h string, b []byte) error {
	err := m.store.WriteHash(h, b)
	if err == nil {
		if err := m.AddEntry(h); err != nil {
			return err
		}
	}
	return err
}
