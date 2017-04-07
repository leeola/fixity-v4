package kala

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"io"
	"time"

	"github.com/leeola/kala/index"
	"github.com/leeola/kala/store"
)

type Commit struct {
	Id                  string    `json:"id,omitempty"`
	PreviousVersionHash string    `json:"previousVersion,omitempty"`
	UploadedAt          time.Time `json:"uploadedAt,omitempty"`
	ChangeLog           string    `json:"changeLog,omitempty"`
}

// Meta is a type alias for store.Meta for the UX of Kala package users.
type Meta store.Meta

type Config struct {
	Index index.Index
	Store store.Store
}

type Kala struct {
	config Config
	index  index.Index
	store  store.Store
}

func New(c Config) (*Kala, error) {
	if c.Index == nil {
		return nil, errors.New("missing reqired config: Index")
	}
	if c.Store == nil {
		return nil, errors.New("missing reqired config: Store")
	}

	return &Kala{
		config: c,
		index:  c.Index,
		store:  c.Store,
	}, nil
}

func (k *Kala) Write(c Commit, m Meta, r io.Reader) ([]string, error) {
	// For quicker prototyping, only supporting metadata atm
	if r != nil {
		return nil, errors.New("reader not yet implemented")
	}

	if structs.IsZero(m) && r == nil {
		return nil, errors.New("No data given to write")
	}

	// k.store.Write(

	return nil, nil
}

// NewId is a helper to generate a new default length Id.
//
// Note that the Id is encoded as hex to easily interact with it, rather
// than plain bytes.
func NewId() string {
	b := make([]byte, 32)
	if _, err := io.ReadFull(rand.Reader, b); err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}
