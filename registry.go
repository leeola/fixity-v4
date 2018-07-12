package fixity

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"sync"

	"github.com/leeola/fixity/config"
	"github.com/leeola/fixity/q"
)

var (
	blobstoreRegistry   map[string]BlobstoreCreator
	blobstoreRegistryMu sync.Mutex
	indexRegistry       map[string]IndexCreator
	indexRegistryMu     sync.Mutex
	storeRegistry       map[string]StoreCreator
	storeRegistryMu     sync.Mutex
	registeredDefault   func() (config.Config, error)
)

type Writer interface {
	Write(context.Context, []byte) (Ref, error)
}

type BlobReader interface {
	Read(context.Context, Ref) (io.ReadCloser, error)
}

type ReadWriter interface {
	BlobReader
	Writer
}

type Input map[string]interface{}

type BlobstoreCreator interface {
	New(name string, c config.Config) (ReadWriter, error)
}

type IndexCreator interface {
	New(name string, c config.Config) (QueryIndexer, error)
}

type StoreCreator interface {
	New(name string, c config.Config) (Store, error)
}

type BlobstoreCreatorFunc func(string, config.Config) (ReadWriter, error)

type IndexCreatorFunc func(string, config.Config) (QueryIndexer, error)

type StoreCreatorFunc func(string, config.Config) (Store, error)

type Inputer interface {
	// Init a config with optional user input.
	Init(Input) (json.RawMessage, error)
}

type InputError struct {
	Field   string
	Message string
}

func RegisterBlobstore(key string, c BlobstoreCreator) {
	if key == "" {
		panic(fmt.Sprintf("key cannot be empty"))
	}

	blobstoreRegistryMu.Lock()
	defer blobstoreRegistryMu.Unlock()

	_, ok := blobstoreRegistry[key]
	if ok {
		panic(fmt.Sprintf("already registered blobstore: %s", key))
	}

	blobstoreRegistry[key] = c
}

func RegisterIndex(key string, c IndexCreator) {
	if key == "" {
		panic(fmt.Sprintf("key cannot be empty"))
	}

	indexRegistryMu.Lock()
	defer indexRegistryMu.Unlock()

	if _, ok := indexRegistry[key]; ok {
		panic(fmt.Sprintf("already registered index: %s", key))
	}

	indexRegistry[key] = c
}

func RegisterStore(key string, c StoreCreator) {
	if key == "" {
		panic(fmt.Sprintf("key cannot be empty"))
	}

	storeRegistryMu.Lock()
	defer storeRegistryMu.Unlock()

	if _, ok := storeRegistry[key]; ok {
		panic(fmt.Sprintf("already registered store: %s", key))
	}

	storeRegistry[key] = c
}

type Store interface {
	Blob(ctx context.Context, ref Ref) (io.ReadCloser, error)
	Read(ctx context.Context, id string) (Mutation, Values, Reader, error)
	ReadRef(context.Context, Ref) (Mutation, Values, Reader, error)
	Write(ctx context.Context, id string, v Values, r io.Reader) ([]Ref, error)
	Querier
}

func New() (Store, error) {
	//
	// TODO(leeola): try to load the config here
	//

	if registeredDefault == nil {
		return nil, errors.New("no default config generator specified")
	}

	// config doesn't exist, generate a default.
	c, err := registeredDefault()
	if err != nil {
		return nil, fmt.Errorf("defaultConfigGen: %v", err)
	}

	return NewFromConfig(c.Store, c)
}

func NewFromConfig(key string, c config.Config) (Store, error) {
	if key == "" {
		return nil, fmt.Errorf("missing required config: store")
	}

	tc, ok := c.StoreConfigs[key]
	if !ok {
		return nil, fmt.Errorf("store name not found: %q", key)
	}

	constructor, ok := storeRegistry[tc.Type]
	if !ok {
		return nil, fmt.Errorf("store type not found: %q", tc.Type)
	}

	return constructor.New(key, c)
}

func NewBlobstoreFromConfig(key string, c config.Config) (ReadWriter, error) {
	return nil, errors.New("not implemented")
}

func NewIndexFromConfig(key string, c config.Config) (QueryIndexer, error) {
	return nil, errors.New("not implemented")
}

type QueryIndexer interface {
	Indexer
	Querier
}

type Indexer interface {
	Index(mutRef Ref, m Mutation, d *DataSchema, v Values) error
}

// TODO(leeola): articulate a mechanism to query against unique ids or
// versions.
type Querier interface {
	Query(q.Query) ([]Match, error)
}

func SetDefaultConfig(f func() (config.Config, error)) {
	registeredDefault = f
}

func (f BlobstoreCreatorFunc) New(name string, c config.Config) (ReadWriter, error) {
	return f(name, c)
}

func (f IndexCreatorFunc) New(name string, c config.Config) (QueryIndexer, error) {
	return f(name, c)
}

func (f StoreCreatorFunc) New(name string, c config.Config) (Store, error) {
	return f(name, c)
}
